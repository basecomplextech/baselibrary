package slices

import (
	"math/rand"
	"sort"

	"github.com/epochtimeout/basekit/library/constraints"
)

// Reverse reverse the slice in place.
func Reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
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
