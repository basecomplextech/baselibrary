// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package alloc

import (
	"unsafe"
)

// Append appends a new item to a slice, grows the slice if required, and returns the modified slice.
func Append[S ~[]T, T any](a Arena, s S, item T) S {
	dst := growSlice(a, s, len(s)+1)
	dst = dst[:len(s)+1]
	dst[len(s)] = item
	return dst
}

// AppendN appends a new slice to a slice, grows the slice if required, and returns the modified slice.
func AppendN[S ~[]T, T any](a Arena, s S, items ...T) S {
	dst := growSlice(a, s, len(s)+len(items))
	dst = dst[:len(s)+len(items)]
	copy(dst[len(s):], items)
	return dst
}

// Copy allocates a new slice and copies items from src into it.
// The slice capacity is len(src).
func Copy[S ~[]T, T any](a Arena, src S) S {
	dst := allocSlice[S](a, len(src), len(src))
	copy(dst, src)
	return dst
}

// Grow grows the slice to at least the given capacity.
func Grow[S ~[]T, T any](a Arena, s S, capacity int) S {
	return growSlice(a, s, capacity)
}

// Slice allocates a new slice of a generic type.
//
// Usage:
//
//	s := Slice[[]MyStruct](arena, 0, 16)
func Slice[S ~[]T, T any](a Arena, len int, cap int) S {
	return allocSlice[S](a, len, cap)
}

// Slice1 allocates a new slice with a single item.
//
// Usage:
//
//	elem := 123
//	s := Slice[[]int](arena, elem)
func Slice1[S ~[]T, T any](a Arena, item T) S {
	s := allocSlice[S](a, 1, 1)
	s[0] = item
	return s
}

// private

func allocSlice[S ~[]T, T any](a Arena, len int, cap int) S {
	if len > cap {
		panic("len > cap")
	}
	if cap == 0 {
		return nil
	}

	var zero T
	elem := int(unsafe.Sizeof(zero))
	size := elem * cap

	ptr := a.Alloc(size)
	s := unsafe.Slice((*T)(ptr), cap)
	return s[:len]
}

func growSlice[S ~[]T, T any](a Arena, src S, capacity int) S {
	if cap(src) >= capacity {
		return src
	}

	oldCap := cap(src)
	newCap := growCapacity(oldCap, capacity)

	dst := allocSlice[S](a, len(src), newCap)
	copy(dst, src)
	return dst
}

func growCapacity(oldCap int, capacity int) int {
	var newCap int

	if oldCap < 1024 {
		// Double until 1024
		newCap = oldCap + oldCap
	} else {
		// Grow by 25% after 1024
		// Detect overflow and prevent an infinite loop.
		newCap = oldCap
		for newCap < capacity {
			newCap += newCap / 4
			if newCap <= 0 {
				newCap = capacity
				break
			}
		}
	}

	if capacity > newCap {
		newCap = capacity
	}
	return newCap
}
