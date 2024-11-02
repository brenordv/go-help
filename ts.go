package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"
)

// List of common date-time formats to guess from
var layouts = []string{
	"2006-01-02T15:04:05",       // YYYY-MM-DDTHH:MM:SS
	"2006-01-02T15:04:05-07:00", // YYYY-MM-DDTHH:MM:SS<timezone>
	"2006-01-02 15:04:05",       // YYYY-MM-DD HH:MM:SS
	"2006-01-02 15:04",          // YYYY-MM-DD HH:MM
	"2006-01-02",                // YYYY-MM-DD
	"02-01-2006 15:04:05",       // DD-MM-YYYY HH:MM:SS
	"02-01-2006 15:04",          // DD-MM-YYYY HH:MM
	"02-01-2006",                // DD-MM-YYYY
	"2006/01/02 15:04:05",       // YYYY/MM/DD HH:MM:SS
	"2006/01/02 15:04",          // YYYY/MM/DD HH:MM
	"2006/01/02",                // YYYY/MM/DD
	"02/01/2006 15:04:05",       // DD/MM/YYYY HH:MM:SS
	"02/01/2006 15:04",          // DD/MM/YYYY HH:MM
	"02/01/2006",                // DD/MM/YYYY
	time.RFC3339,                // RFC3339 (2006-01-02T15:04:05Z07:00)
	time.UnixDate,               // Unix date format
	time.RubyDate,               // Ruby date format
}

// Prints the help message with usage instructions
func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("  no parameters          : Prints the current time as a Unix timestamp")
	fmt.Println("  <unix timestamp>       : Converts the Unix timestamp to UTC and local date time")
	fmt.Println("  <YYYY-MM-DD HH:MM:SS>  : Converts the date time to a Unix timestamp")
	fmt.Println("  -h, --help             : Prints this help message")
}

// Prints the current time as a Unix timestamp
func printCurrentUnixTimestamp() {
	currentTime := time.Now().Unix()
	fmt.Printf("%d", currentTime)
}

// Converts a Unix timestamp to UTC and local time
func convertUnixToDateTime(unixTimestamp int64) {
	t := time.Unix(unixTimestamp, 0)
	fmt.Printf("UTC Time    : %s\n", t.UTC().Format(time.RFC3339))
	fmt.Printf("Local Time  : %s\n", t.Local().Format(time.RFC3339))
}

// Converts a date-time string to a Unix timestamp
func convertDateTimeToUnix(dateTimeStr string) {
	// First, try the default format
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, dateTimeStr)
	if err != nil {
		// If the default format fails, attempt to guess the format
		t, err = guessDateTimeFormat(dateTimeStr)
		if err != nil {
			fmt.Println("Invalid date-time format. Unable to parse the input.")
			return
		}
	}
	fmt.Printf("Unix Timestamp: %d\n", t.Unix())
}

// Tries to guess the date-time format from a list of common formats
func guessDateTimeFormat(input string) (time.Time, error) {
	for _, layout := range layouts {
		if t, err := time.Parse(layout, input); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unknown date-time format")
}

func main() {
	// Parse flags
	help := flag.Bool("h", false, "Show help")
	helpLong := flag.Bool("help", false, "Show help")
	flag.Parse()

	// Show help if -h or --help is used
	if *help || *helpLong {
		printHelp()
		return
	}

	// No arguments, print current time as Unix timestamp
	if len(flag.Args()) == 0 {
		printCurrentUnixTimestamp()
		return
	}

	input := flag.Args()[0]

	// Try to parse the input as a Unix timestamp (integer)
	if unixTimestamp, err := strconv.ParseInt(input, 10, 64); err == nil {
		convertUnixToDateTime(unixTimestamp)
		return
	}

	// If not a Unix timestamp, treat it as a date-time string
	convertDateTimeToUnix(input)
}
