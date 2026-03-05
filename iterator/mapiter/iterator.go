// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package mapiter

import (
	"github.com/basecomplextech/baselibrary/ref"
	"github.com/basecomplextech/baselibrary/status"
)

type (
	// Iter iterates over key-value pairs.
	Iter[K, V any] interface {
		// Next returns the next key-value pair, or false on the end.
		Next() (K, V, bool)

		// Free frees the iterator.
		Free()
	}

	// IterError iterates over key-value pairs, and may return an error.
	IterError[K, V any] interface {
		// Next returns the next key-value pair, or false on the end, or an error.
		Next() (K, V, bool, error)

		// Free frees the iterator.
		Free()
	}

	// IterStatus iterates over key-value pairs, and may return an error status.
	IterStatus[K, V any] interface {
		// Next returns the next key-value pair, or false on the end, or an error status.
		Next() (K, V, bool, status.Status)

		// Free frees the iterator.
		Free()
	}
)

type (
	// NextFunc yields the next key-value pair, or false on the end.
	NextFunc[K, V any] func() (K, V, bool)

	// NextFuncError yields the next key-value pair, or false on the end, or an error.
	NextFuncError[K, V any] func() (K, V, bool, error)

	// NextFuncStatus yields the next key-value pair, or false on the end, or an error status.
	NextFuncStatus[K, V any] func() (K, V, bool, status.Status)
)

// New returns a new map iterator.
func New[K, V any](freer ref.Freer, next NextFunc[K, V]) Iter[K, V] {
	return &iter[K, V]{next: next, freer: freer}
}

// NewError returns a new map iterator with an error.
func NewError[K, V any](next NextFuncError[K, V], free func()) IterError[K, V] {
	return &iterError[K, V]{next: next, free: free}
}

// NewStatus returns a new map iterator with a status.
func NewStatus[K, V any](next NextFuncStatus[K, V], free func()) IterStatus[K, V] {
	return &iterStatus[K, V]{next: next, free: free}
}

// private

var _ Iter[any, any] = (*iter[any, any])(nil)

type iter[K, V any] struct {
	next  NextFunc[K, V]
	freer ref.Freer
}

func (it *iter[K, V]) Next() (K, V, bool) {
	return it.next()
}

func (it *iter[K, V]) Free() {
	if it.freer != nil {
		it.freer.Free()
		it.freer = nil
	}
}

// error

var _ IterError[any, any] = (*iterError[any, any])(nil)

type iterError[K, V any] struct {
	next NextFuncError[K, V]
	free func()
}

func (it *iterError[K, V]) Next() (K, V, bool, error) {
	return it.next()
}

func (it *iterError[K, V]) Free() {
	if it.free != nil {
		it.free()
		it.free = nil
	}
}

// status

var _ IterStatus[any, any] = (*iterStatus[any, any])(nil)

type iterStatus[K, V any] struct {
	next NextFuncStatus[K, V]
	free func()
}

func (it *iterStatus[K, V]) Next() (K, V, bool, status.Status) {
	return it.next()
}

func (it *iterStatus[K, V]) Free() {
	if it.free != nil {
		it.free()
		it.free = nil
	}
}
