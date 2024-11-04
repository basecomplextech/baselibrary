// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package slices2

import (
	"math/rand"
	"slices"
	"sort"
)

// Clear zeros the slice and returns it.
func Clear[S ~[]T, T any](s S) S {
	clear(s)
	return s
}

// Insert inserts an item at an index.
func Insert[S ~[]T, T any](s S, index int, item T) []T {
	var zero T
	s = append(s, zero)

	copy(s[index+1:], s[index:])
	s[index] = item
	return s
}

// LeftShift shifts the slice left by n, and returns the truncated slice.
func LeftShift[S ~[]T, T any](s S, n int) S {
	copy(s, s[n:])

	// Clear n last elements for GC
	var zero T
	for i := len(s) - n; i < len(s); i++ {
		s[i] = zero
	}

	// Truncate slice
	return s[:len(s)-n]
}

// Random returns a random item from the slice, panics if the slice is empty.
// Internally uses math/rand.Intn to generate a random index.
func Random[S ~[]T, T any](s S) T {
	return s[rand.Intn(len(s))]
}

// Remove removes the first occurrence of the item from the slice.
func Remove[S ~[]T, T comparable](s S, item T) S {
	index := slices.Index(s, item)
	if index == -1 {
		return s
	}
	return RemoveAt(s, index, 1)
}

// RemoveAt removes n items at an index, and truncates the slice.
func RemoveAt[S ~[]T, T any](s S, index int, n int) S {
	// Shift left
	copy(s[index:], s[index+n:])

	// Clear n last elements for GC
	for i := len(s) - n; i < len(s); i++ {
		var zero T
		s[i] = zero
	}

	// Truncate slice
	return s[:len(s)-n]
}

// SortLess sorts a slice.
func SortLess[S ~[]T, T any](s S, less func(a, b T) bool) {
	sort.Slice(s, func(i, j int) bool {
		a, b := s[i], s[j]
		return less(a, b)
	})
}

// Shuffle pseudo-randomizes the order of elements using rand.Shuffle.
func Shuffle[S ~[]T, T any](s S) {
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
}

// Truncate clears and truncates the slice, returns nil if the slice is nil.
//
// Usage:
//
//	s := []int{1, 2, 3}
//	s = slices2.Truncate(s)
func Truncate[S ~[]T, T any](s S) S {
	if s == nil {
		return nil
	}

	clear(s)
	return s[:0]
}
