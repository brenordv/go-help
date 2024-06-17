package main

import (
	"flag"
	"fmt"
	"github.com/brenordv/go-help/constants"
	"io"
	"os"
)

// Usage function to display custom help message
func usage() {
	fmt.Println("Mimics the behavior of the Linux 'cat' command.")
	fmt.Println(constants.StringLine)
	fmt.Println("Usage: gocat [-i] [filename]...")
	fmt.Println("  -i Read from stdin if no files are provided")
	fmt.Println(constants.StringLine)
	fmt.Println("version:", constants.AppCatVersion)
	flag.PrintDefaults()
}

func main() {
	iFlag := flag.Bool("i", false, "Execute with -i to read from stdin if no files are provided")
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()

	if *iFlag && len(args) == 0 {
		// If -i is specified or no files are provided, read from stdin
		_, err := io.Copy(os.Stdout, os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
		}
		return
	}

	// Handle each file provided as an argument
	for _, filename := range args {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", filename, err)
			continue
		}
		defer file.Close()

		_, err = io.Copy(os.Stdout, file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", filename, err)
		}
	}

}
