package arena

import "unsafe"

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

	arena := a.(*arena)
	ptr := arena.alloc(size)
	return (*T)(ptr)
}

// Bytes allocates a new byte slice.
func Bytes(a Arena, cap int) []byte {
	if cap == 0 {
		return nil
	}

	arena := a.(*arena)
	ptr := arena.alloc(cap)
	return unsafe.Slice((*byte)(ptr), cap)
}

// Slice allocates a new slice of a generic type.
//
// Usage:
//
//	var foo []MyStruct
//	foo = Slice[MyStruct](arena, 16)
func Slice[T any](a Arena, cap int) []T {
	if cap == 0 {
		return nil
	}

	var zero T
	elem := int(unsafe.Sizeof(zero))
	size := elem * cap

	arena := a.(*arena)
	ptr := arena.alloc(size)
	return unsafe.Slice((*T)(ptr), cap)
}

// Copy allocates a new slice and copies items from src into it.
// The slice capacity is len(src).
func Copy[T any](a Arena, src []T) []T {
	dst := Slice[T](a, len(src))
	copy(dst, src)
	return dst
}

// CopyBytes allocates a new byte slice and copies items from src into it.
// The slice capacity is len(src).
func CopyBytes(a Arena, src []byte) []byte {
	dst := Bytes(a, len(src))
	copy(dst, src)
	return dst
}

// String allocates a new string and copies data from src into it.
func String(a Arena, src string) string {
	if len(src) == 0 {
		return ""
	}

	dst := Bytes(a, len(src))
	copy(dst, src)
	return *(*string)(unsafe.Pointer(&dst))
}
