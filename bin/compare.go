// Copyright 2024 Ivan Korobkov. All rights reserved.

package bin

// Compare64 compares two bin64 values.
//
// Returns:
//
//	-1 if a < b
//	 0 if a == b
//	 1 if a > b
func Compare64(a, b Bin64) int {
	return a.Compare(b)
}

// Compare128 compares two bin128 values.
//
// Returns:
//
//	-1 if a < b
//	 0 if a == b
//	 1 if a > b
func Compare128(a, b Bin128) int {
	return a.Compare(b)
}

// Compare256 compares two bin256 values.
//
// Returns:
//
//	-1 if a < b
//	 0 if a == b
//	 1 if a > b
func Compare256(a, b Bin256) int {
	return a.Compare(b)
}
