// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package iterator

import "github.com/basecomplextech/baselibrary/status"

// Flatten returns an iterator that flattens an iterator of iterators.
// The returned iterator owns the input iterator and the sub-iterators and frees them.
func Flatten[T any](it Iter[Iter[T]]) Iter[T] {
	return &flatten[T]{it: it}
}

// FlattenError returns an iterator that flattens an iterator of iterators.
// The returned iterator owns the input iterator and the sub-iterators and frees them.
func FlattenError[T any](it IterError[IterError[T]]) IterError[T] {
	return &flattenError[T]{it: it}
}

// FlattenStatus returns an iterator that flattens an iterator of iterators.
// The returned iterator owns the input iterator and the sub-iterators and frees them.
func FlattenStatus[T any](it IterStatus[IterStatus[T]]) IterStatus[T] {
	return &flattenStatus[T]{it: it}
}

// internal

var _ Iter[any] = (*flatten[any])(nil)

type flatten[T any] struct {
	it  Iter[Iter[T]]
	cur Iter[T]
}

// Next returns the next item, or false on an end, or an error.
func (it *flatten[T]) Next() (v T, ok bool) {
	for {
		if it.cur != nil {
			v, ok := it.cur.Next()
			if !ok {
				return v, false
			}
			if ok {
				return v, true
			}

			it.cur.Free()
			it.cur = nil
		}

		next, ok := it.it.Next()
		if !ok {
			return v, false
		}

		it.cur = next
	}
}

// Free frees the iterator.
func (it *flatten[T]) Free() {
	if it.cur != nil {
		it.cur.Free()
		it.cur = nil
	}
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}

// error

var _ IterError[any] = (*flattenError[any])(nil)

type flattenError[T any] struct {
	it  IterError[IterError[T]]
	cur IterError[T]
}

// Next returns the next item, or false on an end, or an error.
func (it *flattenError[T]) Next() (v T, ok bool, _ error) {
	for {
		if it.cur != nil {
			v, ok, err := it.cur.Next()
			if err != nil {
				return v, false, err
			}
			if ok {
				return v, true, nil
			}

			it.cur.Free()
			it.cur = nil
		}

		next, ok, err := it.it.Next()
		if !ok || err != nil {
			return v, false, err
		}

		it.cur = next
	}
}

// Free frees the iterator.
func (it *flattenError[T]) Free() {
	if it.cur != nil {
		it.cur.Free()
		it.cur = nil
	}
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}

// status

var _ IterStatus[any] = (*flattenStatus[any])(nil)

type flattenStatus[T any] struct {
	it  IterStatus[IterStatus[T]]
	cur IterStatus[T]
}

// Next returns the next item, or false on an end, or an error.
func (it *flattenStatus[T]) Next() (v T, ok bool, _ status.Status) {
	for {
		if it.cur != nil {
			v, ok, st := it.cur.Next()
			if !st.OK() {
				return v, false, st
			}
			if ok {
				return v, true, status.OK
			}

			it.cur.Free()
			it.cur = nil
		}

		next, ok, st := it.it.Next()
		if !ok || !st.OK() {
			return v, false, st
		}

		it.cur = next
	}
}

// Free frees the iterator.
func (it *flattenStatus[T]) Free() {
	if it.cur != nil {
		it.cur.Free()
		it.cur = nil
	}
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}
