// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package orderedmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew__should_create_new_map(t *testing.T) {
	items := []Item[int, int]{
		{3, 3},
		{2, 2},
		{1, 1},
	}

	m := New[int, int](items...)

	items1 := []Item[int, int]{}
	m.Iterate(func(key int, value int) bool {
		item := Item[int, int]{key, value}
		items1 = append(items1, item)
		return true
	})

	assert.Equal(t, items, items1)
}

// Index

func TestMap_Index__should_return_index_of_key(t *testing.T) {
	items := []Item[int, int]{
		{3, 3},
		{2, 2},
		{1, 1},
	}

	m := New[int, int](items...)
	assert.Equal(t, 0, m.Index(3))
	assert.Equal(t, 1, m.Index(2))
	assert.Equal(t, 2, m.Index(1))
}

// Get

func TestMap_Get__should_return_value_for_key(t *testing.T) {
	items := []Item[int, int]{
		{3, 3},
		{2, 2},
		{1, 1},
	}

	m := New[int, int](items...)
	value, ok := m.Get(2)
	assert.True(t, ok)
	assert.Equal(t, 2, value)

	value, ok = m.Get(4)
	assert.False(t, ok)
	assert.Equal(t, 0, value)
}

// Put

func TestMap_Put__should_add_new_key_value_pair(t *testing.T) {
	items := []Item[int, int]{
		{3, 3},
		{2, 2},
		{1, 1},
	}

	m := New[int, int](items...)
	m.Put(4, 4)

	value, ok := m.Get(4)
	assert.True(t, ok)
	assert.Equal(t, 4, value)
}

func TestMap_Put__should_update_existing_key_value_pair(t *testing.T) {
	items := []Item[int, int]{
		{3, 3},
		{2, 2},
		{1, 1},
	}

	m := New[int, int](items...)
	m.Put(1, 3)
	m.Put(2, 2)
	m.Put(3, 1)

	items1 := []Item[int, int]{
		{3, 1},
		{2, 2},
		{1, 3},
	}

	items2 := []Item[int, int]{}
	m.Iterate(func(key int, value int) bool {
		item := Item[int, int]{key, value}
		items2 = append(items2, item)
		return true
	})
	assert.Equal(t, items1, items2)
}

// Delete

func TestMap_Delete__should_delete_key_value_pair(t *testing.T) {
	items := []Item[int, int]{
		{3, 3},
		{2, 2},
		{1, 1},
	}

	m := New[int, int](items...)
	m.Delete(2)

	items1 := []Item[int, int]{
		{3, 3},
		{1, 1},
	}

	items2 := []Item[int, int]{}
	m.Iterate(func(key int, value int) bool {
		item := Item[int, int]{key, value}
		items2 = append(items2, item)
		return true
	})
	assert.Equal(t, items1, items2)
}

// Keys

func TestMap_Keys__should_return_keys(t *testing.T) {
	items := []Item[int, int]{
		{3, 3},
		{2, 2},
		{1, 1},
	}

	m := New[int, int](items...)
	keys := m.Keys()

	assert.Equal(t, []int{3, 2, 1}, keys)
}

// Values

func TestMap_Values__should_return_values(t *testing.T) {
	items := []Item[int, int]{
		{3, 3},
		{2, 2},
		{1, 1},
	}

	m := New[int, int](items...)
	values := m.Values()

	assert.Equal(t, []int{3, 2, 1}, values)
}
