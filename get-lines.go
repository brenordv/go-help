package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Result struct {
	LineNumber int
	Line       string
}

func main() {
	// Define flags
	search := flag.String("search", "", "Text to search for (required)")
	file := flag.String("file", "", "File to search in (required)")
	output := flag.String("output", "", "Output file (optional)")
	workers := flag.Int("workers", 1, "Number of parallel workers (default 1)")
	help := flag.Bool("help", false, "Show help")
	flag.Parse()

	// Show help if requested
	if *help || *search == "" || *file == "" {
		fmt.Println("Usage: goapp --search <value> --file <filename> [--output <filename>] [--workers <number>]")
		fmt.Println("Example: goapp -s \"search text\" -f input.txt -o output.txt -w 4")
		fmt.Println("Options:")
		fmt.Println("  -s, --search     Text to search for (required)")
		fmt.Println("  -f, --file       File to search in (required)")
		fmt.Println("  -o, --output     Output file (optional)")
		fmt.Println("  -w, --workers    Number of parallel workers (default 1)")
		fmt.Println("  -h, --help       Show help")
		return
	}

	// Ensure workers count is valid
	if *workers <= 0 {
		fmt.Println("Error: workers must be greater than 0.")
		return
	}

	// Notify user about unordered output if using multiple workers
	if *workers > 1 {
		fmt.Println("Warning: Using multiple workers, output order will not match input order.")
	}

	// Open file for reading
	fileHandle, err := os.Open(*file)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer fileHandle.Close()

	// Prepare output (console or file)
	var outputWriter *os.File
	if *output != "" {
		outputWriter, err = os.Create(*output)
		if err != nil {
			fmt.Printf("Error creating output file: %v\n", err)
			return
		}
		defer outputWriter.Close()
	} else {
		outputWriter = os.Stdout
	}

	// Convert search term to lowercase for case-insensitive search
	searchTerm := strings.ToLower(*search)

	// Channel for distributing lines to workers
	lines := make(chan Result, *workers)
	var wg sync.WaitGroup

	// Launch workers
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for lineResult := range lines {
				if strings.Contains(strings.ToLower(lineResult.Line), searchTerm) {
					fmt.Fprintf(outputWriter, "%d\t%s\n", lineResult.LineNumber, lineResult.Line)
				}
			}
		}()
	}

	// Read lines and send to workers
	scanner := bufio.NewScanner(fileHandle)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		lines <- Result{LineNumber: lineNumber, Line: scanner.Text()}
	}
	close(lines)
	wg.Wait()

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
}
