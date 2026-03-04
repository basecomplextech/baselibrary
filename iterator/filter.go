// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package iterator

import "github.com/basecomplextech/baselibrary/status"

type (
	// FilterFunc returns true if the element should be included.
	FilterFunc[T any] func(T) bool

	// FilterFuncErr returns true if the element should be included, or an error.
	FilterFuncErr[T any] func(T) (bool, error)

	// FilterFuncStat returns true if the element should be included, and an error status.
	FilterFuncStat[T any] func(T) (bool, status.Status)
)

// Filter returns an iterator that filters elements from the input iterator.
// The returned iterator owns the input iterator and frees it when done.
func Filter[T any](it Iter[T], fn FilterFunc[T]) Iter[T] {
	return &filterIter[T]{
		it: it,
		fn: fn,
	}
}

// FilterErr returns an iterator that filters elements from the input iterator.
// The returned iterator owns the input iterator and frees it when done.
func FilterErr[T any](it IterErr[T], fn FilterFuncErr[T]) IterErr[T] {
	return &filterIterErr[T]{
		it: it,
		fn: fn,
	}
}

// FilterStat returns an iterator that filters elements from the input iterator.
// The returned iterator owns the input iterator and frees it when done.
func FilterStat[T any](it IterStat[T], fn FilterFuncStat[T]) IterStat[T] {
	return &filterIterStat[T]{
		it: it,
		fn: fn,
	}
}

// private

var _ Iter[any] = (*filterIter[any])(nil)

type filterIter[T any] struct {
	it Iter[T]
	fn FilterFunc[T]
}

func (it *filterIter[T]) Next() (v T, _ bool) {
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

func (it *filterIter[T]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}

// error

var _ IterErr[any] = (*filterIterErr[any])(nil)

type filterIterErr[T any] struct {
	it IterErr[T]
	fn FilterFuncErr[T]
}

func (it *filterIterErr[T]) Next() (v T, _ bool, _ error) {
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

func (it *filterIterErr[T]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}

// status

var _ IterStat[any] = (*filterIterStat[any])(nil)

type filterIterStat[T any] struct {
	it IterStat[T]
	fn FilterFuncStat[T]
}

func (it *filterIterStat[T]) Next() (v T, _ bool, _ status.Status) {
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

func (it *filterIterStat[T]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}
