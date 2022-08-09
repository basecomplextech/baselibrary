package bin256

import "sort"

// Sort sorts values.
func Sort(vv []B256) {
	sort.Slice(vv, func(i, j int) bool {
		a := vv[i]
		b := vv[j]
		return a.Less(b)
	})
}

// Sorted returns a sorted slice of values.
func Sorted(vv ...B256) []B256 {
	dst := make([]B256, len(vv))
	copy(dst, vv)
	Sort(dst)
	return dst
}
