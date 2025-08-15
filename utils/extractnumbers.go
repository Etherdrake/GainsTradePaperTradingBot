package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

func ExtractNumber(s string) (int, error) {
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(s)
	if match == "" {
		return 0, fmt.Errorf("No number found in the string")
	}
	number, err := strconv.Atoi(match)
	if err != nil {
		return 0, err
	}
	return number, nil
}
