// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package iterator

import "github.com/basecomplextech/baselibrary/status"

type (
	// MapFunc maps an iterator element, or returns false to skip it.
	MapFunc[T, V any] func(T) (V, bool)

	// MapFuncError maps an iterator element, or returns false to skip it, or an error.
	MapFuncError[T, V any] func(T) (V, bool, error)

	// MapFuncStatus maps an iterator element, or returns false to skip it, or a status.
	MapFuncStatus[T, V any] func(T) (V, bool, status.Status)
)

// Map returns an iterator that maps elements from the input iterator.
// The returned iterator owns the input iterator and frees it.
func Map[T any, V any](it Iter[T], fn MapFunc[T, V]) Iter[V] {
	return &mapImpl[T, V]{
		it: it,
		fn: fn,
	}
}

// MapErr returns an iterator that maps elements from the input iterator.
// The returned iterator owns the input iterator and frees it.
func MapErr[T any, V any](it IterError[T], fn MapFuncError[T, V]) IterError[V] {
	return &mapError[T, V]{
		it: it,
		fn: fn,
	}
}

// MapStat returns an iterator that maps elements from the input iterator.
// The returned iterator owns the input iterator and frees it.
func MapStat[T any, V any](it IterStatus[T], fn MapFuncStatus[T, V]) IterStatus[V] {
	return &mapStatus[T, V]{
		it: it,
		fn: fn,
	}
}

// private

var _ Iter[any] = (*mapImpl[any, any])(nil)

type mapImpl[T any, V any] struct {
	it Iter[T]
	fn MapFunc[T, V]
}

func (it *mapImpl[T, V]) Next() (v V, _ bool) {
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

func (it *mapImpl[T, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}

// error

var _ IterError[any] = (*mapError[any, any])(nil)

type mapError[T any, V any] struct {
	it IterError[T]
	fn MapFuncError[T, V]
}

func (it *mapError[T, V]) Next() (v V, _ bool, _ error) {
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

func (it *mapError[T, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}

// status

var _ IterStatus[any] = (*mapStatus[any, any])(nil)

type mapStatus[T any, V any] struct {
	it IterStatus[T]
	fn MapFuncStatus[T, V]
}

func (it *mapStatus[T, V]) Next() (v V, _ bool, _ status.Status) {
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

func (it *mapStatus[T, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}
