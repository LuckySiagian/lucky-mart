package utils

import (
	"regexp"
	"strings"
)

var (
	nonAlphanumericRegex = regexp.MustCompile(`[^a-z0-9]+`)
)

func GenerateSlug(s string) string {
	str := strings.ToLower(strings.TrimSpace(s))
	str = nonAlphanumericRegex.ReplaceAllString(str, "-")
	return strings.Trim(str, "-")
}
