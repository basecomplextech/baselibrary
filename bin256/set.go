package bin256

import "sort"

// Set is a set of bin256 values.
type Set map[B256]struct{}

func NewSet(uu ...B256) Set {
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

func (s Set) Add(u B256) {
	s[u] = struct{}{}
}

func (s Set) Clone() Set {
	return cloneSet(s)
}

func (s Set) Contains(u B256) bool {
	_, ok := s[u]
	return ok
}

func (s Set) Intersect(uu ...B256) []B256 {
	result := make([]B256, 0, len(uu))

	for _, u := range uu {
		if s.Contains(u) {
			result = append(result, u)
		}
	}

	return result
}

func (s Set) Delete(u B256) {
	delete(s, u)
}

func (s Set) Slice() []B256 {
	vv := make([]B256, 0, len(s))
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
