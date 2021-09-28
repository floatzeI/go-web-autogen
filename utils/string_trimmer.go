package utils

import "strings"

// StringTrimmer removes any beginning/ending whitespace, and start/end quotes, from a string
func StringTrimmer(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "\"")
	s = strings.TrimSuffix(s, "\"")
	return s
}
