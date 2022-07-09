package u128

import "sort"

// Sort sorts uids.
func Sort(uu []U128) {
	sort.Slice(uu, func(i, j int) bool {
		a := uu[i]
		b := uu[j]
		return a.Less(b)
	})
}

// Sorted returns a sorted slice of uids.
func Sorted(uu ...U128) []U128 {
	dst := make([]U128, len(uu))
	copy(dst, uu)
	Sort(dst)
	return dst
}
