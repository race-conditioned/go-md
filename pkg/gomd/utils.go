package gomd

// inlineWrap concatenates the wrap to each side of s.
func inlineWrap(wrap, s string) string {
	return wrap + s + wrap
}

// btoi converts a boolean to an integer.
func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
