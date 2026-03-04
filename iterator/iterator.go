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

	// IterErr iterates over items of type T, and may return an error.
	IterErr[T any] interface {
		// Next returns the next item, or false on the end, or an error.
		Next() (T, bool, error)

		// Free frees the iterator.
		Free()
	}

	// IterStat iterates over items of type T, and may return an error status.
	IterStat[T any] interface {
		// Next returns the next item, or false on the end, or an error status.
		Next() (T, bool, status.Status)

		// Free frees the iterator.
		Free()
	}
)

type (
	// NextFunc yields the next item, or false on the end.
	NextFunc[T any] func() (T, bool)

	// NextFuncErr yields the next item, or false on the end, or an error.
	NextFuncErr[T any] func() (T, bool, error)

	// NextFuncStat yields the next item, or false on the end, or an error status.
	NextFuncStat[T any] func() (T, bool, status.Status)
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

// NewErr

// NewErr returns a new iterator.
func NewErr[T any](freer ref.Freer, next NextFuncErr[T]) IterErr[T] {
	return &iterErr[T]{next, freer}
}

// NewFreeErr returns a new iterator with a free function.
func NewFreeErr[T any](next NextFuncErr[T], free func()) IterErr[T] {
	freer := ref.FreeFunc(free)
	return &iterErr[T]{next, freer}
}

// NewNoopErr returns a new iterator with and without a freer.
func NewNoopErr[T any](next NextFuncErr[T]) IterErr[T] {
	return &iterErr[T]{next, nil}
}

// NewStat

// NewStat returns a new iterator.
func NewStat[T any](freer ref.Freer, next NextFuncStat[T]) IterStat[T] {
	return &iterStat[T]{next, freer}
}

// NewFreeStat returns a new iterator with a free function.
func NewFreeStat[T any](next NextFuncStat[T], free func()) IterStat[T] {
	freer := ref.FreeFunc(free)
	return &iterStat[T]{next, freer}
}

// NewNoopStat returns a new iterator without a freer.
func NewNoopStat[T any](next NextFuncStat[T]) IterStat[T] {
	return &iterStat[T]{next, nil}
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

var _ IterErr[any] = (*iterErr[any])(nil)

type iterErr[T any] struct {
	next  NextFuncErr[T]
	freer ref.Freer // maybe nil
}

func (it *iterErr[T]) Next() (T, bool, error) {
	return it.next()
}

func (it *iterErr[T]) Free() {
	if it.freer != nil {
		it.freer.Free()
		it.freer = nil
	}
}

// status

var _ IterStat[any] = (*iterStat[any])(nil)

type iterStat[T any] struct {
	next  NextFuncStat[T]
	freer ref.Freer // maybe nil
}

func (it *iterStat[T]) Next() (T, bool, status.Status) {
	return it.next()
}

func (it *iterStat[T]) Free() {
	if it.freer != nil {
		it.freer.Free()
		it.freer = nil
	}
}
