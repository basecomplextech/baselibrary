package ranges

import "github.com/complex1tech/baselibrary/compare"

// Binary is a type alias for Range[[]byte].
type Binary = Range[[]byte]

// ExpandBinary expands a binary range, and returns a new range, skips nil values.
func ExpandBinary(r Binary, r1 Binary, compare compare.Func[[]byte]) Binary {
	if r.Start == nil || (r1.Start != nil && compare(r1.Start, r.Start) < 0) {
		r.Start = r1.Start
	}
	if r.End == nil || (r1.End != nil && compare(r1.End, r.End) > 0) {
		r.End = r1.End
	}
	return r
}
