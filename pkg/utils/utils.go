package utils

import (
	"os"
	"strconv"
)

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

// StringToInt returns the integer representation of the given string.
//
// If the string cannot be converted to an integer, this function panics.
func StringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return i
}

// GetEnvOrDefault returns the value of the given environment variable
// if it exists, otherwise returns the given default value.
func GetEnvOrDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}
