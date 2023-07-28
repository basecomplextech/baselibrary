package ranges

import "github.com/basecomplextech/baselibrary/compare"

// BytesRange is a type alias for Range[[]byte].
type BytesRange = Range[[]byte]

// ExpandBytes expands a bytes range, and returns a new range, skips nil values.
func ExpandBytes(r BytesRange, r1 BytesRange, compare compare.Func[[]byte]) BytesRange {
	if r.Start == nil || (r1.Start != nil && compare(r1.Start, r.Start) < 0) {
		r.Start = r1.Start
	}
	if r.End == nil || (r1.End != nil && compare(r1.End, r.End) > 0) {
		r.End = r1.End
	}
	return r
}
