// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package immap

import (
	"slices"
	"sort"
	"testing"

	"github.com/basecomplextech/baselibrary/collect/slices2"
	"github.com/basecomplextech/baselibrary/compare"
	"github.com/basecomplextech/baselibrary/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testMap(t tests.T, items ...Item[int, int64]) *btree[int, int64] {
	compare := compare.Int
	m := newBtree[int, int64](compare)

	for _, item := range items {
		m.Set(item.Key, item.Value)
	}
	return m
}

func testUnwrap[K, V any](m Map[K, V]) *btree[K, V] {
	return m.(*btree[K, V])
}

func testItem(v int) Item[int, int64] {
	return Item[int, int64]{
		Key:   v,
		Value: int64(v),
	}
}

func testItems() []Item[int, int64] {
	n := 2 * maxItems * maxItems
	return testItemsN(n)
}

func testItemsN(n int) []Item[int, int64] {
	items := make([]Item[int, int64], 0, n)
	for i := 0; i < n; i++ {
		item := testItem(i)
		items = append(items, item)
	}
	return items
}

func sortItems(items []Item[int, int64]) {
	sort.Slice(items, func(i, j int) bool {
		a, b := items[i].Key, items[j].Key
		return a < b
	})
}

func testInsert(t tests.T, m *btree[int, int64], items ...Item[int, int64]) {
	for _, item := range items {
		m.Set(item.Key, item.Value)
	}
}

// Free

func TestMap_Free__should_free_map(t *testing.T) {
	items := testItems()
	m := testMap(t, items...)
	m.Free()
}

// Length

func TestMap_Length__should_return_item_count(t *testing.T) {
	t.Skip()

	m := testMap(t)
	items := testItems()
	slices2.Shuffle(items)

	for i, item := range items {
		m.Set(item.Key, item.Value)

		n := m.Length()
		assert.Equal(t, int64(i+1), n)
	}
}

// Clone

func TestMap_Clone__should_return_mutable_clone(t *testing.T) {
	items := testItems()
	m := testMap(t, items...)
	m.Freeze()

	btree1 := m.Clone()
	assert.NotSame(t, m, btree1)
}

func TestMap_Clone__should_allow_clone_mutation(t *testing.T) {
	items := testItems()
	middle := len(items) / 2
	itemsLeft := items[:middle]
	itemsRight := items[middle:]

	m := testMap(t, itemsLeft...)
	m.Freeze()

	btree1 := testUnwrap(m.Clone())
	testInsert(t, btree1, itemsRight...)

	items0 := m.items()
	items1 := btree1.items()
	assert.Equal(t, itemsLeft, items0)
	assert.Equal(t, items, items1)
}

func TestMap_Clone__should_retain_root_branch_children_but_not_values(t *testing.T) {
	items := testItems()
	m := testMap(t, items...)

	m.Freeze()
	m.Clone()

	// Check children refs
	root := m.root.(*branchNode[int, int64])
	for _, item := range root.items {
		child := item.node
		require.Equal(t, int64(2), child.refcount())
	}
}

// Freeze

func TestMap_Freeze__should_recursively_freeze_btree(t *testing.T) {
	items := testItems()
	m := testMap(t, items...)
	m.Freeze()

	walk(m.root, func(n node[int, int64]) {
		require.False(t, n.mutable())
	})
}

// Get

func TestMap_Get__should_return_item_value(t *testing.T) {
	m := testMap(t)
	items := testItems()
	slices2.Shuffle(items)

	for _, item := range items {
		m.Set(item.Key, item.Value)

		value, ok := m.Get(item.Key)
		if !ok {
			t.Fatal(ok, item.Key)
		}

		require.True(t, ok)
		require.Equal(t, item.Value, value)
	}
}

// Set

func TestMap_Set__should_insert_items_in_correct_order(t *testing.T) {
	m := testMap(t)
	items := testItems()
	slices2.Shuffle(items)

	for _, item := range items {
		m.Set(item.Key, item.Value)
	}

	sortItems(items)
	items1 := m.items()
	assert.Equal(t, items, items1)
}

// Delete

func TestMap_Delete__should_delete_items(t *testing.T) {
	m := testMap(t)
	items := testItems()
	slices2.Shuffle(items)
	testInsert(t, m, items...)

	for _, item := range items {
		m.Delete(item.Key)
	}

	assert.Equal(t, 0, m.root.length())
}

// Move

func TestMap_Move__should_move_item_to_another_key(t *testing.T) {
	m := testMap(t)
	items := testItems()

	item0 := items[0]
	item1 := items[1]
	key0 := item0.Key
	key1 := item1.Key
	testInsert(t, m, item0)

	// Move item
	ok := m.Move(key0, key1)
	require.True(t, ok)

	// Check item
	value, ok := m.Get(key1)
	assert.True(t, ok)
	assert.Equal(t, item0.Value, value)

	_, ok = m.Get(key0)
	assert.False(t, ok)
}

func TestMap_Move__should_return_false_if_not_found(t *testing.T) {
	m := testMap(t)
	items := testItems()

	key0 := items[0].Key
	key1 := items[1].Key

	ok := m.Move(key0, key1)
	assert.False(t, ok)
}

// Contains

func TestMap_Contains__should_return_true_when_btree_contains_item(t *testing.T) {
	m := testMap(t)
	items := testItems()
	slices2.Shuffle(items)

	for _, item := range items {
		m.Set(item.Key, item.Value)

		ok := m.Contains(item.Key)
		require.True(t, ok)
	}
}

// items

func TestMap_items__should_returns_items_as_slice(t *testing.T) {
	m := testMap(t)
	items := testItems()
	slices2.Shuffle(items)
	testInsert(t, m, items...)

	sortItems(items)
	items1 := m.items()
	assert.Equal(t, items, items1)
}

// keys

func TestMap_keys__should_return_keys_as_slice(t *testing.T) {
	m := testMap(t)
	items := testItems()
	testInsert(t, m, items...)

	keys := make([]int, 0, len(items))
	for _, item := range items {
		keys = append(keys, item.Key)
	}
	slices.Sort(keys)

	keys1 := m.keys()
	assert.Equal(t, keys, keys1)
}

// values

func TestMap_values__should_return_values_as_slice(t *testing.T) {
	m := testMap(t)
	items := testItems()
	slices2.Shuffle(items)
	testInsert(t, m, items...)

	values := make([]int64, 0, len(items))
	for _, item := range items {
		values = append(values, item.Value)
	}
	slices2.SortLess(values, func(a, b int64) bool {
		return a < b
	})

	values1 := m.values()
	assert.Equal(t, values, values1)
}
