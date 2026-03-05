// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package refmap

import (
	"github.com/basecomplextech/baselibrary/tests"
)

func testIterate[K, V any](_ tests.T, it Iterator[K, V]) []testTuple[K, V] {
	result := []testTuple[K, V]{}

	for {
		key, value, ok := it.Next()
		if !ok {
			break
		}

		tuple := testTuple[K, V]{
			Key:   key,
			Value: value,
		}
		result = append(result, tuple)
	}

	return result
}

func testIterateN[K, V any](_ tests.T, it Iterator[K, V], n int) []testTuple[K, V] {
	result := []testTuple[K, V]{}

	for {
		key, value, ok := it.Next()
		if !ok {
			break
		}

		tuple := testTuple[K, V]{
			Key:   key,
			Value: value,
		}
		result = append(result, tuple)

		if len(result) == n {
			break
		}
	}

	return result
}

func testIterateBackward[K, V any](_ tests.T, it Iterator[K, V]) []testTuple[K, V] {
	result := []testTuple[K, V]{}

	for {
		key, value, ok := it.Previous()
		if !ok {
			break
		}

		tuple := testTuple[K, V]{
			Key:   key,
			Value: value,
		}
		result = append(result, tuple)
	}

	return result
}

func testIterateBackwardN[K, V any](_ tests.T, it Iterator[K, V], n int) []testTuple[K, V] {
	result := []testTuple[K, V]{}

	for {
		key, value, ok := it.Previous()
		if !ok {
			break
		}

		tuple := testTuple[K, V]{
			Key:   key,
			Value: value,
		}
		result = append(result, tuple)

		if len(result) == n {
			break
		}
	}

	return result
}
