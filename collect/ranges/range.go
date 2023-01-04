package ranges

import "github.com/complex1tech/baselibrary/compare"

// Range is a closed range [start:end] which defines boundaries around a continuous span of values.
type Range[V any] struct {
	Start V
	End   V
}

// Contains returns true if a values is inside the range.
func (r Range[V]) Contains(v V, compare compare.Func[V]) bool {
	return compare(v, r.Start) >= 0 && compare(v, r.End) <= 0
}

// Inside returns if the current range is inside another range.
func (r Range[V]) Inside(r1 Range[V], compare compare.Func[V]) bool {
	return compare(r.Start, r1.Start) >= 0 && compare(r.End, r1.End) <= 0
}

// Expand expands the current range and returns a new range.
func (r Range[V]) Expand(r1 Range[V], compare compare.Func[V]) Range[V] {
	if compare(r1.Start, r.Start) < 0 {
		r.Start = r1.Start
	}
	if compare(r1.End, r.End) > 0 {
		r.End = r1.End
	}
	return r
}

// Overlaps returns if the current range overlaps another range.
//
// One range overlaps another when either start/end of the first
// is inside the second, or vice versa.
func (r Range[V]) Overlaps(r1 Range[V], compare compare.Func[V]) bool {
	switch {
	case r.Contains(r1.Start, compare):
		return true
	case r.Contains(r1.End, compare):
		return true
	}

	switch {
	case r1.Contains(r.Start, compare):
		return true
	case r1.Contains(r.End, compare):
		return true
	}

	return false
}
