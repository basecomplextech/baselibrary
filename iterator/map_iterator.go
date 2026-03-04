// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package iterator

import (
	"github.com/basecomplextech/baselibrary/ref"
	"github.com/basecomplextech/baselibrary/status"
)

type (
	// MapIter iterates over key-value pairs.
	MapIter[K, V any] interface {
		// Next returns the next key-value pair, or false on the end.
		Next() (K, V, bool)

		// Free frees the iterator.
		Free()
	}

	// MapIterErr iterates over key-value pairs, and may return an error.
	MapIterErr[K, V any] interface {
		// Next returns the next key-value pair, or false on the end, or an error.
		Next() (K, V, bool, error)

		// Free frees the iterator.
		Free()
	}

	// MapIterStat iterates over key-value pairs, and may return an error status.
	MapIterStat[K, V any] interface {
		// Next returns the next key-value pair, or false on the end, or an error status.
		Next() (K, V, bool, status.Status)

		// Free frees the iterator.
		Free()
	}
)

type (
	// MapNextFunc yields the next key-value pair, or false on the end.
	MapNextFunc[K, V any] func() (K, V, bool)

	// MapNextFuncErr yields the next key-value pair, or false on the end, or an error.
	MapNextFuncErr[K, V any] func() (K, V, bool, error)

	// MapNextFuncStat yields the next key-value pair, or false on the end, or an error status.
	MapNextFuncStat[K, V any] func() (K, V, bool, status.Status)
)

// NewMap returns a new map iterator.
func NewMap[K, V any](freer ref.Freer, next MapNextFunc[K, V]) MapIter[K, V] {
	return &mapIter[K, V]{next: next, freer: freer}
}

// NewMapErr returns a new map iterator with an error.
func NewMapErr[K, V any](next MapNextFuncErr[K, V], free func()) MapIterErr[K, V] {
	return &mapIterErr[K, V]{next: next, free: free}
}

// NewMapStat returns a new map iterator with a status.
func NewMapStat[K, V any](next MapNextFuncStat[K, V], free func()) MapIterStat[K, V] {
	return &mapIterStat[K, V]{next: next, free: free}
}

// private

var _ MapIter[any, any] = (*mapIter[any, any])(nil)

type mapIter[K, V any] struct {
	next  MapNextFunc[K, V]
	freer ref.Freer
}

func (it *mapIter[K, V]) Next() (K, V, bool) {
	return it.next()
}

func (it *mapIter[K, V]) Free() {
	if it.freer != nil {
		it.freer.Free()
		it.freer = nil
	}
}

// error

var _ MapIterErr[any, any] = (*mapIterErr[any, any])(nil)

type mapIterErr[K, V any] struct {
	next MapNextFuncErr[K, V]
	free func()
}

func (it *mapIterErr[K, V]) Next() (K, V, bool, error) {
	return it.next()
}

func (it *mapIterErr[K, V]) Free() {
	if it.free != nil {
		it.free()
		it.free = nil
	}
}

// status

var _ MapIterStat[any, any] = (*mapIterStat[any, any])(nil)

type mapIterStat[K, V any] struct {
	next MapNextFuncStat[K, V]
	free func()
}

func (it *mapIterStat[K, V]) Next() (K, V, bool, status.Status) {
	return it.next()
}

func (it *mapIterStat[K, V]) Free() {
	if it.free != nil {
		it.free()
		it.free = nil
	}
}
