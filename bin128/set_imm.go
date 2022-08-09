package bin128

// ImmutableSet is an immutable set of bin128 values.
type ImmutableSet struct {
	items Set
}

// NewImmutableSet returns a new immutable set.
func NewImmutableSet(uu ...B128) *ImmutableSet {
	items := NewSet(uu...)
	return &ImmutableSet{items}
}

// Add returns a set clone with added items.
func (s *ImmutableSet) Add(uu ...B128) *ImmutableSet {
	items := s.items.Clone()
	items.Add(uu...)

	return &ImmutableSet{items}
}

// Remove returns a set clone without removed items.
func (s *ImmutableSet) Remove(uu ...B128) *ImmutableSet {
	items := s.items.Clone()
	items.Remove(uu...)

	return &ImmutableSet{items}
}

// Contains returns whether a value is present in the set.
func (s *ImmutableSet) Contains(u B128) bool {
	return s.items.Contains(u)
}

// Len returns the number of values in this set.
func (s *ImmutableSet) Len() int {
	return len(s.items)
}

// ToSet returns a mutable set with this set values.
func (s *ImmutableSet) ToSet() Set {
	items := s.items.ToSlice()
	return NewSet(items...)
}

// ToSlice returns a slice with this set values.
func (s *ImmutableSet) ToSlice() []B128 {
	return s.items.ToSlice()
}
