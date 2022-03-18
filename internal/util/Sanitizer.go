package util

import (
	"regexp"
	"strings"
)

func Sanitize(ss []string) []string {
	return Map(ss, SanitizeString)
}

func SanitizeString(s string) string {
	allowedChars := regexp.MustCompile(`[^a-zA-Z0-9\.\-@]`)
	return strings.ReplaceAll(allowedChars.ReplaceAllString(s, ""), "..", "")
}
