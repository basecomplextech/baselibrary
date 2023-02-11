package basic

import (
	"bytes"
)

// BytesView is a []byte wrapper which indicates an unowned view of an underlying byte slice.
// The slice must not be modified or stored for later use, clone it instead.
type BytesView []byte

// Clone returns a fresh copy allocated on the heap.
func (b BytesView) Clone() []byte {
	return bytes.Clone(b)
}

// Unwrap returns the underlying slice.
func (b BytesView) Unwrap() []byte {
	return []byte(b)
}
