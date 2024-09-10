package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func IsFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	fmt.Println(err)
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

func SanitizeFileName(fileName string) string {
	forbiddenChars := []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*"}
	for _, char := range forbiddenChars {
		fileName = strings.ReplaceAll(fileName, char, " ")
	}
	return fileName
}
