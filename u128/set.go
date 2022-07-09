package u128

import "sort"

// Set is a set of u128 values.
type Set map[U128]struct{}

// NewSet returns a new set.
func NewSet(uu ...U128) Set {
	set := make(Set, len(uu))
	for _, u := range uu {
		set.Add(u)
	}
	return set
}

func cloneSet(s0 Set) Set {
	s1 := make(Set, len(s0))
	for k := range s0 {
		s1[k] = struct{}{}
	}
	return s1
}

// Add adds values to the set.
func (s Set) Add(uu ...U128) {
	for _, u := range uu {
		s[u] = struct{}{}
	}
}

// Remove removes values from the set.
func (s Set) Remove(uu ...U128) {
	for _, u := range uu {
		delete(s, u)
	}
}

// Clone returns a set clone.
func (s Set) Clone() Set {
	return cloneSet(s)
}

// Contains returns whether a value is present in the set.
func (s Set) Contains(u U128) bool {
	_, ok := s[u]
	return ok
}

// ToSlice returns a slice with this set values.
func (s Set) ToSlice() []U128 {
	vv := make([]U128, 0, len(s))
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
