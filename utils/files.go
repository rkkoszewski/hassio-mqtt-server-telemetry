package utils

import (
	"io/ioutil"
	"os"
)

// Read a file to string
func ReadFileToString(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// Check if file can be read
func CanReadFile(path string) bool {
	_, err := ReadFileToString(path)
	return err == nil
}