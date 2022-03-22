package alloc

import (
	"unsafe"
)

// Alloc allocates a new object and returns a pointer to it.
//
// Usage:
//	var foo *float64
//	var bar *MyStruct
//	foo = Alloc[float64](arena)
//	bar = Alloc[MyStruct](arena)
//
func Alloc[T any](a *Arena) *T {
	var zero T

	size := int(unsafe.Sizeof(zero))
	ptr := a.alloc(size)
	return (*T)(ptr)
}

// AllocBytes allocates a byte slice with a `size` capacity in the arena.
func AllocBytes(a *Arena, cap int) []byte {
	if cap == 0 {
		return nil
	}

	ptr := a.alloc(cap)
	return unsafe.Slice((*byte)(ptr), cap)
}

// CopyBytes allocates a byte slice and copies items from src into it.
// The slice capacity is len(src).
func CopyBytes(a *Arena, src []byte) []byte {
	dst := AllocBytes(a, len(src))
	copy(dst, src)
	return dst
}

// AllocSlice allocates a new slice of a generic type.
//
// Usage:
//	var foo []MyStruct
//	foo = AllocSlice[MyStruct](arena, 16)
//
func AllocSlice[T any](a *Arena, cap int) []T {
	if cap == 0 {
		return nil
	}

	var zero T
	elem := int(unsafe.Sizeof(zero))
	size := elem * cap
	ptr := a.alloc(size)

	return unsafe.Slice((*T)(ptr), cap)
}

// CopySlice allocates a new slice and copies items from src into it.
// The slice capacity is len(src).
func CopySlice[T any](a *Arena, src []T) []T {
	dst := AllocSlice[T](a, len(src))
	copy(dst, src)
	return dst
}

// AllocString returns a string copy allocated in the arena.
func AllocString(a *Arena, s string) string {
	if len(s) == 0 {
		return ""
	}

	b := AllocBytes(a, len(s))
	copy(b, s)
	return *(*string)(unsafe.Pointer(&b))
}

// AllocFreeList allocates and returns a free list in the arena.
// The method returns the same free list for the same type.
func AllocFreeList[T any](a *Arena) *FreeList[T] {
	return allocFreeList[T](a)
}

// internal

func init() {
	initBlockClasses()
	initGlobalHeap()
}
