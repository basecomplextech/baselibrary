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

// Reverse reverse the slice in place.
func Reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// RemoveAt removes an item at an index.
func RemoveAt[T any](s []T, index int) []T {
	copy(s[index:], s[index+1:])
	return s[:len(s)-1]
}

// Shuffle pseudo-randomizes the order of elements using rand.Shuffle.
func Shuffle[T any](s []T) {
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
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
