package alloc

import (
	"unsafe"
)

// Alloc allocates a new object and returns a pointer to it.
//
// Usage:
//
//	var foo *float64
//	var bar *MyStruct
//	foo = Alloc[float64](arena)
//	bar = Alloc[MyStruct](arena)
func Alloc[T any](a Arena) *T {
	var zero T
	size := int(unsafe.Sizeof(zero))

	ptr := a.Alloc(size)
	return (*T)(ptr)
}

// Bytes allocates a new byte slice.
func Bytes(a Arena, len int) []byte {
	if len == 0 {
		return nil
	}

	ptr := a.Alloc(len)
	return unsafe.Slice((*byte)(ptr), len)
}

// CopyBytes allocates a new byte slice and copies items from src into it.
// The slice capacity is len(src).
func CopyBytes(a Arena, src []byte) []byte {
	dst := Bytes(a, len(src))
	copy(dst, src)
	return dst
}
