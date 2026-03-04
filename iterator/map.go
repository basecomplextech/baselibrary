// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package iterator

import "github.com/basecomplextech/baselibrary/status"

type (
	// MapFunc maps an iterator element, or returns false to skip it.
	MapFunc[T, V any] func(T) (V, bool)

	// MapFuncErr maps an iterator element, or returns false to skip it, or an error.
	MapFuncErr[T, V any] func(T) (V, bool, error)

	// MapFuncStat maps an iterator element, or returns false to skip it, or a status.
	MapFuncStat[T, V any] func(T) (V, bool, status.Status)
)

// Map returns an iterator that maps elements from the input iterator.
// The returned iterator owns the input iterator and frees it.
func Map[T any, V any](it Iter[T], fn MapFunc[T, V]) Iter[V] {
	return &mapIter[T, V]{
		it: it,
		fn: fn,
	}
}

// MapErr returns an iterator that maps elements from the input iterator.
// The returned iterator owns the input iterator and frees it.
func MapErr[T any, V any](it IterErr[T], fn MapFuncErr[T, V]) IterErr[V] {
	return &mapIterErr[T, V]{
		it: it,
		fn: fn,
	}
}

// MapStat returns an iterator that maps elements from the input iterator.
// The returned iterator owns the input iterator and frees it.
func MapStat[T any, V any](it IterStat[T], fn MapFuncStat[T, V]) IterStat[V] {
	return &mapIterStat[T, V]{
		it: it,
		fn: fn,
	}
}

// private

var _ Iter[any] = (*mapIter[any, any])(nil)

type mapIter[T any, V any] struct {
	it Iter[T]
	fn MapFunc[T, V]
}

func (it *mapIter[T, V]) Next() (v V, _ bool) {
	for {
		v0, ok := it.it.Next()
		if !ok {
			return v, false
		}

		v, ok := it.fn(v0)
		if !ok {
			continue
		}
		return v, true
	}
}

func (it *mapIter[T, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}

// error

var _ IterErr[any] = (*mapIterErr[any, any])(nil)

type mapIterErr[T any, V any] struct {
	it IterErr[T]
	fn MapFuncErr[T, V]
}

func (it *mapIterErr[T, V]) Next() (v V, _ bool, _ error) {
	for {
		v0, ok, err := it.it.Next()
		if !ok || err != nil {
			return v, false, err
		}

		v, ok, err := it.fn(v0)
		if err != nil {
			return v, false, err
		}
		if !ok {
			continue
		}
		return v, true, nil
	}
}

func (it *mapIterErr[T, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}

// status

var _ IterStat[any] = (*mapIterStat[any, any])(nil)

type mapIterStat[T any, V any] struct {
	it IterStat[T]
	fn MapFuncStat[T, V]
}

func (it *mapIterStat[T, V]) Next() (v V, _ bool, _ status.Status) {
	for {
		v0, ok, st := it.it.Next()
		if !ok || !st.OK() {
			return v, false, st
		}

		v, ok, st := it.fn(v0)
		if !st.OK() {
			return v, false, st
		}
		if !ok {
			continue
		}
		return v, true, status.OK
	}
}

func (it *mapIterStat[T, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}
