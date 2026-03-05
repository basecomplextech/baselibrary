// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package iterator

// FlatMap returns an iterator that maps elements from the input iterator to sub-iterators
// and flattens them. The returned iterator owns the sub-iterators and frees them when done.
func FlatMap[T any, V any](it Iter[T], fn MapFunc[T, Iter[V]]) Iter[V] {
	it1 := Map(it, fn)
	return Flatten(it1)
}

// FlatMapError returns an iterator that maps elements from the input iterator to sub-iterators
// and flattens them. The returned iterator owns the sub-iterators and frees them when done.
func FlatMapError[T any, V any](it IterError[T], fn MapFuncError[T, IterError[V]]) IterError[V] {
	it1 := MapErr(it, fn)
	return FlattenError(it1)
}

// FlatMap returns an iterator that maps elements from the input iterator to sub-iterators
// and flattens them. The returned iterator owns the sub-iterators and frees them when done.
func FlatMapStat[T any, V any](it IterStatus[T], fn MapFuncStatus[T, IterStatus[V]]) IterStatus[V] {
	it1 := MapStat(it, fn)
	return FlattenStatus(it1)
}
