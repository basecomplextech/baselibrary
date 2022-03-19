package alloc

import (
	"reflect"
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
func Alloc[T any](a Arena) *T {
	var zero T

	size := int(unsafe.Sizeof(zero))
	ptr := a.Alloc(size)
	return (*T)(ptr)
}

// AllocBytes allocates a byte slice with a `size` capacity in the arena.
func AllocBytes(a Arena, cap int) []byte {
	if cap == 0 {
		return nil
	}

	ptr := a.Alloc(cap)
	return unsafe.Slice((*byte)(ptr), cap)
}

// AllocBytes allocates a byte slice copy in the arena.
func CopyBytes(a Arena, b []byte) []byte {
	b1 := AllocBytes(a, len(b))
	copy(b1, b)
	return b1
}

// AllocSlice allocates a new slice of a generic type.
//
// Usage:
//	var foo []MyStruct
//	foo = AllocSlice[MyStruct](arena, 0, 16)
//
func AllocSlice[T any](a Arena, cap int) (result []T) {
	if cap == 0 {
		return
	}

	var zero T
	elem := int(unsafe.Sizeof(zero))
	size := elem * cap
	ptr := a.Alloc(size)

	// set slice header
	header := (*reflect.SliceHeader)(unsafe.Pointer(&result))
	header.Data = uintptr(ptr)
	header.Len = cap
	header.Cap = cap
	return result
}

// AllocString returns a string copy allocated in the arena.
func AllocString(a Arena, s string) string {
	if len(s) == 0 {
		return ""
	}

	b := AllocBytes(a, len(s))
	copy(b, s)
	return *(*string)(unsafe.Pointer(&b))
}

// internal

func init() {
	initBlockClasses()
	initGlobalHeap()
}
