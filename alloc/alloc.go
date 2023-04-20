package alloc

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

	ptr := a.Alloc(size)
	return (*T)(ptr)
}

// Bytes

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

// Slice

// Append appends a new item to a slice, grows the slice if required, and returns the modified slice.
func Append[T any](a Arena, s []T, item T) []T {
	dst := growSlice(a, s, len(s)+1)
	dst = dst[:len(s)+1]
	dst[len(s)] = item
	return dst
}

// Copy allocates a new slice and copies items from src into it.
// The slice capacity is len(src).
func Copy[T any](a Arena, src []T) []T {
	dst := Slice[T](a, len(src))
	copy(dst, src)
	return dst
}

// Slice allocates a new slice of a generic type.
//
// Usage:
//
//	var foo []MyStruct
//	foo = Slice[MyStruct](arena, 16)
func Slice[T any](a Arena, len int) []T {
	return allocSlice[T](a, len, len)
}

// String

// String allocates a new string and copies data from src into it.
func String(a Arena, src string) string {
	if len(src) == 0 {
		return ""
	}

	dst := Bytes(a, len(src))
	copy(dst, src)
	return *(*string)(unsafe.Pointer(&dst))
}

// StringBytes allocates a new string and copies data from src into it.
func StringBytes(a Arena, src []byte) string {
	if len(src) == 0 {
		return ""
	}

	dst := Bytes(a, len(src))
	copy(dst, src)
	return *(*string)(unsafe.Pointer(&dst))
}

// private

func allocSlice[T any](a Arena, len int, cap int) []T {
	if len > cap {
		panic("len > cap")
	}
	if cap == 0 {
		return nil // TODO: Maybe return a zero-length slice
	}

	var zero T
	elem := int(unsafe.Sizeof(zero))
	size := elem * cap

	ptr := a.Alloc(size)
	s := unsafe.Slice((*T)(ptr), cap)
	return s[:len]
}

func growSlice[T any](a Arena, src []T, capacity int) []T {
	if cap(src) >= capacity {
		return src
	}

	oldCap := cap(src)
	newCap := growCapacity(oldCap, capacity)

	dst := allocSlice[T](a, oldCap, newCap)
	copy(dst, src)
	return dst
}

func growCapacity(oldCap int, capacity int) int {
	var newCap int

	if oldCap < 1024 {
		newCap = oldCap + oldCap
	} else {
		// Detect overflow and prevent an infinite loop.
		for 0 < newCap && newCap < capacity {
			newCap += newCap / 4
		}

		// Handle overflow.
		if newCap <= 0 {
			newCap = capacity
		}
	}

	if capacity > newCap {
		newCap = capacity
	}
	return newCap
}
