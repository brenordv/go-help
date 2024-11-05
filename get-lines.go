package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
)

func main() {
	// Command-line flags
	searchPtr := flag.String("search", "", "Text to search for (case-insensitive)")
	searchAlias := flag.String("s", "", "Alias for --search")

	filePtr := flag.String("file", "", "Input file name")
	fileAlias := flag.String("f", "", "Alias for --file")

	outputPtr := flag.String("output", "", "Output file name")
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

	// Open input file
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error: Unable to open file %s\n", fileName)
		return
	}
	defer file.Close()

	// Prepare output
	var output *os.File
	if outputName != "" {
		output, err = os.Create(outputName)
		if err != nil {
			fmt.Printf("Error: Unable to create output file %s\n", outputName)
			return
		}
		defer output.Close()
	} else {
		output = os.Stdout
	}

	// Channels and WaitGroup for concurrency
	lines := make(chan string, workers*2)
	results := make(chan string, workers*2)
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range lines {
				parts := strings.SplitN(line, "\t", 2)
				lineNum := parts[0]
				content := parts[1]

				if strings.Contains(strings.ToLower(content), strings.ToLower(searchValue)) {
					results <- fmt.Sprintf("%s\t%s\n", lineNum, content)
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

	// Collect results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Write results to output
	for result := range results {
		_, err := output.WriteString(result)
		if err != nil {
			fmt.Printf("Error: Unable to write to output\n")
			return
		}
	}
}
