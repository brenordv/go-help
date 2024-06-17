package main

import (
	"flag"
	"fmt"
	"github.com/brenordv/go-help/constants"
	"github.com/brenordv/go-help/utils"
)

func main() {
	flag.Usage = func() {
		fmt.Println("Generates a new GUID and optionally copies it to the clipboard.")
		fmt.Println(constants.StringLine)
		fmt.Println("Usage: guid [-c]")
		flag.PrintDefaults()
		fmt.Println(constants.StringLine+"\nversion: ", constants.AppGuidVersion)
	}
	cFlag := flag.Bool("c", false, "Execute with -c to copy the guid to the clipboard")
	flag.Parse()

	guid := utils.NewGuid()
	fmt.Print(guid)

	if !*cFlag {
		return
	}

	utils.CopyToClipboard(guid)
}
