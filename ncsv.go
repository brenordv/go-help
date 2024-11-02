package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Normalize CSV tool
// 1. if --header <value> or -e <value> is passed, will use value as the header. This must be a string with the header columns separated by commas.
// 2. If --header(-e) is not passed, will use the first line of the file as the header.
// 3. The filename will be informed with the argument --file <path to file> (or -f <path to file>). If not informed, will fail execution.
// 4. With the header info, will read all the other lines, and try to fit the data into the header columns.
// 5. If a line has more columns than the header, will ignore the extra columns.
// 6. If a line has less columns than the header, will fill the missing columns with empty strings.
// 7. Any time there's a mismatch between the header and the line, will add that info to a mismatches.text file.
// 8. In the mismatches file, will add the line number, the header columns, and the line columns.
// 9. The output will be a new file with the same name as the input file, but with the suffix "_normalized".
// 10. The output file will have the same header as the input file.
// 11. The output file will have all the lines normalized according to the header.
// 12. Optionally, the argument --split <number> (or -s <number)> can be passed, to split the output file in multiple files with the number of lines passed. (Pay attention to odd numbers. In the end 1005 of the lines of the input file must be saved to the output files.)
// 13. When splitting the output, the resulting files will have the same name as the input file, but with the suffix "_normalized_<number>".
// 14. The output file or files will always have the extension csv.
// 15. Make safe use of channels, wait groups, and go routines to process the data. Always assume the input file can be huge and it must be done as fast as possible.
// 16. While processing the file, print the progress to the console. You can use the fmt package for that and show the number of lines processed.

type LineData struct {
	LineNumber int
	Record     []string
}

func main() {
	// Start time measurement
	startTime := time.Now()

	// Parse command-line arguments
	headerArg := flag.String("header", "", "Header columns separated by commas")
	headerArgShort := flag.String("e", "", "Header columns separated by commas (short)")
	fileArg := flag.String("file", "", "Path to the input CSV file")
	fileArgShort := flag.String("f", "", "Path to the input CSV file (short)")
	splitArg := flag.Int("split", 0, "Number of lines per output file")
	splitArgShort := flag.Int("s", 0, "Number of lines per output file (short)")
	printEveryArg := flag.Int("print-every", 500, "Print status every X lines")
	printEveryArgShort := flag.Int("p", 500, "Print status every X lines (short)")
	concurrencyArg := flag.Int("concurrency", 1000, "Concurrency level (buffer sizes)")
	concurrencyArgShort := flag.Int("c", 1000, "Concurrency level (buffer sizes) (short)")

	flag.Parse()

	// Determine header value
	headerValue := *headerArg
	if headerValue == "" {
		headerValue = *headerArgShort
	}

	// Determine file path
	filePath := *fileArg
	if filePath == "" {
		filePath = *fileArgShort
	}
	if filePath == "" {
		fmt.Println("Error: --file (-f) argument is required")
		os.Exit(1)
	}

	// Determine split lines
	splitLines := *splitArg
	if splitLines == 0 {
		splitLines = *splitArgShort
	}

	// Determine print every
	printEvery := *printEveryArg
	if printEvery == 500 {
		printEvery = *printEveryArgShort
	}

	// Determine concurrency
	concurrency := *concurrencyArg
	if concurrency == 1000 {
		concurrency = *concurrencyArgShort
	}

	// Open the input file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Read header
	var header []string
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow variable number of fields
	lineNumber := 0

	if headerValue != "" {
		header = strings.Split(headerValue, ",")
	} else {
		firstLine, err := reader.Read()
		if err != nil {
			fmt.Printf("Error reading header: %v\n", err)
			os.Exit(1)
		}
		lineNumber++
		header = firstLine
	}

	// Prepare output file naming
	inputFileName := filepath.Base(filePath)
	inputFileExt := filepath.Ext(inputFileName)
	inputFileNameWithoutExt := strings.TrimSuffix(inputFileName, inputFileExt)
	outputFileNameTemplate := fmt.Sprintf("%s_normalized", inputFileNameWithoutExt)
	if splitLines > 0 {
		outputFileNameTemplate += "_%d.csv"
	} else {
		outputFileNameTemplate += ".csv"
	}

	// Prepare mismatches file
	mismatchesFile, err := os.Create("mismatches.text")
	if err != nil {
		fmt.Printf("Error creating mismatches file: %v\n", err)
		os.Exit(1)
	}
	defer mismatchesFile.Close()
	mismatchesWriter := bufio.NewWriter(mismatchesFile)
	defer mismatchesWriter.Flush()

	// Channels and wait groups for concurrency
	linesChan := make(chan LineData, concurrency)
	resultsChan := make(chan LineData, concurrency)
	mismatchesChan := make(chan string, concurrency)
	var wg sync.WaitGroup

	// Worker pool
	numWorkers := runtime.NumCPU()
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for lineData := range linesChan {
				record := lineData.Record
				lineNum := lineData.LineNumber

				if len(record) > len(header) {
					record = record[:len(header)]
					mismatchesChan <- fmt.Sprintf("Line %d: Extra columns detected. Header: %d columns, Line: %d columns\n", lineNum, len(header), len(lineData.Record))
				} else if len(record) < len(header) {
					for len(record) < len(header) {
						record = append(record, "")
					}
					mismatchesChan <- fmt.Sprintf("Line %d: Missing columns detected. Header: %d columns, Line: %d columns\n", lineNum, len(header), len(lineData.Record))
				}
				resultsChan <- LineData{
					LineNumber: lineNum,
					Record:     record,
				}
			}
		}()
	}

	// Goroutine to write mismatches
	go func() {
		for mismatch := range mismatchesChan {
			mismatchesWriter.WriteString(mismatch)
		}
	}()

	// Read lines and send to linesChan
	go func() {
		for {
			record, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Printf("Error reading line %d: %v\n", lineNumber+1, err)
				continue
			}
			lineNumber++
			linesChan <- LineData{
				LineNumber: lineNumber,
				Record:     record,
			}
			if lineNumber%printEvery == 0 {
				fmt.Printf("\rProcessed %d lines", lineNumber)
				// Flush stdout to update the console
				os.Stdout.Sync()
			}
		}
		close(linesChan)
	}()

	// Close channels after processing
	go func() {
		wg.Wait()
		close(resultsChan)
		close(mismatchesChan)
	}()

	// Prepare output files
	var outputFile *os.File
	var outputWriter *csv.Writer
	currentFileIndex := 1
	currentLineCount := 0

	createNewOutputFile := func() {
		if outputFile != nil {
			outputWriter.Flush()
			outputFile.Close()
		}
		var outputFileName string
		if splitLines > 0 {
			outputFileName = fmt.Sprintf(outputFileNameTemplate, currentFileIndex)
		} else {
			outputFileName = outputFileNameTemplate
		}
		var err error
		outputFile, err = os.Create(outputFileName)
		if err != nil {
			fmt.Printf("Error creating output file: %v\n", err)
			os.Exit(1)
		}
		outputWriter = csv.NewWriter(outputFile)
		outputWriter.Write(header)
		currentFileIndex++
		currentLineCount = 0
	}

	createNewOutputFile()
	totalLinesWritten := 0

	// Write processed lines to output
	for lineData := range resultsChan {
		if splitLines > 0 && currentLineCount >= splitLines {
			createNewOutputFile()
		}
		outputWriter.Write(lineData.Record)
		currentLineCount++
		totalLinesWritten++
		if totalLinesWritten%printEvery == 0 {
			fmt.Printf("\rTotal lines written: %d", totalLinesWritten)
			// Flush stdout to update the console
			os.Stdout.Sync()
		}
	}

	// Finalize
	outputWriter.Flush()
	mismatchesWriter.Flush()
	if outputFile != nil {
		outputFile.Close()
	}

	// Print final status and elapsed time
	elapsedTime := time.Since(startTime)
	fmt.Printf("\rProcessing completed. Total lines processed: %d\n", lineNumber)
	fmt.Printf("Elapsed time: %s\n", elapsedTime.Round(time.Millisecond))
}
