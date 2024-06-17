package utils

import (
	"errors"
	"os"
	"time"
)

func CreateFileIfNotExists(filename string) (bool, error) {
	if _, err := os.Stat(filename); err == nil {
		return true, nil

	} else if errors.Is(err, os.ErrNotExist) {
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return false, err
		}
		file.Close()
		return false, nil

	} else {
		// Some other error occurred, maybe permission or something like this.
		// not going to overthink this too much...
		return false, err

	}
}

func UpdateFileTimeStamps(filename string) error {
	now := time.Now()

	return os.Chtimes(filename, now, now)
}

func GetFileTimeStamp(filename string) time.Time {
	fileInfo, _ := os.Stat(filename)
	return fileInfo.ModTime()
}
