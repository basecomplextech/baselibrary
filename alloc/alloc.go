package alloc

import (
	"unicode/utf8"
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
func Append[S ~[]T, T any](a Arena, s []T, item T) S {
	dst := growSlice[S, T](a, s, len(s)+1)
	dst = dst[:len(s)+1]
	dst[len(s)] = item
	return dst
}

// AppendN appends a new slice to a slice, grows the slice if required, and returns the modified slice.
func AppendN[S ~[]T, T any](a Arena, s []T, items ...T) S {
	dst := growSlice[S, T](a, s, len(s)+len(items))
	dst = dst[:len(s)+len(items)]
	copy(dst[len(s):], items)
	return dst
}

// Copy allocates a new slice and copies items from src into it.
// The slice capacity is len(src).
//
// TODO: Fix me, use S ~[]T
func Copy[T any](a Arena, src []T) []T {
	dst := allocSlice[[]T, T](a, len(src), len(src))
	copy(dst, src)
	return dst
}

// Grow grows the slice to at least the given capacity.
func Grow[S ~[]T, T any](a Arena, s []T, capacity int) S {
	return growSlice[S, T](a, s, capacity)
}

// Slice allocates a new slice of a generic type.
//
// Usage:
//
//	s := Slice[[]MyStruct](arena, 0, 16)
func Slice[S ~[]T, T any](a Arena, len int, cap int) S {
	return allocSlice[S, T](a, len, cap)
}

// Slice1 allocates a new slice with a single item.
//
// Usage:
//
//	elem := 123
//	s := Slice[[]int](arena, elem)
func Slice1[S ~[]T, T any](a Arena, item T) S {
	s := allocSlice[S, T](a, 1, 1)
	s[0] = item
	return s
}

// String

// String allocates a new string and copies data from src into it.
func String(a Arena, src string) string {
	if len(src) == 0 {
		return ""
	}

	dst := Bytes(a, len(src))
	copy(dst, src)
	return unsafeString(dst)
}

// StringBytes allocates a new string and copies data from src into it.
func StringBytes(a Arena, src []byte) string {
	if len(src) == 0 {
		return ""
	}

	dst := Bytes(a, len(src))
	copy(dst, src)
	return unsafeString(dst)
}

// StringRunes allocates a new string and copies data from src into it.
func StringRunes(a Arena, src []rune) string {
	size := 0
	for _, r := range src {
		size += utf8.RuneLen(r)
	}

	dst := allocSlice[[]byte](a, size, size)
	n := 0

	for _, r := range src {
		n += utf8.EncodeRune(dst[n:], r)
	}

	return unsafeString(dst)
}

// StringJoin allocates a new string and joins items from src into it.
func StringJoin(a Arena, src []string, sep string) string {
	size := 0
	for i, s := range src {
		if i > 0 {
			size++
		}
		size += len(s)
	}
	if size == 0 {
		return ""
	}

	b := allocSlice[[]byte](a, 0, size)
	for i, s := range src {
		if i > 0 {
			b = append(b, '.')
		}
		b = append(b, s...)
	}

	return unsafeString(b)
}

// private

func allocSlice[S ~[]T, T any](a Arena, len int, cap int) S {
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

func growSlice[S ~[]T, T any](a Arena, src []T, capacity int) S {
	if cap(src) >= capacity {
		return src
	}

	oldCap := cap(src)
	newCap := growCapacity(oldCap, capacity)

	dst := allocSlice[S, T](a, oldCap, newCap)
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

// unsafe

func unsafeString(b []byte) string {
	if b == nil {
		return ""
	}
	return *(*string)(unsafe.Pointer(&b))
}
