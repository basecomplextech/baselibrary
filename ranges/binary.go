package ranges

import "github.com/epochtimeout/baselibrary/compare"

// BinaryRange is a type alias for a range of bytes.
type BinaryRange = Range[[]byte]

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
