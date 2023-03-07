package slices

import (
	"math/rand"
	"sort"

	"github.com/complex1tech/baselibrary/constraints"
)

// Clone returns a copy of the slice.
func Clone[T any](s []T) []T {
	s1 := make([]T, len(s))
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
func Contains[T comparable](s []T, item T) bool {
	i := IndexOf(s, item)
	return i >= 0
}

// IndexOf returns the index of the first occurrence of the item in the slice.
// If the item is not found, -1 is returned.
func IndexOf[T comparable](s []T, item T) int {
	for i, v := range s {
		if v == item {
			return i
		}
	}
	return -1
}

// Insert inserts an item at an index.
func Insert[T any](s []T, index int, item T) []T {
	var zero T
	s = append(s, zero)

	copy(s[index+1:], s[index:])
	s[index] = item
	return s
}

// InsertAt inserts items at an index.
func InsertAt[T any](s []T, index int, items ...T) []T {
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
func Remove[T comparable](s []T, item T) []T {
	index := IndexOf(s, item)
	if index == -1 {
		return s
	}
	return RemoveAt(s, index, 1)
}

// RemoveAt removes n items at an index, and truncates the slice.
func RemoveAt[T any](s []T, index int, n int) []T {
	// shift left
	copy(s[index:], s[index+n:])

	// clear n last elements for GC
	for i := len(s) - n; i < len(s); i++ {
		var zero T
		s[i] = zero
	}

	// truncate slice
	return s[:len(s)-n]
}

// Reverse reverse the slice in place.
func Reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// Sort sorts a slice of ordered items.
func Sort[T constraints.Ordered](s []T) {
	sort.Slice(s, func(i, j int) bool {
		a, b := s[i], s[j]
		return a < b
	})
}

// SortCompare sorts a slice.
func SortCompare[T any](s []T, compare func(a, b T) bool) {
	sort.Slice(s, func(i, j int) bool {
		a, b := s[i], s[j]
		return compare(a, b)
	})
}

// Shuffle pseudo-randomizes the order of elements using rand.Shuffle.
func Shuffle[T any](s []T) {
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
}
