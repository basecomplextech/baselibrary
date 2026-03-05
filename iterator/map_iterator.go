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

	// MapIterErroror iterates over key-value pairs, and may return an error.
	MapIterErroror[K, V any] interface {
		// Next returns the next key-value pair, or false on the end, or an error.
		Next() (K, V, bool, error)

		// Free frees the iterator.
		Free()
	}

	// MapIterStatusus iterates over key-value pairs, and may return an error status.
	MapIterStatusus[K, V any] interface {
		// Next returns the next key-value pair, or false on the end, or an error status.
		Next() (K, V, bool, status.Status)

		// Free frees the iterator.
		Free()
	}
)

type (
	// MapNextFunc yields the next key-value pair, or false on the end.
	MapNextFunc[K, V any] func() (K, V, bool)

	// MapNextFuncErroror yields the next key-value pair, or false on the end, or an error.
	MapNextFuncErroror[K, V any] func() (K, V, bool, error)

	// MapNextFuncStatusus yields the next key-value pair, or false on the end, or an error status.
	MapNextFuncStatusus[K, V any] func() (K, V, bool, status.Status)
)

// NewMap returns a new map iterator.
func NewMap[K, V any](freer ref.Freer, next MapNextFunc[K, V]) MapIter[K, V] {
	return &mapIter[K, V]{next: next, freer: freer}
}

// NewMapError returns a new map iterator with an error.
func NewMapError[K, V any](next MapNextFuncErroror[K, V], free func()) MapIterErroror[K, V] {
	return &mapIterErroror[K, V]{next: next, free: free}
}

// NewMapStatus returns a new map iterator with a status.
func NewMapStatus[K, V any](next MapNextFuncStatusus[K, V], free func()) MapIterStatusus[K, V] {
	return &mapIterStatusus[K, V]{next: next, free: free}
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

var _ MapIterErroror[any, any] = (*mapIterErroror[any, any])(nil)

type mapIterErroror[K, V any] struct {
	next MapNextFuncErroror[K, V]
	free func()
}

func (it *mapIterErroror[K, V]) Next() (K, V, bool, error) {
	return it.next()
}

func (it *mapIterErroror[K, V]) Free() {
	if it.free != nil {
		it.free()
		it.free = nil
	}
}

// status

var _ MapIterStatusus[any, any] = (*mapIterStatusus[any, any])(nil)

type mapIterStatusus[K, V any] struct {
	next MapNextFuncStatusus[K, V]
	free func()
}

func (it *mapIterStatusus[K, V]) Next() (K, V, bool, status.Status) {
	return it.next()
}

func (it *mapIterStatusus[K, V]) Free() {
	if it.free != nil {
		it.free()
		it.free = nil
	}
}
