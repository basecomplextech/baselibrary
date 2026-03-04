// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package iterator

import "github.com/basecomplextech/baselibrary/status"

// Flatten returns an iterator that flattens an iterator of iterators.
// The returned iterator owns the input iterator and the sub-iterators and frees them.
func Flatten[T any](it Iter[Iter[T]]) Iter[T] {
	return &flattenIter[T]{it: it}
}

// FlattenErr returns an iterator that flattens an iterator of iterators.
// The returned iterator owns the input iterator and the sub-iterators and frees them.
func FlattenErr[T any](it IterErr[IterErr[T]]) IterErr[T] {
	return &flattenIterErr[T]{it: it}
}

// FlattenStat returns an iterator that flattens an iterator of iterators.
// The returned iterator owns the input iterator and the sub-iterators and frees them.
func FlattenStat[T any](it IterStat[IterStat[T]]) IterStat[T] {
	return &flattenIterStat[T]{it: it}
}

// internal

var _ Iter[any] = (*flattenIter[any])(nil)

type flattenIter[T any] struct {
	it  Iter[Iter[T]]
	cur Iter[T]
}

// Next returns the next item, or false on an end, or an error.
func (it *flattenIter[T]) Next() (v T, ok bool) {
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
func (it *flattenIter[T]) Free() {
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

var _ IterErr[any] = (*flattenIterErr[any])(nil)

type flattenIterErr[T any] struct {
	it  IterErr[IterErr[T]]
	cur IterErr[T]
}

// Next returns the next item, or false on an end, or an error.
func (it *flattenIterErr[T]) Next() (v T, ok bool, _ error) {
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
func (it *flattenIterErr[T]) Free() {
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

var _ IterStat[any] = (*flattenIterStat[any])(nil)

type flattenIterStat[T any] struct {
	it  IterStat[IterStat[T]]
	cur IterStat[T]
}

// Next returns the next item, or false on an end, or an error.
func (it *flattenIterStat[T]) Next() (v T, ok bool, _ status.Status) {
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
func (it *flattenIterStat[T]) Free() {
	if it.cur != nil {
		it.cur.Free()
		it.cur = nil
	}
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}
