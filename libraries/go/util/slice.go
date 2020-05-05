package util

// UniqueStr removes duplicates from a string slice
func UniqueStr(slice []string) []string {
	m := map[string]struct{}{}
	for _, s := range slice {
		m[s] = struct{}{}
	}

	result := make([]string, 0, len(m))
	for s := range m {
		result = append(result, s)
	}

	return result
}
