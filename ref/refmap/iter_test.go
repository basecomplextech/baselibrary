// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package refmap

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Next

func TestIterator_Next__should_iterate_items_in_direct_order(t *testing.T) {
	items := makeTestItems()
	btree := testBtree(t, items...)

	it := btree.Iterator()
	defer it.Free()

	tuples := items.tuples()
	tuples1 := testIterate(t, it)

	assert.Equal(t, tuples, tuples1)
}

func TestIterator_Next__should_end_when_empty(t *testing.T) {
	btree := testBtree(t)
	it := btree.Iterator()

	tuples := makeTestItemsN(0).tuples()
	tuples1 := testIterate(t, it)

	assert.Equal(t, tuples, tuples1)
}

func TestIterator_Next_SeekBefore__should_iterate_forward_from_key(t *testing.T) {
	items := makeTestItems()
	btree := testBtree(t, items...)
	it := btree.Iterator()

	for i, item := range items {
		if ok := it.SeekBefore(item.Key); !ok {
			t.Fatal(ok)
		}

		tuples := testIterate(t, it)
		tuples1 := items[i:].tuples()

		assert.Equal(t, tuples1, tuples)
	}
}

func TestIterator_Next_Previous_Next__should_iterate_in_both_directions(t *testing.T) {
	items := makeTestItems()
	btree := testBtree(t, items...)
	it := btree.Iterator()

	tuples1 := testIterate(t, it)
	tuples2 := testIterateBackward(t, it)
	tuples3 := testIterate(t, it)
	slices.Reverse(tuples2)

	tuples := items.tuples()
	assert.Equal(t, tuples, tuples1)
	assert.Equal(t, tuples, tuples2)
	assert.Equal(t, tuples, tuples3)
}

func TestIterator_Next_Previous_Next__should_switch_directions(t *testing.T) {
	items := makeTestItems()
	btree := testBtree(t, items...)
	it := btree.Iterator()

	n := len(items) / 2
	tuples1 := testIterateN(t, it, n)
	tuples2 := testIterateBackward(t, it)
	tuples3 := testIterateN(t, it, n)
	slices.Reverse(tuples2)

	tuples := items.tuples()
	assert.Equal(t, tuples[:n], tuples1)
	assert.Equal(t, tuples[:n-1], tuples2)
	assert.Equal(t, tuples[:n], tuples3)
}

// Previous

func TestIterator_Previous__should_iterate_items_in_reverse_order(t *testing.T) {
	items := makeTestItems()
	btree := testBtree(t, items...)

	it := btree.Iterator()
	it.SeekToEnd()

	tuples := items.tuples()
	tuples1 := testIterateBackward(t, it)

	slices.Reverse(tuples1)
	assert.Equal(t, tuples, tuples1)
}

func TestIterator_Previous__should_end_when_empty(t *testing.T) {
	btree := testBtree(t)

	it := btree.Iterator()
	it.SeekToEnd()

	tuples := makeTestItemsN(0).tuples()
	tuples1 := testIterateBackward(t, it)

	assert.Equal(t, tuples, tuples1)
}

func TestIterator_Previous_SeekBefore__should_iterate_backward_from_key(t *testing.T) {
	items := makeTestItems()
	btree := testBtree(t, items...)
	it := btree.Iterator()

	for i, item := range items {
		if ok := it.SeekBefore(item.Key); !ok {
			t.Fatal(ok)
		}

		tuples := items[:i].tuples()
		tuples1 := testIterateBackward(t, it)

		slices.Reverse(tuples1)
		assert.Equal(t, tuples, tuples1)
	}
}

func TestIterator_Previous_Next_Previous__should_switch_directions(t *testing.T) {
	items := makeTestItems()
	btree := testBtree(t, items...)

	it := btree.Iterator()
	it.SeekToEnd()

	n := len(items) / 2
	tuples1 := testIterateBackwardN(t, it, n)
	tuples2 := testIterate(t, it)
	tuples3 := testIterateBackwardN(t, it, n)
	slices.Reverse(tuples1)
	slices.Reverse(tuples3)

	assert.Equal(t, items[n:].tuples(), tuples1)
	assert.Equal(t, items[n+1:].tuples(), tuples2)
	assert.Equal(t, items[n:].tuples(), tuples3)
}

// SeekBefore

func TestIterator_SeekBefore__should_position_before_key(t *testing.T) {
	items := makeTestItems()
	btree := testBtree(t, items...)
	it := btree.Iterator().(*iter[int, *Value])

	for _, item := range items {
		if ok := it.SeekBefore(item.Key); !ok {
			t.Fatal(ok)
		}

		require.Equal(t, positionBefore, it.pos, "key=%v", item.Key)

		key1, _, ok := it.Next()
		require.True(t, ok, "key=%v", item.Key)
		require.Equal(t, item.Key, key1, "key=%v", item.Key)
	}
}

func TestIterator_SeekBefore__should_position_at_start_when_first_key_greater_than_sought_key(t *testing.T) {
	items := makeTestItemsN(10)
	btree := testBtree(t, items[1:]...)
	it := btree.Iterator().(*iter[int, *Value])

	ok := it.SeekBefore(0)
	assert.True(t, ok)
	assert.Equal(t, positionBefore, it.pos)
}

func TestIterator_SeekBefore__should_position_at_end_when_all_keys_less_than_sought_key(t *testing.T) {
	items := makeTestItemsN(10)
	btree := testBtree(t, items...)
	it := btree.Iterator().(*iter[int, *Value])

	ok := it.SeekBefore(11)
	require.False(t, ok)
	require.Equal(t, positionEnd, it.pos)

	tuples1 := testIterateBackward(t, it)
	slices.Reverse(tuples1)

	tuples := items.tuples()
	assert.Equal(t, tuples, tuples1)
}

// Concurrent modification

func TestIterator_Next__should_panic_on_concurrent_modification(t *testing.T) {
	items := makeTestItemsN(10)
	btree := testBtree(t, items[:5]...)
	it := btree.Iterator().(*iter[int, *Value])

	_, _, ok := it.Next()
	require.True(t, ok)

	btree.SetRetain(items[5].Key, items[5].Value)

	assert.Panics(t, func() {
		it.Next()
	})
}

func TestIterator_Previous__should_panic_on_concurrent_modification(t *testing.T) {
	items := makeTestItemsN(10)
	btree := testBtree(t, items[:5]...)
	it := btree.Iterator().(*iter[int, *Value])

	it.SeekToEnd()
	btree.SetRetain(items[5].Key, items[5].Value)

	assert.Panics(t, func() {
		it.Previous()
	})
}
