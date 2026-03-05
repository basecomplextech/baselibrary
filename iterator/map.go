// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package iterator

import "github.com/basecomplextech/baselibrary/status"

type (
	// MapFunc maps an iterator element.
	MapFunc[T, V any] func(T) V

	// MapFuncError maps an iterator element.
	MapFuncError[T, V any] func(T) (V, error)

	// MapFuncStatus maps an iterator element.
	MapFuncStatus[T, V any] func(T) (V, status.Status)
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
	v0, ok := it.it.Next()
	if !ok {
		return v, false
	}

	v = it.fn(v0)
	return v, true
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
	v0, ok, err := it.it.Next()
	if !ok || err != nil {
		return v, false, err
	}

	v, err = it.fn(v0)
	if err != nil {
		return v, false, err
	}
	return v, true, nil
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
	v0, ok, st := it.it.Next()
	if !ok || !st.OK() {
		return v, false, st
	}

	v, st = it.fn(v0)
	if !st.OK() {
		return v, false, st
	}
	return v, true, status.OK
}

func (it *mapStatus[T, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}
