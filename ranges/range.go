package ranges

import "github.com/epochtimeout/baselibrary/compare"

// Range is a closed range [start:end] which defines boundaries around a continuous span of values.
type Range[V any] struct {
	Start V
	End   V
}

// Contains returns true if a values is inside the range.
func (r Range[V]) Contains(v V, cmp compare.Compare[V]) bool {
	return cmp(v, r.Start) >= 0 && cmp(v, r.End) <= 0
}

// Inside returns if the current range is inside another range.
func (r Range[V]) Inside(r1 Range[V], cmp compare.Compare[V]) bool {
	return cmp(r.Start, r1.Start) >= 0 && cmp(r.End, r1.End) <= 0
}

// Expand expands the current range and returns a new range.
func (r Range[V]) Expand(r1 Range[V], cmp compare.Compare[V]) Range[V] {
	if cmp(r1.Start, r.Start) < 0 {
		r.Start = r1.Start
	}
	if cmp(r1.End, r.End) > 0 {
		r.End = r1.End
	}
	return r
}

// Overlaps returns if the current range overlaps another range.
//
// One range overlaps another when either start/end of the first
// is inside the second, or vice versa.
func (r Range[V]) Overlaps(r1 Range[V], cmp compare.Compare[V]) bool {
	switch {
	case r.Contains(r1.Start, cmp):
		return true
	case r.Contains(r1.End, cmp):
		return true
	}

	switch {
	case r1.Contains(r.Start, cmp):
		return true
	case r1.Contains(r.End, cmp):
		return true
	}

	return false
}

// ExpandBinary expands a binary range, and returns a new range, skips nil values.
func ExpandBinary(r Range[[]byte], r1 Range[[]byte], cmp compare.Compare[[]byte]) Range[[]byte] {
	if r.Start == nil || (r1.Start != nil && cmp(r1.Start, r.Start) < 0) {
		r.Start = r1.Start
	}
	if r.End == nil || (r1.End != nil && cmp(r1.End, r.End) > 0) {
		r.End = r1.End
	}
	return r
}
