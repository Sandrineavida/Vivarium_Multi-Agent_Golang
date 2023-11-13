package utils

// min returns the minimum of two integers.
func Intmin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum of two integers.
func Intmax(a, b int) int {
	if a > b {
		return a
	}
	return b
}
