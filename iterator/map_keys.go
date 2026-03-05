// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package iterator

import "github.com/basecomplextech/baselibrary/status"

// MapToKeys converts a map iterator to a key-only iterator.
// The returned iterator owns the input iterator and frees it.
func MapToKeys[K, V any](it MapIter[K, V]) Iter[K] {
	return &mapToKeys[K, V]{it: it}
}

// MapToKeysError converts a map iterator to a key-only iterator.
// The returned iterator owns the input iterator and frees it.
func MapToKeysError[K, V any](it MapIterErroror[K, V]) IterError[K] {
	return &mapToKeysError[K, V]{it: it}
}

// MapToKeysStatus converts a map iterator to a key-only iterator.
// The returned iterator owns the input iterator and frees it.
func MapToKeysStatus[K, V any](it MapIterStatusus[K, V]) IterStatus[K] {
	return &mapToKeysStatus[K, V]{it: it}
}

// internal

type mapToKeys[K, V any] struct {
	it MapIter[K, V]
}

func (it *mapToKeys[K, V]) Next() (K, bool) {
	k, _, ok := it.it.Next()
	return k, ok
}

func (it *mapToKeys[K, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}

// error

type mapToKeysError[K, V any] struct {
	it MapIterErroror[K, V]
}

func (it *mapToKeysError[K, V]) Next() (K, bool, error) {
	k, _, ok, err := it.it.Next()
	return k, ok, err
}

func (it *mapToKeysError[K, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}

// status

type mapToKeysStatus[K, V any] struct {
	it MapIterStatusus[K, V]
}

func (it *mapToKeysStatus[K, V]) Next() (K, bool, status.Status) {
	k, _, ok, stat := it.it.Next()
	return k, ok, stat
}

func (it *mapToKeysStatus[K, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}
