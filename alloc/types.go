package alloc

import (
	"bytes"
	"strings"
)

// Bytes is a []byte wrapper which indicates the slice is allocated outside of the heap.
// The slice must not be modified or stored for later use, clone it instead.
type Bytes []byte

// Clone returns a fresh copy allocated on the heap.
func (b Bytes) Clone() []byte {
	return bytes.Clone(b)
}

// Unwrap returns the underlying slice.
func (b Bytes) Unwrap() []byte {
	return []byte(b)
}

// String

// String is a string wrapper which indicates the string is allocated outside of the heap.
// The string must not be stored for later use, clone it instead.
type String string

// Clone returns a fresh copy allocated on the heap.
func (s String) Clone() string {
	return strings.Clone(string(s))
}

// Unwrap returns the underlying string.
func (s String) Unwrap() string {
	return string(s)
}
