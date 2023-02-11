package basic

import "strings"

// StringView is a string wrapper which indicates an unowned view of an underlying string.
// The string must not be stored for later use, clone it instead.
type StringView string

// Clone returns a fresh copy allocated on the heap.
func (s StringView) Clone() string {
	return strings.Clone(string(s))
}

// Unwrap returns the underlying string.
func (s StringView) Unwrap() string {
	return string(s)
}
