package bin

import "sort"

// Set256 is a set of bin256 values.
type Set256 map[Bin256]struct{}

func NewSet256(uu ...Bin256) Set256 {
	set := make(Set256, len(uu))
	for _, u := range uu {
		set.Add(u)
	}
	return set
}

func cloneSet(s0 Set256) Set256 {
	s1 := make(Set256, len(s0))
	for k := range s0 {
		s1[k] = struct{}{}
	}
	return s1
}

func (s Set256) Add(u Bin256) {
	s[u] = struct{}{}
}

func (s Set256) Clone() Set256 {
	return cloneSet(s)
}

func (s Set256) Contains(u Bin256) bool {
	_, ok := s[u]
	return ok
}

func (s Set256) Intersect(uu ...Bin256) []Bin256 {
	result := make([]Bin256, 0, len(uu))

	for _, u := range uu {
		if s.Contains(u) {
			result = append(result, u)
		}
	}

	return result
}

func (s Set256) Delete(u Bin256) {
	delete(s, u)
}

func (s Set256) Slice() []Bin256 {
	vv := make([]Bin256, 0, len(s))
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
