package utils

import "strings"

func RemoveQuotes(inputString string) string {
	return strings.ReplaceAll(inputString, "\"", "")
}
