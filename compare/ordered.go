// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package compare

import (
	"github.com/basecomplextech/baselibrary/status"
	"golang.org/x/exp/constraints"
)

// Ordered returns a comparison function for a natually ordered type.
func Ordered[T constraints.Ordered]() Compare[T] {
	return func(a, b T) int {
		switch {
		case a < b:
			return -1
		case a > b:
			return 1
		}
		return 0
	}
}

// OrderedError returns a comparison function for a natually ordered type.
func OrderedError[T constraints.Ordered]() CompareError[T] {
	return func(a, b T) (int, error) {
		switch {
		case a < b:
			return -1, nil
		case a > b:
			return 1, nil
		}
		return 0, nil
	}
}

// OrderedStatus returns a comparison function for a natually ordered type.
func OrderedStatus[T constraints.Ordered]() CompareStatus[T] {
	return func(a, b T) (int, status.Status) {
		switch {
		case a < b:
			return -1, status.OK
		case a > b:
			return 1, status.OK
		}
		return 0, status.OK
	}
}
