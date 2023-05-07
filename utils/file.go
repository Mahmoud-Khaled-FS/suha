package utils

import (
	"errors"
	"fmt"
	"os"
)

func IsFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return !errors.Is(err, os.ErrNotExist)
}

func PanicIfExist(fileName string) {
	if IsFileExist(fileName) {
		PanicRed(fmt.Sprintf("ERROR: File '%s' already exist", fileName))
	}
}

func CreateDir(dir string) {
	err := os.MkdirAll(dir, 0750)
	if err != nil {
		PanicRed(fmt.Sprintf("ERROR: make dir failed\n%s", err))
	}
}