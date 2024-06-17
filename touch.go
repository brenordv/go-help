package main

import (
	"flag"
	"fmt"
	"github.com/brenordv/go-help/constants"
	"github.com/brenordv/go-help/utils"
)

func main() {
	flag.Usage = func() {
		fmt.Println("Mimics the functionality of the 'touch' command in Linux.")
		fmt.Println(constants.StringLine)
		fmt.Println("Usage: touch [-t] <filename1> <filename2> ... <filenameN>")
		flag.PrintDefaults()
		fmt.Println(constants.StringLine+"\nversion: ", constants.AppTouchVersion)
	}
	tFlag := flag.Bool("t", false, "Execute with -t to display the current time after touching the file")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Error: Missing filename")
		flag.Usage()
		return
	}

	// Handle each file provided as an argument
	for _, filename := range flag.Args() {
		fileAlreadyExisted, err := utils.CreateFileIfNotExists(filename)

		if err != nil {
			fmt.Println("Error creating file '", filename, "': ", err)
			continue
		}

		if !fileAlreadyExisted {
			// If we just created the file, there's no need to update the timestamps
			err = utils.UpdateFileTimeStamps(filename)

			if err != nil {
				fmt.Println("Error updating file '", filename, "' timestamps: ", err)
				return
			}
		}

		if !*tFlag {
			continue
		}

		fmt.Println("File touched '", filename, "' at: ", utils.GetFileTimeStamp(filename))
	}
}
