package utils

import (
	"io/ioutil"
	"os"
	"os/exec"
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

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Check if Command Exists
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}