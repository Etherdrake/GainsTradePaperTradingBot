package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ExtractTradeID(inputString string, splitString string) string {
	// Split the input string based on "/closeposition+"
	parts := strings.Split(inputString, splitString)

	var finalID string
	// Check if there are two parts (before and after "/closeposition+")
	if len(parts) == 2 {
		// Return the second part as the trade ID
		finalID = parts[1]
	} else {
		return ""
	}

	if strings.Contains(finalID, "realtrade") {
		finalID = strings.Trim(finalID, "realtrade")
	}

	if strings.Contains(finalID, "papertrade") {
		finalID = strings.Trim(finalID, "papertrade")
	}

	// Return an empty string if "/closeposition+" is not found or if there are more than two parts
	return finalID
}

func ExtractTradeIDPlusMethod(inputString string) string {
	// Find the first '+' character
	firstPlus := strings.Index(inputString, "+")
	if firstPlus == -1 {
		return "" // No '+' found
	}

	// Find the second '+' character after the first one
	secondPlus := strings.Index(inputString[firstPlus+1:], "+")
	if secondPlus == -1 {
		return "" // Only one '+' found
	}

	// Extract the trade ID between the two '+' characters
	tradeID := inputString[firstPlus+1 : firstPlus+1+secondPlus]

	return tradeID
}

func ExtractSendIdxPlusMethod(inputString string) string {
	// Find the first '+' character
	firstPlus := strings.Index(inputString, "+")
	if firstPlus == -1 {
		return "" // No '+' found
	}

	// Extract the trade ID between the two '+' characters
	tradeID := inputString[0 : firstPlus+1]

	return tradeID
}

// ExtractIndex extracts the part of the string that comes after "index="
func ExtractIndex(longString string) string {
	// Find the position of "index="
	index := strings.Index(longString, "index=")

	if index != -1 {
		// Extract the part after "index="
		result := longString[index+len("index="):]
		return result
	}

	// Return an empty string if "index=" is not found
	return ""
}

func ExtractNumberPlus(s string) (int, error) {
	re := regexp.MustCompile(`\+(\d+)`)
	match := re.FindStringSubmatch(s)

	if len(match) > 1 {
		return strconv.Atoi(match[1])
	}

	return 0, fmt.Errorf("No number found after +")
}
