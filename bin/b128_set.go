package bin

import "sort"

// Set128 is a set of bin128 values.
type Set128 map[Bin128]struct{}

// NewSet128 returns a new set.
func NewSet128(uu ...Bin128) Set128 {
	set := make(Set128, len(uu))
	for _, u := range uu {
		set.Add(u)
	}
	return set
}

func cloneSet128(s0 Set128) Set128 {
	s1 := make(Set128, len(s0))
	for k := range s0 {
		s1[k] = struct{}{}
	}
	return s1
}

// Add adds values to the set.
func (s Set128) Add(uu ...Bin128) {
	for _, u := range uu {
		s[u] = struct{}{}
	}
}

// Remove removes values from the set.
func (s Set128) Remove(uu ...Bin128) {
	for _, u := range uu {
		delete(s, u)
	}
}

// Clone returns a set clone.
func (s Set128) Clone() Set128 {
	return cloneSet128(s)
}

// Contains returns whether a value is present in the set.
func (s Set128) Contains(u Bin128) bool {
	_, ok := s[u]
	return ok
}

// ToSlice returns a slice with this set values.
func (s Set128) ToSlice() []Bin128 {
	vv := make([]Bin128, 0, len(s))
	for v := range s {
		vv = append(vv, v)
	}
	sort.Slice(vv, func(i, j int) bool {
		a := vv[i]
		b := vv[j]
		return a.Less(b)
	})
	return vv
}

// Immutable set

// ImmSet128 is an immutable set of bin128 values.
type ImmSet128 struct {
	items Set128
}

// NewImmSet128 returns a new immutable set.
func NewImmSet128(uu ...Bin128) *ImmSet128 {
	items := NewSet128(uu...)
	return &ImmSet128{items}
}

// Add returns a set clone with added items.
func (s *ImmSet128) Add(uu ...Bin128) *ImmSet128 {
	items := s.items.Clone()
	items.Add(uu...)

	return &ImmSet128{items}
}

// Remove returns a set clone without removed items.
func (s *ImmSet128) Remove(uu ...Bin128) *ImmSet128 {
	items := s.items.Clone()
	items.Remove(uu...)

	return &ImmSet128{items}
}

// Contains returns whether a value is present in the set.
func (s *ImmSet128) Contains(u Bin128) bool {
	return s.items.Contains(u)
}

// Len returns the number of values in this set.
func (s *ImmSet128) Len() int {
	return len(s.items)
}

// ToSet returns a mutable set with this set values.
func (s *ImmSet128) ToSet() Set128 {
	items := s.items.ToSlice()
	return NewSet128(items...)
}

// ToSlice returns a slice with this set values.
func (s *ImmSet128) ToSlice() []Bin128 {
	return s.items.ToSlice()
}
