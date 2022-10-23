package utils

import "strconv"

// ToPointer takes any value and returns a pointer to that value.
//
// This function is useful when you want to pass a literal value to a
// function that takes a pointer to a value.
func ToPointer[T any](v T) *T {
	return &v
}

// SliceContains returns true if the slice holds the value given value
func SliceContains[T any](slice []T, v T, eq func(T, T) bool) bool {
	for _, s := range slice {
		if eq(s, v) {
			return true
		}
	}

	return false
}

// BoolToString returns "true" if the given boolean is true, otherwise
// returns "false".
func BoolToString(b bool) string {
	if b {
		return "true"
	}

	return "false"
}

// IntToString returns the string representation of the given integer.
func IntToString(i int) string {
	return strconv.Itoa(i)
}
