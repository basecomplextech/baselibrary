// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package mapiter

import (
	"github.com/basecomplextech/baselibrary/iterator"
	"github.com/basecomplextech/baselibrary/status"
)

// Values converts a map iterator to a value-only iterator.
// The returned iterator owns the input iterator and frees it.
func Values[K, V any](it Iter[K, V]) iterator.Iter[V] {
	return &values[K, V]{it: it}
}

// ValuesError converts a map iterator to a value-only iterator.
// The returned iterator owns the input iterator and frees it.
func ValuesError[K, V any](it IterError[K, V]) iterator.IterError[V] {
	return &valuesError[K, V]{it: it}
}

// ValuesStatus converts a map iterator to a value-only iterator.
// The returned iterator owns the input iterator and frees it.
func ValuesStatus[K, V any](it IterStatus[K, V]) iterator.IterStatus[V] {
	return &valuesStatus[K, V]{it: it}
}

// internal

type values[K, V any] struct {
	it Iter[K, V]
}

func (it *values[K, V]) Next() (V, bool) {
	_, v, ok := it.it.Next()
	return v, ok
}

func (it *values[K, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}

// error

type valuesError[K, V any] struct {
	it IterError[K, V]
}

func (it *valuesError[K, V]) Next() (V, bool, error) {
	_, v, ok, err := it.it.Next()
	return v, ok, err
}

func (it *valuesError[K, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}

// status

type valuesStatus[K, V any] struct {
	it IterStatus[K, V]
}

func (it *valuesStatus[K, V]) Next() (V, bool, status.Status) {
	_, v, ok, stat := it.it.Next()
	return v, ok, stat
}

func (it *valuesStatus[K, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}
