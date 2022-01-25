package u256

import "sort"

// Sort sorts uids.
func Sort(uu []U256) {
	sort.Slice(uu, func(i, j int) bool {
		a := uu[i]
		b := uu[j]
		return a.Less(b)
	})
}

// Sorted returns a sorted slice of uids.
func Sorted(uu ...U256) []U256 {
	dst := make([]U256, len(uu))
	copy(dst, uu)
	Sort(dst)
	return dst
}
