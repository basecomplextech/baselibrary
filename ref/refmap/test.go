// Copyright 2023 Ivan Korobkov. All rights reserved.

package refmap

import (
	"github.com/basecomplextech/baselibrary/tests"
)

func testIterate[K, V any](t tests.T, it Iterator[K, V]) []Item[K, V] {
	result := []Item[K, V]{}

	for it.Next() {
		key := it.Key()
		value := it.Value()

		item := Item[K, V]{
			Key:   key,
			Value: value,
		}
		result = append(result, item)
	}

	return result
}

func testIterateN[K, V any](t tests.T, it Iterator[K, V], n int) []Item[K, V] {
	result := []Item[K, V]{}

	for it.Next() {
		key := it.Key()
		value := it.Value()

		item := Item[K, V]{
			Key:   key,
			Value: value,
		}
		result = append(result, item)

		if len(result) == n {
			break
		}
	}

	return result
}

func testIterateBackward[K, V any](t tests.T, it Iterator[K, V]) []Item[K, V] {
	result := []Item[K, V]{}

	for it.Previous() {
		key := it.Key()
		value := it.Value()

		item := Item[K, V]{
			Key:   key,
			Value: value,
		}
		result = append(result, item)
	}

	return result
}

func testIterateBackwardN[K, V any](t tests.T, it Iterator[K, V], n int) []Item[K, V] {
	result := []Item[K, V]{}

	for it.Previous() {
		key := it.Key()
		value := it.Value()

		item := Item[K, V]{
			Key:   key,
			Value: value,
		}
		result = append(result, item)

		if len(result) == n {
			break
		}
	}

	return result
}
