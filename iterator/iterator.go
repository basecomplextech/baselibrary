// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package iterator

import (
	"github.com/basecomplextech/baselibrary/ref"
	"github.com/basecomplextech/baselibrary/status"
)

type (
	// Iter iterates over items of type T.
	Iter[T any] interface {
		// Next returns the next item, or false on the end.
		Next() (T, bool)

		// Free frees the iterator.
		Free()
	}

	// IterError iterates over items of type T, and may return an error.
	IterError[T any] interface {
		// Next returns the next item, or false on the end, or an error.
		Next() (T, bool, error)

		// Free frees the iterator.
		Free()
	}

	// IterStatus iterates over items of type T, and may return an error status.
	IterStatus[T any] interface {
		// Next returns the next item, or false on the end, or an error status.
		Next() (T, bool, status.Status)

		// Free frees the iterator.
		Free()
	}
)

type (
	// NextFunc yields the next item, or false on the end.
	NextFunc[T any] func() (T, bool)

	// NextFuncError yields the next item, or false on the end, or an error.
	NextFuncError[T any] func() (T, bool, error)

	// NextFuncStatus yields the next item, or false on the end, or an error status.
	NextFuncStatus[T any] func() (T, bool, status.Status)
)

// New

// New returns a new iterator.
func New[T any](freer ref.Freer, next NextFunc[T]) Iter[T] {
	return &iter[T]{next, freer}
}

// NewFree returns a new iterator with a free function.
func NewFree[T any](next NextFunc[T], free func()) Iter[T] {
	freer := ref.FreeFunc(free)
	return &iter[T]{next, freer}
}

// NewNoop returns a new iterator without a freer.
func NewNoop[T any](next NextFunc[T]) Iter[T] {
	return &iter[T]{next, nil}
}

// NewError

// NewError returns a new iterator.
func NewError[T any](freer ref.Freer, next NextFuncError[T]) IterError[T] {
	return &iterError[T]{next, freer}
}

// NewFreeError returns a new iterator with a free function.
func NewFreeError[T any](next NextFuncError[T], free func()) IterError[T] {
	freer := ref.FreeFunc(free)
	return &iterError[T]{next, freer}
}

// NewNoopError returns a new iterator with and without a freer.
func NewNoopError[T any](next NextFuncError[T]) IterError[T] {
	return &iterError[T]{next, nil}
}

// NewStatus

// NewStatus returns a new iterator.
func NewStatus[T any](freer ref.Freer, next NextFuncStatus[T]) IterStatus[T] {
	return &iterStatus[T]{next, freer}
}

// NewFreeStatus returns a new iterator with a free function.
func NewFreeStatus[T any](next NextFuncStatus[T], free func()) IterStatus[T] {
	freer := ref.FreeFunc(free)
	return &iterStatus[T]{next, freer}
}

// NewNoopStatus returns a new iterator without a freer.
func NewNoopStatus[T any](next NextFuncStatus[T]) IterStatus[T] {
	return &iterStatus[T]{next, nil}
}

// internal

var _ Iter[any] = (*iter[any])(nil)

type iter[T any] struct {
	next  NextFunc[T]
	freer ref.Freer // maybe nil
}

func (it *iter[T]) Next() (T, bool) {
	return it.next()
}

func (it *iter[T]) Free() {
	if it.freer != nil {
		it.freer.Free()
		it.freer = nil
	}
}

// error

var _ IterError[any] = (*iterError[any])(nil)

type iterError[T any] struct {
	next  NextFuncError[T]
	freer ref.Freer // maybe nil
}

func (it *iterError[T]) Next() (T, bool, error) {
	return it.next()
}

func (it *iterError[T]) Free() {
	if it.freer != nil {
		it.freer.Free()
		it.freer = nil
	}
}

// status

var _ IterStatus[any] = (*iterStatus[any])(nil)

type iterStatus[T any] struct {
	next  NextFuncStatus[T]
	freer ref.Freer // maybe nil
}

func (it *iterStatus[T]) Next() (T, bool, status.Status) {
	return it.next()
}

func (it *iterStatus[T]) Free() {
	if it.freer != nil {
		it.freer.Free()
		it.freer = nil
	}
}
