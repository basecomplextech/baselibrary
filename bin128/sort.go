package bin128

import "sort"

// Sort sorts values.
func Sort(vv []B128) {
	sort.Slice(vv, func(i, j int) bool {
		a := vv[i]
		b := vv[j]
		return a.Less(b)
	})
}

// Sorted returns a sorted slice of values.
func Sorted(vv ...B128) []B128 {
	dst := make([]B128, len(vv))
	copy(dst, vv)
	Sort(dst)
	return dst
}
