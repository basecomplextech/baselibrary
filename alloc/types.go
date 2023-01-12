package alloc

import (
	"bytes"
	"strings"
)

// Bytes is a []byte wrapper which indicates that the underlying slice is stored outside of the heap.
type Bytes []byte

// Clone returns a fresh copy allocated on the heap.
func (b Bytes) Clone() []byte {
	return bytes.Clone(b)
}

// Unwrap returns the underlying slice.
func (b Bytes) Unwrap() []byte {
	return []byte(b)
}

// String is a string wrapper which indicates that the underlying string is stored outside of the heap.
type String string

// Clone returns a fresh copy allocated on the heap.
func (s String) Clone() string {
	return strings.Clone(string(s))
}

// Unwrap returns the underlying string.
func (s String) Unwrap() string {
	return string(s)
}
