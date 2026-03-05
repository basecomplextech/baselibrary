// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package iterator

import "github.com/basecomplextech/baselibrary/status"

type (
	// FilterFunc returns true if the element should be included.
	FilterFunc[T any] func(T) bool

	// FilterFuncError returns true if the element should be included, or an error.
	FilterFuncError[T any] func(T) (bool, error)

	// FilterFuncStatus returns true if the element should be included, and an error status.
	FilterFuncStatus[T any] func(T) (bool, status.Status)
)

// Filter returns an iterator that filters elements from the input iterator.
// The returned iterator owns the input iterator and frees it when done.
func Filter[T any](it Iter[T], fn FilterFunc[T]) Iter[T] {
	return &filter[T]{
		it: it,
		fn: fn,
	}
}

// FilterError returns an iterator that filters elements from the input iterator.
// The returned iterator owns the input iterator and frees it when done.
func FilterError[T any](it IterError[T], fn FilterFuncError[T]) IterError[T] {
	return &filterError[T]{
		it: it,
		fn: fn,
	}
}

// FilterStat returns an iterator that filters elements from the input iterator.
// The returned iterator owns the input iterator and frees it when done.
func FilterStat[T any](it IterStatus[T], fn FilterFuncStatus[T]) IterStatus[T] {
	return &filterStatus[T]{
		it: it,
		fn: fn,
	}
}

// private

var _ Iter[any] = (*filter[any])(nil)

type filter[T any] struct {
	it Iter[T]
	fn FilterFunc[T]
}

func (it *filter[T]) Next() (v T, _ bool) {
	for {
		v, ok := it.it.Next()
		if !ok {
			return v, false
		}

		ok = it.fn(v)
		if !ok {
			continue
		}
		return v, true
	}
}

func (it *filter[T]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}

// error

var _ IterError[any] = (*filterError[any])(nil)

type filterError[T any] struct {
	it IterError[T]
	fn FilterFuncError[T]
}

func (it *filterError[T]) Next() (v T, _ bool, _ error) {
	for {
		v, ok, err := it.it.Next()
		if !ok || err != nil {
			return v, false, err
		}

		ok, err = it.fn(v)
		if err != nil {
			return v, false, err
		}
		if !ok {
			continue
		}
		return v, true, nil
	}
}

func (it *filterError[T]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}

// status

var _ IterStatus[any] = (*filterStatus[any])(nil)

type filterStatus[T any] struct {
	it IterStatus[T]
	fn FilterFuncStatus[T]
}

func (it *filterStatus[T]) Next() (v T, _ bool, _ status.Status) {
	for {
		v, ok, st := it.it.Next()
		if !ok || !st.OK() {
			return v, false, st
		}

		ok, st = it.fn(v)
		if !st.OK() {
			return v, false, st
		}
		if !ok {
			continue
		}
		return v, true, status.OK
	}
}

func (it *filterStatus[T]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}
