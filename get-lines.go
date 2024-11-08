package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	// Command-line flags
	searchPtr := flag.String("search", "", "Comma-separated list of texts to search for (case-insensitive)")
	searchAlias := flag.String("s", "", "Alias for --search")

	filePtr := flag.String("file", "", "Input file name")
	fileAlias := flag.String("f", "", "Alias for --file")

	outputPtr := flag.String("output", "", "Output folder name")
	outputAlias := flag.String("o", "", "Alias for --output")

	workersPtr := flag.Int("workers", 1, "Number of workers for parallel processing")
	workersAlias := flag.Int("w", 1, "Alias for --workers")

	helpPtr := flag.Bool("help", false, "Display help information")
	helpAlias := flag.Bool("h", false, "Alias for --help")

	flag.Parse()

	// Handle help flag
	if *helpPtr || *helpAlias {
		flag.Usage()
		return
	}

	// Consolidate aliases
	searchValue := *searchPtr
	if searchValue == "" {
		searchValue = *searchAlias
	}

	fileName := *filePtr
	if fileName == "" {
		fileName = *fileAlias
	}

	outputName := *outputPtr
	if outputName == "" {
		outputName = *outputAlias
	}

	workers := *workersPtr
	if *workersAlias != 1 {
		workers = *workersAlias
	}

	// Validate inputs
	if searchValue == "" || fileName == "" {
		fmt.Println("Error: --search and --file are required.")
		flag.Usage()
		return
	}

	if workers <= 0 {
		fmt.Println("Error: --workers must be greater than 0.")
		return
	}

	// Notify user about unordered output if workers > 1
	if workers > 1 {
		fmt.Println("Note: Output will not be in the same order as the input due to parallel processing.")
	}

	// Parse search terms
	rawTerms := strings.Split(searchValue, ",")
	searchTerms := make([]string, 0, len(rawTerms))
	for _, term := range rawTerms {
		trimmedTerm := strings.TrimSpace(strings.ToLower(term))
		if trimmedTerm != "" {
			searchTerms = append(searchTerms, trimmedTerm)
		}
	}

	if len(searchTerms) == 0 {
		fmt.Println("Error: No valid search terms provided.")
		return
	}

	// Open input file
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error: Unable to open file %s\n", fileName)
		return
	}
	defer file.Close()

	// Prepare outputs
	var outputFiles map[string]*os.File
	var outputChans map[string]chan string
	var outputWg sync.WaitGroup

	outputFiles = make(map[string]*os.File)
	outputChans = make(map[string]chan string)

	if outputName != "" {
		// Create output directory if it doesn't exist
		err = os.MkdirAll(outputName, os.ModePerm)
		if err != nil {
			fmt.Printf("Error: Unable to create output directory %s\n", outputName)
			return
		}

		for _, term := range searchTerms {
			outputFilePath := filepath.Join(outputName, term+".txt")
			outputFile, err := os.Create(outputFilePath)
			if err != nil {
				fmt.Printf("Error: Unable to create output file %s\n", outputFilePath)
				return
			}
			outputFiles[term] = outputFile

			// Create a channel and start a goroutine to write to the file
			outputChans[term] = make(chan string, workers*2)
			outputWg.Add(1)
			go func(ch chan string, file *os.File) {
				defer outputWg.Done()
				for line := range ch {
					_, err := file.WriteString(line)
					if err != nil {
						fmt.Printf("Error: Unable to write to output file %s\n", file.Name())
						return
					}
				}
			}(outputChans[term], outputFile)
		}
	} else {
		// Prepare channels for console output
		for _, term := range searchTerms {
			outputChans[term] = make(chan string, workers*2)
			outputWg.Add(1)
			go func(ch chan string, term string) {
				defer outputWg.Done()
				for line := range ch {
					fmt.Printf("[%s] %s", term, line)
				}
			}(outputChans[term], term)
		}
	}

	// Channels and WaitGroup for concurrency
	lines := make(chan string, workers*2)
	var workerWg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < workers; i++ {
		workerWg.Add(1)
		go func() {
			defer workerWg.Done()
			for line := range lines {
				parts := strings.SplitN(line, "\t", 2)
				lineNum := parts[0]
				content := parts[1]
				lowerContent := strings.ToLower(content)

				for _, term := range searchTerms {
					if strings.Contains(lowerContent, term) {
						resultLine := fmt.Sprintf("%s\t%s\n", lineNum, content)
						outputChans[term] <- resultLine
					}
				}
			}
		}()
	}

	// Reading lines from file
	go func() {
		scanner := bufio.NewScanner(file)
		lineNumber := 1
		for scanner.Scan() {
			line := fmt.Sprintf("%d\t%s", lineNumber, scanner.Text())
			lines <- line
			lineNumber++
		}
		close(lines)
	}()

	// Wait for workers to finish processing
	workerWg.Wait()

	// Close output channels
	for _, ch := range outputChans {
		close(ch)
	}

	// Wait for output goroutines to finish
	outputWg.Wait()

	// Close output files
	if outputName != "" {
		for _, file := range outputFiles {
			file.Close()
		}
	}
}
