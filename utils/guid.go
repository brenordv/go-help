package utils

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/beevik/guid"
)

func NewGuid() string {
	return guid.NewString()
}

func CopyToClipboard(text string) {
	err := clipboard.WriteAll(text)

	if err == nil {
		return
	}

	fmt.Println("Error copying to clipboard: ", err)
}
