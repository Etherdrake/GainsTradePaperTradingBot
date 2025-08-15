package utils

import "strings"

// replaceSlashWithDash takes a string and replaces all '/' with '-'
func ReplaceSlashWithDash(s string) string {
	return strings.ReplaceAll(s, "/", "-")
}
