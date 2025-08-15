package utils

// trimToLength trims the input string to the specified length
func TrimToLength(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length]
}
