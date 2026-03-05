// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package iterator

import "github.com/basecomplextech/baselibrary/status"

// MapToValues converts a map iterator to a value-only iterator.
// The returned iterator owns the input iterator and frees it.
func MapToValues[K, V any](it MapIter[K, V]) Iter[V] {
	return &mapToValues[K, V]{it: it}
}

// MapToValuesError converts a map iterator to a value-only iterator.
// The returned iterator owns the input iterator and frees it.
func MapToValuesError[K, V any](it MapIterErroror[K, V]) IterError[V] {
	return &mapToValuesError[K, V]{it: it}
}

// MapToValuesStatus converts a map iterator to a value-only iterator.
// The returned iterator owns the input iterator and frees it.
func MapToValuesStatus[K, V any](it MapIterStatusus[K, V]) IterStatus[V] {
	return &mapToValuesStatus[K, V]{it: it}
}

// internal

type mapToValues[K, V any] struct {
	it MapIter[K, V]
}

func (it *mapToValues[K, V]) Next() (V, bool) {
	_, v, ok := it.it.Next()
	return v, ok
}

func (it *mapToValues[K, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}

// error

type mapToValuesError[K, V any] struct {
	it MapIterErroror[K, V]
}

func (it *mapToValuesError[K, V]) Next() (V, bool, error) {
	_, v, ok, err := it.it.Next()
	return v, ok, err
}

func (it *mapToValuesError[K, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}

// status

type mapToValuesStatus[K, V any] struct {
	it MapIterStatusus[K, V]
}

func (it *mapToValuesStatus[K, V]) Next() (V, bool, status.Status) {
	_, v, ok, stat := it.it.Next()
	return v, ok, stat
}

func (it *mapToValuesStatus[K, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}
