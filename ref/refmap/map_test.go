// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package refmap

import (
	"slices"
	"sort"
	"testing"

	"github.com/basecomplextech/baselibrary/collect/slices2"
	"github.com/basecomplextech/baselibrary/ref"
	"github.com/basecomplextech/baselibrary/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testBtree(t tests.T, items ...Item[int, *Value]) *btree[int, *Value] {
	compare := func(a, b int) int { return a - b }
	btree := newBtree[int, *Value](compare)

	for _, item := range items {
		btree.SetRetain(item.Key, item.Value)
	}
	return btree
}

func testUnwrap[K, V any](m Map[K, V]) *btree[K, V] {
	return m.(*btree[K, V])
}

func testItem(v int) Item[int, *Value] {
	return Item[int, *Value]{
		Key:   v,
		Value: testValue(v),
	}
}

func testItems() []Item[int, *Value] {
	n := 2 * maxItems * maxItems
	return testItemsN(n)
}

func testItemsN(n int) []Item[int, *Value] {
	items := make([]Item[int, *Value], 0, n)
	for i := 0; i < n; i++ {
		item := testItem(i)
		items = append(items, item)
	}
	return items
}

func sortItems(items []Item[int, *Value]) {
	sort.Slice(items, func(i, j int) bool {
		a, b := items[i].Key, items[j].Key
		return a < b
	})
}

func testInsert(t tests.T, btree *btree[int, *Value], items ...Item[int, *Value]) {
	for _, item := range items {
		btree.SetRetain(item.Key, item.Value)
	}
}

// New

func TestNew__should_retain_values(t *testing.T) {
	items := testItems()
	testBtree(t, items...)

	for _, item := range items {
		v := item.Value
		require.Equal(t, int64(2), v.Refcount())
	}
}

// Free

func TestMap_Free__should_release_values(t *testing.T) {
	items := testItems()
	btree := testBtree(t, items...)
	btree.Free()

	for _, item := range items {
		v := item.Value
		require.Equal(t, int64(1), v.Refcount())
	}
}

// Length

func TestMap_Length__should_return_item_count(t *testing.T) {
	t.Skip()

	btree := testBtree(t)
	items := testItems()
	slices2.Shuffle(items)

	for i, item := range items {
		btree.SetRetain(item.Key, item.Value)

		n := btree.Length()
		assert.Equal(t, int64(i+1), n)
	}
}

// Clone

func TestMap_Clone__should_return_mutable_clone(t *testing.T) {
	items := testItems()
	btree := testBtree(t, items...)
	btree.Freeze()

	btree1 := btree.Clone()
	assert.NotSame(t, btree, btree1)
}

func TestMap_Clone__should_allow_clone_mutation(t *testing.T) {
	items := testItems()
	middle := len(items) / 2
	itemsLeft := items[:middle]
	itemsRight := items[middle:]

	btree := testBtree(t, itemsLeft...)
	btree.Freeze()

	btree1 := testUnwrap(btree.Clone())
	testInsert(t, btree1, itemsRight...)

	items0 := btree.items()
	items1 := btree1.items()
	assert.Equal(t, itemsLeft, items0)
	assert.Equal(t, items, items1)
}

func TestMap_Clone__should_retain_root_leaf_values(t *testing.T) {
	items := testItemsN(maxItems)
	btree := testBtree(t, items...)

	for _, item := range items {
		v := item.Value
		require.Equal(t, int64(2), v.Refcount())
	}

	btree.Freeze()
	btree.Clone()

	for _, item := range items {
		v := item.Value
		require.Equal(t, int64(3), v.Refcount())
	}
}

func TestMap_Clone__should_retain_root_branch_children_but_not_values(t *testing.T) {
	items := testItems()
	btree := testBtree(t, items...)

	btree.Freeze()
	btree.Clone()

	// Check values refs
	for _, item := range items {
		v := item.Value
		require.Equal(t, int64(2), v.Refcount())
	}

	// Check children refs
	root := btree.root.(*branchNode[int, *Value])
	for _, item := range root.items {
		child := item.node
		require.Equal(t, int64(2), child.refcount())
	}
}

// Freeze

func TestMap_Freeze__should_recursively_freeze_btree(t *testing.T) {
	items := testItems()
	btree := testBtree(t, items...)
	btree.Freeze()

	walk(btree.root, func(n node[int, *Value]) {
		require.False(t, n.mutable())
	})
}

// Get

func TestMap_Get__should_return_item_value(t *testing.T) {
	btree := testBtree(t)
	items := testItems()
	slices2.Shuffle(items)

	for _, item := range items {
		btree.SetRetain(item.Key, item.Value)

		value, ok := btree.Get(item.Key)
		if !ok {
			t.Fatal(ok, item.Key)
		}

		require.True(t, ok)
		require.Equal(t, item.Value, value)
	}
}

func TestMap_Get__should_not_retain_value(t *testing.T) {
	items := testItems()
	btree := testBtree(t, items...)

	for _, item := range items {
		value, ok := btree.Get(item.Key)
		if !ok {
			t.Fatal(ok)
		}

		require.True(t, ok)
		require.Equal(t, int64(2), value.Refcount())
	}
}

// Set

func TestMap_Set__should_insert_items_in_correct_order(t *testing.T) {
	btree := testBtree(t)
	items := testItems()
	slices2.Shuffle(items)

	for _, item := range items {
		btree.SetRetain(item.Key, item.Value)
	}

	sortItems(items)
	items1 := btree.items()
	assert.Equal(t, items, items1)
}

func TestMap_Set__should_retain_value(t *testing.T) {
	items := testItems()
	btree := testBtree(t)

	// Set items
	for _, item := range items {
		require.Equal(t, int64(1), item.Value.Refcount())
		btree.SetRetain(item.Key, item.Value)

		require.Equal(t, int64(2), item.Value.Refcount())
	}

	// Check refcounts
	for _, item := range items {
		require.Equal(t, int64(2), item.Value.Refcount())
	}
}

func TestMap_Set__should_retain_release_item_on_replace(t *testing.T) {
	items0 := testItems()
	items1 := testItems()
	btree := testBtree(t, items0...)

	for i, item0 := range items0 {
		item1 := items1[i]

		require.Equal(t, int64(2), item0.Value.Refcount())
		require.Equal(t, int64(1), item1.Value.Refcount())

		btree.SetRetain(item0.Key, item1.Value)

		require.Equal(t, int64(1), item0.Value.Refcount())
		require.Equal(t, int64(2), item1.Value.Refcount())
	}
}

// Delete

func TestMap_Delete__should_delete_items(t *testing.T) {
	btree := testBtree(t)
	items := testItems()
	slices2.Shuffle(items)
	testInsert(t, btree, items...)

	for _, item := range items {
		btree.Delete(item.Key)
	}

	assert.Equal(t, 0, btree.root.length())
}

func TestMap_Delete__should_release_values(t *testing.T) {
	items := testItems()
	btree := testBtree(t, items...)

	for _, item := range items {
		require.Equal(t, int64(2), item.Value.Refcount())
		btree.Delete(item.Key)

		require.Equal(t, int64(1), item.Value.Refcount())
	}
}

// Move

func TestMap_Move__should_move_item_to_another_key(t *testing.T) {
	btree := testBtree(t)
	items := testItems()

	item0 := items[0]
	item1 := items[1]
	key0 := item0.Key
	key1 := item1.Key
	testInsert(t, btree, item0)

	// Move item
	ok := btree.Move(key0, key1)
	require.True(t, ok)

	// Check item
	value, ok := btree.Get(key1)
	assert.True(t, ok)
	assert.Equal(t, item0.Value, value)

	_, ok = btree.Get(key0)
	assert.False(t, ok)
}

func TestMap_Move__should_not_change_refcount(t *testing.T) {
	btree := testBtree(t)
	items := testItems()

	item0 := items[0]
	item1 := items[1]
	key0 := item0.Key
	key1 := item1.Key

	testInsert(t, btree, item0)
	require.Equal(t, int64(2), item0.Value.Refcount())

	// Move item
	ok := btree.Move(key0, key1)
	require.True(t, ok)

	// Check item
	value, ok := btree.Get(key1)
	assert.True(t, ok)
	assert.Equal(t, int64(2), value.Refcount())
}

func TestMap_Move__should_release_overwritten_value(t *testing.T) {
	btree := testBtree(t)
	items := testItems()

	item0 := items[0]
	item1 := items[1]
	key0 := item0.Key
	key1 := item1.Key

	testInsert(t, btree, item0, item1)
	require.Equal(t, int64(2), item0.Value.Refcount())
	require.Equal(t, int64(2), item1.Value.Refcount())

	ok := btree.Move(key0, key1)
	require.True(t, ok)

	require.Equal(t, int64(2), item0.Value.Refcount())
	require.Equal(t, int64(1), item1.Value.Refcount())
}

func TestMap_Move__should_return_false_if_not_found(t *testing.T) {
	btree := testBtree(t)
	items := testItems()

	key0 := items[0].Key
	key1 := items[1].Key

	ok := btree.Move(key0, key1)
	assert.False(t, ok)
}

// Contains

func TestMap_Contains__should_return_true_when_btree_contains_item(t *testing.T) {
	btree := testBtree(t)
	items := testItems()
	slices2.Shuffle(items)

	for _, item := range items {
		btree.SetRetain(item.Key, item.Value)

		ok := btree.Contains(item.Key)
		require.True(t, ok)
	}
}

// items

func TestMap_items__should_returns_items_as_slice(t *testing.T) {
	btree := testBtree(t)
	items := testItems()
	slices2.Shuffle(items)
	testInsert(t, btree, items...)

	sortItems(items)
	items1 := btree.items()
	assert.Equal(t, items, items1)
}

// keys

func TestMap_keys__should_return_keys_as_slice(t *testing.T) {
	btree := testBtree(t)
	items := testItems()
	testInsert(t, btree, items...)

	keys := make([]int, 0, len(items))
	for _, item := range items {
		keys = append(keys, item.Key)
	}
	slices.Sort(keys)

	keys1 := btree.keys()
	assert.Equal(t, keys, keys1)
}

// values

func TestMap_values__should_return_values_as_slice(t *testing.T) {
	btree := testBtree(t)
	items := testItems()
	slices2.Shuffle(items)
	testInsert(t, btree, items...)

	values := make([]ref.R[*Value], 0, len(items))
	for _, item := range items {
		values = append(values, item.Value)
	}
	slices2.SortLess(values, func(a, b ref.R[*Value]) bool {
		return a.Unwrap().val < b.Unwrap().val
	})

	values1 := btree.values()
	assert.Equal(t, values, values1)
}

// test value

type Value struct {
	val   int
	freed bool
}

func testValue(v int) ref.R[*Value] {
	return ref.New(&Value{val: v})
}

func (v *Value) Free() {
	v.freed = true
}
