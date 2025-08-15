package utils

import "strings"

func CutSlash(s string) string {
	parts := strings.Split(s, "/")
	return parts[0]
}
