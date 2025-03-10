// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package immap

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Next

func TestIterator_Next__should_iterate_items_in_direct_order(t *testing.T) {
	items := testItems()
	m := testMap(t, items...)

	it := m.Iterator()
	defer it.Free()

	items1 := testIterate(t, it)
	assert.Equal(t, items, items1)
}

func TestIterator_Next__should_end_when_empty(t *testing.T) {
	m := testMap(t)

	it := m.Iterator()
	items0 := testItemsN(0)
	items1 := testIterate(t, it)
	assert.Equal(t, items0, items1)
}

func TestIterator_Next_SeekBefore__should_iterate_forward_from_key(t *testing.T) {
	items := testItems()
	m := testMap(t, items...)
	it := m.Iterator()

	for i, item := range items {
		if ok := it.SeekBefore(item.Key); !ok {
			t.Fatal(ok)
		}

		items1 := testIterate(t, it)
		items2 := items[i:]
		assert.Equal(t, items2, items1)
	}
}

func TestIterator_Next_Previous_Next__should_iterate_in_both_directions(t *testing.T) {
	items := testItems()
	m := testMap(t, items...)
	it := m.Iterator()

	items1 := testIterate(t, it)
	items2 := testIterateBackward(t, it)
	items3 := testIterate(t, it)
	slices.Reverse(items2)

	assert.Equal(t, items, items1)
	assert.Equal(t, items, items2)
	assert.Equal(t, items, items3)
}

func TestIterator_Next_Previous_Next__should_switch_directions(t *testing.T) {
	items := testItems()
	m := testMap(t, items...)
	it := m.Iterator()

	n := len(items) / 2
	items1 := testIterateN(t, it, n)
	items2 := testIterateBackward(t, it)
	items3 := testIterateN(t, it, n)
	slices.Reverse(items2)

	assert.Equal(t, items[:n], items1)
	assert.Equal(t, items[:n-1], items2)
	assert.Equal(t, items[:n], items3)
}

// Previous

func TestIterator_Previous__should_iterate_items_in_reverse_order(t *testing.T) {
	items := testItems()
	m := testMap(t, items...)
	it := m.Iterator()

	if ok := it.SeekToEnd(); !ok {
		t.Fatal(ok)
	}

	items1 := testIterateBackward(t, it)
	slices.Reverse(items)
	assert.Equal(t, items, items1)
}

func TestIterator_Previous__should_end_when_empty(t *testing.T) {
	m := testMap(t)
	it := m.Iterator()

	if ok := it.SeekToEnd(); !ok {
		t.Fatal(ok)
	}

	items0 := testItemsN(0)
	items1 := testIterateBackward(t, it)
	assert.Equal(t, items0, items1)
}

func TestIterator_Previous_SeekBefore__should_iterate_backward_from_key(t *testing.T) {
	items := testItems()
	m := testMap(t, items...)
	it := m.Iterator()

	for i, item := range items {
		if ok := it.SeekBefore(item.Key); !ok {
			t.Fatal(ok)
		}

		items1 := testIterateBackward(t, it)
		items2 := items[:i]

		slices.Reverse(items1)
		assert.Equal(t, items2, items1)
	}
}

func TestIterator_Previous_Next_Previous__should_switch_directions(t *testing.T) {
	items := testItems()
	m := testMap(t, items...)
	it := m.Iterator()

	if ok := it.SeekToEnd(); !ok {
		t.Fatal(ok)
	}

	n := len(items) / 2
	items1 := testIterateBackwardN(t, it, n)
	items2 := testIterate(t, it)
	items3 := testIterateBackwardN(t, it, n)
	slices.Reverse(items1)
	slices.Reverse(items3)

	assert.Equal(t, items[n:], items1)
	assert.Equal(t, items[n+1:], items2)
	assert.Equal(t, items[n:], items3)
}

// SeekBefore

func TestIterator_SeekBefore__should_position_before_key(t *testing.T) {
	items := testItems()
	m := testMap(t, items...)
	it := m.Iterator().(*iterator[int, int64])

	for _, item := range items {
		if ok := it.SeekBefore(item.Key); !ok {
			t.Fatal(ok)
		}

		ok := it.OK()
		assert.False(t, ok)
		assert.Equal(t, positionBefore, it.pos, "key=%v", item.Key)

		ok = it.Next()
		require.True(t, ok, "key=%v", item.Key)

		key1 := it.Key()
		require.Equal(t, item.Key, key1, "key=%v", item.Key)
	}
}

func TestIterator_SeekBefore__should_position_at_start_when_first_key_greater_than_sought_key(t *testing.T) {
	items := testItemsN(10)
	m := testMap(t, items[1:]...)
	it := m.Iterator().(*iterator[int, int64])

	ok := it.SeekBefore(0)
	assert.True(t, ok)

	ok = it.OK()
	assert.False(t, ok)
	assert.Equal(t, positionBefore, it.pos)
}

func TestIterator_SeekBefore__should_position_at_end_when_all_keys_less_than_sought_key(t *testing.T) {
	items := testItemsN(10)
	m := testMap(t, items...)
	it := m.Iterator().(*iterator[int, int64])

	ok := it.SeekBefore(11)
	assert.False(t, ok)

	ok = it.OK()
	assert.False(t, ok)
	assert.Equal(t, positionEnd, it.pos)

	items1 := testIterateBackward[int, int64](t, it)
	slices.Reverse(items1)
	assert.Equal(t, items, items1)
}

// Concurrent modification

func TestIterator_Next__should_panic_on_concurrent_modification(t *testing.T) {
	items := testItemsN(10)
	m := testMap(t, items[:5]...)
	it := m.Iterator().(*iterator[int, int64])

	ok := it.Next()
	assert.True(t, ok)

	m.Set(items[5].Key, items[5].Value)

	assert.Panics(t, func() {
		it.Next()
	})
}

func TestIterator_Previous__should_panic_on_concurrent_modification(t *testing.T) {
	items := testItemsN(10)
	m := testMap(t, items[:5]...)
	it := m.Iterator().(*iterator[int, int64])

	ok := it.SeekToEnd()
	assert.True(t, ok)

	m.Set(items[5].Key, items[5].Value)

	assert.Panics(t, func() {
		it.Previous()
	})
}
