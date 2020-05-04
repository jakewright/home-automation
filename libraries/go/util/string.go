package util

import "strings"

// RemoveWhitespaceStrings removes any strings from  the input slice that
// contain only whitespace and trims whitespace from the remaining lines.
func RemoveWhitespaceStrings(a []string) []string {
	result := make([]string, 0, len(a))

	for _, s := range a {
		s = strings.TrimSpace(s)

		if s == "" {
			continue
		}

		result = append(result, s)
	}

	return result
}
