package hashset

// Set is a collection of unique values implemented as a map[T]struct{}.
type Set[T comparable] map[T]struct{}

// New returns a new set.
func New[T comparable](uu ...T) Set[T] {
	set := make(Set[T], len(uu))
	for _, u := range uu {
		set.Add(u)
	}
	return set
}

// Add adds values to the set.
func (s Set[T]) Add(uu ...T) {
	for _, u := range uu {
		s[u] = struct{}{}
	}
}

// Remove removes values from the set.
func (s Set[T]) Remove(uu ...T) {
	for _, u := range uu {
		delete(s, u)
	}
}

// Contains returns whether a value is present in the set.
func (s Set[T]) Contains(u T) bool {
	_, ok := s[u]
	return ok
}

// Clone returns a set clone.
func (s Set[T]) Clone() Set[T] {
	s1 := make(Set[T], len(s))
	for k := range s {
		s1[k] = struct{}{}
	}
	return s1
}

// ToSlice returns a slice with this set values.
func (s Set[T]) ToSlice() []T {
	vv := make([]T, 0, len(s))
	for v := range s {
		vv = append(vv, v)
	}
	return vv
}
