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

// Remove removes the first occurrence of the item from the slice.
func Remove[T comparable](s []T, item T) []T {
	index := IndexOf(s, item)
	if index == -1 {
		return s
	}
	return RemoveAt(s, index)
}

// RemoveAt removes an item at an index, and truncates the slice.
func RemoveAt[T any](s []T, index int) []T {
	copy(s[index:], s[index+1:])

	// clear the last element for GC
	var zero T
	s[len(s)-1] = zero

	// truncate slice
	return s[:len(s)-1]
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
