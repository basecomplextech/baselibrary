// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package mapiter

import (
	"github.com/basecomplextech/baselibrary/iterator"
	"github.com/basecomplextech/baselibrary/status"
)

// Keys converts a map iterator to a key-only iterator.
// The returned iterator owns the input iterator and frees it.
func Keys[K, V any](it Iter[K, V]) iterator.Iter[K] {
	return &keys[K, V]{it: it}
}

// KeysError converts a map iterator to a key-only iterator.
// The returned iterator owns the input iterator and frees it.
func KeysError[K, V any](it IterError[K, V]) iterator.IterError[K] {
	return &keysError[K, V]{it: it}
}

// KeysStatus converts a map iterator to a key-only iterator.
// The returned iterator owns the input iterator and frees it.
func KeysStatus[K, V any](it IterStatus[K, V]) iterator.IterStatus[K] {
	return &keysStatus[K, V]{it: it}
}

// internal

type keys[K, V any] struct {
	it Iter[K, V]
}

func (it *keys[K, V]) Next() (K, bool) {
	k, _, ok := it.it.Next()
	return k, ok
}

func (it *keys[K, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}

// error

type keysError[K, V any] struct {
	it IterError[K, V]
}

func (it *keysError[K, V]) Next() (K, bool, error) {
	k, _, ok, err := it.it.Next()
	return k, ok, err
}

func (it *keysError[K, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}

// status

type keysStatus[K, V any] struct {
	it IterStatus[K, V]
}

func (it *keysStatus[K, V]) Next() (K, bool, status.Status) {
	k, _, ok, stat := it.it.Next()
	return k, ok, stat
}

func (it *keysStatus[K, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}
