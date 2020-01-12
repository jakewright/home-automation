package util

// ExactlyOne returns true iff one of the given booleans is true
func ExactlyOne(bs ...bool) bool {
	n := 0
	for _, b := range bs {
		if b {
			n++
		}
	}
	return n == 1
}
