// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package orderedmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOrderedMap__should_create_new_map(t *testing.T) {
	items := []item[int, int]{
		{3, 3},
		{2, 2},
		{1, 1},
	}

	m := newOrderedMap(items...)
	assert.Equal(t, items, m.items)
}

// Get

func TestMap_Get__should_return_value_for_key(t *testing.T) {
	items := []item[int, int]{
		{3, 3},
		{2, 2},
		{1, 1},
	}

	m := newOrderedMap(items...)
	value, ok := m.Get(2)
	assert.True(t, ok)
	assert.Equal(t, 2, value)

	value, ok = m.Get(4)
	assert.False(t, ok)
	assert.Equal(t, 0, value)
}

// IndexOf

func TestMap_IndexOf__should_return_index_of_key(t *testing.T) {
	items := []item[int, int]{
		{3, 3},
		{2, 2},
		{1, 1},
	}

	m := newOrderedMap(items...)
	assert.Equal(t, 0, m.IndexOf(3))
	assert.Equal(t, 1, m.IndexOf(2))
	assert.Equal(t, 2, m.IndexOf(1))
}

// Keys

func TestMap_Keys__should_return_keys(t *testing.T) {
	items := []item[int, int]{
		{3, 3},
		{2, 2},
		{1, 1},
	}

	m := newOrderedMap(items...)
	keys := m.Keys()

	assert.Equal(t, []int{3, 2, 1}, keys)
}

// Values

func TestMap_Values__should_return_values(t *testing.T) {
	items := []item[int, int]{
		{3, 3},
		{2, 2},
		{1, 1},
	}

	m := newOrderedMap(items...)
	values := m.Values()

	assert.Equal(t, []int{3, 2, 1}, values)
}

// Clone

func TestMap_Clone__should_return_clone(t *testing.T) {
	items := []item[int, int]{
		{3, 3},
		{2, 2},
		{1, 1},
	}

	m := newOrderedMap(items...)
	m1 := m.Clone().(*orderedMap[int, int])

	assert.Equal(t, items, m1.items)
}

// Delete

func TestMap_Delete__should_delete_item(t *testing.T) {
	items := []item[int, int]{
		{3, 3},
		{2, 2},
		{1, 1},
	}

	m := newOrderedMap(items...)
	m.Delete(2)

	items1 := []item[int, int]{
		{3, 3},
		{1, 1},
	}

	assert.Equal(t, items1, m.items)
}

// Set

func TestMap_Set__should_add_new_item(t *testing.T) {
	items := []item[int, int]{
		{3, 3},
		{2, 2},
		{1, 1},
	}

	m := newOrderedMap(items...)
	m.Set(4, 4)

	value, ok := m.Get(4)
	assert.True(t, ok)
	assert.Equal(t, 4, value)
}

func TestMap_Set__should_update_existing_item(t *testing.T) {
	items := []item[int, int]{
		{3, 3},
		{2, 2},
		{1, 1},
	}

	m := newOrderedMap(items...)
	m.Set(1, 3)
	m.Set(2, 2)
	m.Set(3, 1)

	items1 := []item[int, int]{
		{3, 1},
		{2, 2},
		{1, 3},
	}

	assert.Equal(t, items1, m.items)
}
