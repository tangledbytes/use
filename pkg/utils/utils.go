package utils

// ToPointer takes any value and returns a pointer to that value.
//
// This function is useful when you want to pass a literal value to a
// function that takes a pointer to a value.
func ToPointer[T any](v T) *T {
	return &v
}
