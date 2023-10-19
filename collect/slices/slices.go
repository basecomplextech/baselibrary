package slices

import (
	"math/rand"
	"sort"

	"github.com/basecomplextech/baselibrary/constraints"
)

// CastAny casts a slice of any type to a slice of another type.
func CastAny[T any](s []any) []T {
	s1 := make([]T, 0, len(s))
	for _, v := range s {
		v1 := v.(T)
		s1 = append(s1, v1)
	}
	return s1
}

// Clear zeros the slice, truncates it to zero length and returns.
//
// Usage:
//
//	s := []int{1, 2, 3}
//	s = slices.Clear(s)
func Clear[S ~[]T, T any](s S) S {
	var zero T
	for i := range s {
		s[i] = zero
	}
	return s[:0]
}

// Clone returns a copy of the slice.
func Clone[S ~[]T, T any](s S) S {
	s1 := make(S, len(s))
	copy(s1, s)
	return s1
}

// CloneTree returns a deep copy of the slice of slices.
func CloneTree[T any](tree [][]T) [][]T {
	tree1 := make([][]T, len(tree))
	for i, v := range tree {
		tree1[i] = Clone(v)
	}
	return tree1
}

// Contains returns true if the slice contains an item.
func Contains[S ~[]T, T comparable](s S, item T) bool {
	i := IndexOf(s, item)
	return i >= 0
}

// IndexOf returns the index of the first occurrence of the item in the slice.
// If the item is not found, -1 is returned.
func IndexOf[S ~[]T, T comparable](s S, item T) int {
	for i, v := range s {
		if v == item {
			return i
		}
	}
	return -1
}

// Insert inserts an item at an index.
func Insert[S ~[]T, T any](s S, index int, item T) []T {
	var zero T
	s = append(s, zero)

	copy(s[index+1:], s[index:])
	s[index] = item
	return s
}

// InsertAt inserts items at an index.
func InsertAt[S ~[]T, T any](s S, index int, items ...T) []T {
	total := len(s) + len(items)

	if cap(s) < total {
		s1 := make([]T, len(s), total)
		copy(s1, s)
		s = s1
	}

	s = s[:total]
	copy(s[index+len(items):], s[index:])
	copy(s[index:], items)
	return s
}

// Remove removes the first occurrence of the item from the slice.
func Remove[S ~[]T, T comparable](s S, item T) S {
	index := IndexOf(s, item)
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

// Reverse reverse the slice in place.
func Reverse[S ~[]T, T any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// Sort sorts a slice of ordered items.
func Sort[S ~[]T, T constraints.Ordered](s S) {
	sort.Slice(s, func(i, j int) bool {
		a, b := s[i], s[j]
		return a < b
	})
}

// SortCompare sorts a slice.
func SortCompare[S ~[]T, T any](s S, compare func(a, b T) bool) {
	sort.Slice(s, func(i, j int) bool {
		a, b := s[i], s[j]
		return compare(a, b)
	})
}

// Shuffle pseudo-randomizes the order of elements using rand.Shuffle.
func Shuffle[S ~[]T, T any](s S) {
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
}

// ToAny converts a slice of any type to a slice of any.
func ToAny[S ~[]T, T any](s S) []any {
	s1 := make([]any, len(s))
	for i, v := range s {
		s1[i] = v
	}
	return s1
}

// Zero zeros the slice and returns it.
func Zero[S ~[]T, T any](s S) S {
	var zero T
	for i := range s {
		s[i] = zero
	}
	return s
}
