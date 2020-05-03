package utils

import (
	"log"
	"regexp"
	"strings"
)

// Strip Special Characters
func StripSpecialChars(str string, forceLower bool) string {
	if forceLower {
		str = strings.ToLower(str)
	}

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	return reg.ReplaceAllString(str, "_");
}
