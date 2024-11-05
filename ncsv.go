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
	"unicode/utf8"
)

type LineData struct {
	LineNumber int
	Record     []string
}

func cleanInvalidUTF8(s string) string {
	if utf8.ValidString(s) {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			// Invalid UTF-8 byte sequence; skip or replace
			i++
			// Uncomment the next line to replace invalid bytes with '�' (Unicode replacement character)
			// b.WriteRune('�')
			continue // Skip invalid byte
		}
		b.WriteRune(r)
		i += size
	}
	return b.String()
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
	cleanStringArg := flag.Bool("clean-string", false, "Enable UTF-8 text cleaning (may slow down processing)")
	cleanStringArgShort := flag.Bool("l", false, "Enable UTF-8 text cleaning (short)")

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

	// Determine clean string
	cleanString := *cleanStringArg || *cleanStringArgShort

	if cleanString {
		fmt.Println("Warning: UTF-8 text cleaning is enabled. Processing will take considerably longer.")
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

				if cleanString {
					// Clean each field in the record
					for i := range record {
						record[i] = cleanInvalidUTF8(record[i])
					}
				}

				// Adjust record length
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
