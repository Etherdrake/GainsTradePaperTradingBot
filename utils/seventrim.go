package utils

func TrimToSeven(inputString string) string {
	// Check if the input string is shorter than 7 characters
	if len(inputString) < 7 {
		return inputString
	}

	// Use string slicing to get the first 7 characters
	result := inputString[:7]
	return result
}
