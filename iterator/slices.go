// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package iterator

import "github.com/basecomplextech/baselibrary/status"

// ToSlice converts an iterator to a slice.
func ToSlice[T any](it Iter[T]) []T {
	var items []T
	for {
		item, ok := it.Next()
		if !ok {
			break
		}

		items = append(items, item)
	}
	return items
}

// ToSliceErr converts an iterator to a slice.
func ToSliceErr[T any](it IterError[T]) ([]T, error) {
	var items []T
	for {
		item, ok, err := it.Next()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}

		items = append(items, item)
	}
	return items, nil
}

// ToSliceStatus converts an iterator to a slice.
func ToSliceStatus[T any](it IterStatus[T]) ([]T, status.Status) {
	var items []T
	for {
		item, ok, st := it.Next()
		if !st.OK() {
			return nil, st
		}
		if !ok {
			break
		}

		items = append(items, item)
	}
	return items, status.OK
}
