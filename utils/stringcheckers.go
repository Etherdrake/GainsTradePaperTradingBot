package utils

import (
	"strings"
	"unicode"
)

func IsFirstThreeCaps(s string) bool {
	if len(s) < 4 || s[0] != '/' {
		return false
	}

	// Remove the first character "/"
	s = s[1:]

	for i := 0; i < 3; i++ {
		if !unicode.IsUpper(rune(s[i])) {
			return false
		}
	}
	return true
}

func ExtractPair(s string) string {
	index := strings.Index(s, "/")
	if index == -1 || index == len(s)-1 {
		return ""
	}
	return s[index+1:]
}
