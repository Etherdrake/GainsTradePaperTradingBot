package utils

import "strings"

func EscapeMarkdownV2(input string) string {
	charsToEscape := "_*[]()~>#+-=|{}.!"
	for _, char := range charsToEscape {
		input = strings.ReplaceAll(input, string(char), "\\"+string(char))
	}
	return input
}
