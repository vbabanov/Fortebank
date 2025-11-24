package utils

import "strings"

func NormalizeString(s string) string {
	strings.ToLower(s)
	Regexp.ReplaceAllString(s, "")
	return cleaned_string
}
