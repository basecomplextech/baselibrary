// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package iterator

import "github.com/basecomplextech/baselibrary/status"

// MapToValues converts a map iterator to a value-only iterator.
// The returned iterator owns the input iterator and frees it.
func MapToValues[K, V any](it MapIter[K, V]) Iter[V] {
	return &mapToValuesIter[K, V]{it: it}
}

// MapToValuesErr converts a map iterator to a value-only iterator.
// The returned iterator owns the input iterator and frees it.
func MapToValuesErr[K, V any](it MapIterErr[K, V]) IterErr[V] {
	return &mapToValuesErrIter[K, V]{it: it}
}

// MapToValuesStat converts a map iterator to a value-only iterator.
// The returned iterator owns the input iterator and frees it.
func MapToValuesStat[K, V any](it MapIterStat[K, V]) IterStat[V] {
	return &mapToValuesStatIter[K, V]{it: it}
}

// internal

type mapToValuesIter[K, V any] struct {
	it MapIter[K, V]
}

func (it *mapToValuesIter[K, V]) Next() (V, bool) {
	_, v, ok := it.it.Next()
	return v, ok
}

func (it *mapToValuesIter[K, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}

// error

type mapToValuesErrIter[K, V any] struct {
	it MapIterErr[K, V]
}

func (it *mapToValuesErrIter[K, V]) Next() (V, bool, error) {
	_, v, ok, err := it.it.Next()
	return v, ok, err
}

func (it *mapToValuesErrIter[K, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}

// status

type mapToValuesStatIter[K, V any] struct {
	it MapIterStat[K, V]
}

func (it *mapToValuesStatIter[K, V]) Next() (V, bool, status.Status) {
	_, v, ok, stat := it.it.Next()
	return v, ok, stat
}

func (it *mapToValuesStatIter[K, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}
