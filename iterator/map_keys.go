// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package iterator

import "github.com/basecomplextech/baselibrary/status"

// MapToKeys converts a map iterator to a key-only iterator.
// The returned iterator owns the input iterator and frees it.
func MapToKeys[K, V any](it MapIter[K, V]) Iter[K] {
	return &mapToKeysIter[K, V]{it: it}
}

// MapToKeysErr converts a map iterator to a key-only iterator.
// The returned iterator owns the input iterator and frees it.
func MapToKeysErr[K, V any](it MapIterErr[K, V]) IterErr[K] {
	return &mapToKeysErrIter[K, V]{it: it}
}

// MapToKeysStat converts a map iterator to a key-only iterator.
// The returned iterator owns the input iterator and frees it.
func MapToKeysStat[K, V any](it MapIterStat[K, V]) IterStat[K] {
	return &mapToKeysStatIter[K, V]{it: it}
}

// internal

type mapToKeysIter[K, V any] struct {
	it MapIter[K, V]
}

func (it *mapToKeysIter[K, V]) Next() (K, bool) {
	k, _, ok := it.it.Next()
	return k, ok
}

func (it *mapToKeysIter[K, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}

// error

type mapToKeysErrIter[K, V any] struct {
	it MapIterErr[K, V]
}

func (it *mapToKeysErrIter[K, V]) Next() (K, bool, error) {
	k, _, ok, err := it.it.Next()
	return k, ok, err
}

func (it *mapToKeysErrIter[K, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}

// status

type mapToKeysStatIter[K, V any] struct {
	it MapIterStat[K, V]
}

func (it *mapToKeysStatIter[K, V]) Next() (K, bool, status.Status) {
	k, _, ok, stat := it.it.Next()
	return k, ok, stat
}

func (it *mapToKeysStatIter[K, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}
