package bin

import "sort"

// Sort128 sorts bin128 values.
func Sort128(vv []Bin128) {
	sort.Slice(vv, func(i, j int) bool {
		a := vv[i]
		b := vv[j]
		return a.Less(b)
	})
}
