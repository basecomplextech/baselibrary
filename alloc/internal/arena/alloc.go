// Copyright 2023 Ivan Korobkov. All rights reserved.

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

	ptr := a.Alloc(size)
	return (*T)(ptr)
}
