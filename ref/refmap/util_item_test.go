// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package refmap

import (
	"sort"

	"github.com/basecomplextech/baselibrary/ref"
)

type testItem[K, V any] struct {
	Key   K
	Value ref.R[V]
}

type testItems[K, V any] []testItem[K, V]

type testTuple[K, V any] struct {
	Key   K
	Value V
}

// make

func makeTestItem(v int) testItem[int, *Value] {
	return testItem[int, *Value]{
		Key:   v,
		Value: testValue(v),
	}
}

func makeTestItems() testItems[int, *Value] {
	n := 2 * maxItems * maxItems
	return makeTestItemsN(n)
}

func makeTestItemsN(n int) testItems[int, *Value] {
	items := make(testItems[int, *Value], 0, n)
	for i := 0; i < n; i++ {
		item := makeTestItem(i)
		items = append(items, item)
	}
	return items
}

func sortItems(items testItems[int, *Value]) {
	sort.Slice(items, func(i, j int) bool {
		a, b := items[i].Key, items[j].Key
		return a < b
	})
}

// items

func (m testItems[K, V]) keys() []K {
	result := make([]K, len(m))
	for i, item := range m {
		result[i] = item.Key
	}
	return result
}

func (m testItems[K, V]) values() []V {
	result := make([]V, len(m))
	for i, item := range m {
		result[i] = item.Value.Unwrap()
	}
	return result
}

func (m testItems[K, V]) tuples() []testTuple[K, V] {
	result := make([]testTuple[K, V], len(m))
	for i, item := range m {
		result[i] = testTuple[K, V]{
			Key:   item.Key,
			Value: item.Value.Unwrap(),
		}
	}
	return result
}
