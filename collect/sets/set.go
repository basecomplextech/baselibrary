// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package sets

// Set is a collection of unique items implemented as a map[T]struct{}.
type Set[T comparable] map[T]struct{}

// New returns a new set.
func New[T comparable](items ...T) Set[T] {
	set := make(Set[T], len(items))
	for _, item := range items {
		set.Add(item)
	}
	return set
}

// Add adds items to the set.
func (s Set[T]) Add(item T) {
	s[item] = struct{}{}
}

// AddMany adds items to the set.
func (s Set[T]) AddMany(items ...T) {
	for _, item := range items {
		s[item] = struct{}{}
	}
}

// Remove removes items from the set.
func (s Set[T]) Remove(item T) {
	delete(s, item)
}

// RemoveMany removes items from the set.
func (s Set[T]) RemoveMany(items ...T) {
	for _, item := range items {
		delete(s, item)
	}
}

// Clear clears the set.
func (s Set[T]) Clear() {
	clear(s)
}

// Contains returns whether a value is present in the set.
func (s Set[T]) Contains(item T) bool {
	_, ok := s[item]
	return ok
}

// Clone returns a set clone.
func (s Set[T]) Clone() Set[T] {
	s1 := make(Set[T], len(s))
	for item := range s {
		s1[item] = struct{}{}
	}
	return s1
}

// Diff returns items that are in this set but not in the other.
func (s Set[T]) Diff(s1 Set[T]) Set[T] {
	result := make(Set[T], len(s))
	for item := range s {
		if _, ok := s1[item]; !ok {
			result.Add(item)
		}
	}
	return result
}

// Intersect returns items that are in both sets.
func (s Set[T]) Intersect(s1 Set[T]) Set[T] {
	result := make(Set[T], len(s))
	for item := range s {
		if _, ok := s1[item]; ok {
			result.Add(item)
		}
	}
	return result
}

// ToSlice returns a slice with this set items.
func (s Set[T]) ToSlice() []T {
	vv := make([]T, 0, len(s))
	for v := range s {
		vv = append(vv, v)
	}
	return vv
}
