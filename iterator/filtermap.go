// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package iterator

import "github.com/basecomplextech/baselibrary/status"

type (
	// FilterMapFunc filters and maps an element.
	FilterMapFunc[T, V any] func(T) (V, bool)

	// FilterMapFuncError filters and maps an element.
	FilterMapFuncError[T, V any] func(T) (V, bool, error)

	// FilterMapFuncStatus filters and maps an element.
	FilterMapFuncStatus[T, V any] func(T) (V, bool, status.Status)
)

// FilterMap returns an iterator that filters and maps elements from the input iterator.
// The returned iterator owns the input iterator and frees it when done.
func FilterMap[T, V any](it Iter[T], fn FilterMapFunc[T, V]) Iter[V] {
	return &filterMap[T, V]{
		it: it,
		fn: fn,
	}
}

// FilterMapError returns an iterator that filters and maps elements from the input iterator.
// The returned iterator owns the input iterator and frees it when done.
func FilterMapError[T, V any](it IterError[T], fn FilterMapFuncError[T, V]) IterError[V] {
	return &filterMapError[T, V]{
		it: it,
		fn: fn,
	}
}

// FilterMapStatus returns an iterator that filters and maps elements from the input iterator.
// The returned iterator owns the input iterator and frees it when done.
func FilterMapStatus[T, V any](it IterStatus[T], fn FilterMapFuncStatus[T, V]) IterStatus[V] {
	return &filterMapStatus[T, V]{
		it: it,
		fn: fn,
	}
}

// private

var _ Iter[any] = (*filterMap[any, any])(nil)

type filterMap[T, V any] struct {
	it Iter[T]
	fn FilterMapFunc[T, V]
}

func (it *filterMap[T, V]) Next() (v V, _ bool) {
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

func (it *filterMap[T, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}

// error

var _ IterError[any] = (*filterMapError[any, any])(nil)

type filterMapError[T, V any] struct {
	it IterError[T]
	fn FilterMapFuncError[T, V]
}

func (it *filterMapError[T, V]) Next() (v V, _ bool, _ error) {
	for {
		v0, ok, err := it.it.Next()
		if !ok || err != nil {
			return v, false, err
		}

		v, ok, err = it.fn(v0)
		switch {
		case err != nil:
			return v, false, err
		case !ok:
			continue
		}
		return v, true, nil
	}
}

func (it *filterMapError[T, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}

// status

var _ IterStatus[any] = (*filterMapStatus[any, any])(nil)

type filterMapStatus[T, V any] struct {
	it IterStatus[T]
	fn FilterMapFuncStatus[T, V]
}

func (it *filterMapStatus[T, V]) Next() (v V, _ bool, _ status.Status) {
	for {
		v0, ok, st := it.it.Next()
		if !ok || !st.OK() {
			return v, false, st
		}

		v, ok, st = it.fn(v0)
		switch {
		case !st.OK():
			return v, false, st
		case !ok:
			continue
		}
		return v, true, status.OK
	}
}

func (it *filterMapStatus[T, V]) Free() {
	if it.it != nil {
		it.it.Free()
		it.it = nil
	}
}
