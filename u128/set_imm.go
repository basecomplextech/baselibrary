package u128

// ImmutableSet is an immutable set of u128 values.
type ImmutableSet struct {
	items Set
}

// NewImmutableSet returns a new immutable set.
func NewImmutableSet(uu ...U128) *ImmutableSet {
	items := NewSet(uu...)
	return &ImmutableSet{items}
}

// Add returns a set clone with added items.
func (s *ImmutableSet) Add(uu ...U128) *ImmutableSet {
	items := s.items.Clone()
	items.Add(uu...)

	return &ImmutableSet{items}
}

// Remove returns a set clone without removed items.
func (s *ImmutableSet) Remove(uu ...U128) *ImmutableSet {
	items := s.items.Clone()
	items.Remove(uu...)

	return &ImmutableSet{items}
}

// Contains returns whether a value is present in the set.
func (s *ImmutableSet) Contains(u U128) bool {
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
func (s *ImmutableSet) ToSlice() []U128 {
	return s.items.ToSlice()
}
