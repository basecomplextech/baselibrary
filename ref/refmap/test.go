package refmap

import (
	"github.com/complex1tech/baselibrary/ref"
	"github.com/complex1tech/baselibrary/tests"
)

func testIterate[K any, V ref.Ref](t tests.T, it Iterator[K, V]) []Item[K, V] {
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

func testIterateN[K any, V ref.Ref](t tests.T, it Iterator[K, V], n int) []Item[K, V] {
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

func testIterateBackward[K any, V ref.Ref](t tests.T, it Iterator[K, V]) []Item[K, V] {
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

func testIterateBackwardN[K any, V ref.Ref](t tests.T, it Iterator[K, V], n int) []Item[K, V] {
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
