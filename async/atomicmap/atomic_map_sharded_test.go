// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the Business Source License (BSL 1.1)
// that can be found in the LICENSE file.

package atomics

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAtomicShardedMap_Len__should_return_number_of_items(t *testing.T) {
	m := newAtomicShardedMap[int, int](0)
	n := 1024

	for i := 0; i < n; i++ {
		m.Set(i, i)

		n1 := m.Len()
		require.Equal(t, i+1, n1, i)
	}

	n1 := m.Len()
	assert.Equal(t, n, n1)
}

// Clear

func TestAtomicShardedMap_Clear__should_delete_all_items(t *testing.T) {
	m := newAtomicShardedMap[int, int](0)
	n := 1024

	for i := 0; i < n; i++ {
		m.Set(i, i)
	}

	m.Clear()

	n1 := m.Len()
	assert.Equal(t, 0, n1)
}

// Contains

func TestAtomicShardedMap_Contains__should_return_true_if_key_exists(t *testing.T) {
	m := newAtomicShardedMap[int, int](0)
	n := 1024

	for i := 0; i < n; i++ {
		m.Set(i, i)
	}

	for i := 0; i < n; i++ {
		ok := m.Contains(i)
		require.True(t, ok)
	}

	ok := m.Contains(n)
	assert.False(t, ok)
}

// Get

func TestAtomicShardedMap_Get__should_return_value_by_key(t *testing.T) {
	m := newAtomicShardedMap[int, int](0)
	n := 1024

	for i := 0; i < n; i++ {
		m.Set(i, i)
	}

	for i := 0; i < n; i++ {
		v, ok := m.Get(i)
		require.True(t, ok)
		require.Equal(t, i, v)
	}

	_, ok := m.Get(n)
	assert.False(t, ok)
}

// GetOrSet

func TestAtomicShardedMap_GetOrSet__should_return_value(t *testing.T) {
	m := newAtomicShardedMap[int, int](0)
	n := 1024

	for i := 0; i < n; i++ {
		m.Set(i, i)
	}

	for i := 0; i < n; i++ {
		v, ok := m.GetOrSet(i, i)
		require.True(t, ok)
		require.Equal(t, i, v)
	}
}

func TestAtomicShardedMap_GetOrSet__should_set_value_if_not_exists(t *testing.T) {
	m := newAtomicShardedMap[int, int](0)
	n := 1024

	for i := 0; i < n; i++ {
		v, ok := m.GetOrSet(i, i)
		require.False(t, ok)
		require.Equal(t, i, v)
	}
}

// Delete

func TestAtomicShardedMap_Delete__should_delete_value(t *testing.T) {
	m := newAtomicShardedMap[int, int](0)
	n := 1024

	for i := 0; i < n; i++ {
		m.Set(i, i)
	}

	for i := 0; i < n; i++ {
		m.Delete(i)
	}

	n1 := m.Len()
	assert.Equal(t, 0, n1)
}

func TestAtomicShardedMap_Delete__should_skip_absent_key(t *testing.T) {
	m := newAtomicShardedMap[int, int](0)
	n := 1024

	for i := 0; i < n; i++ {
		m.Delete(i)
	}

	n1 := m.Len()
	assert.Equal(t, 0, n1)
}

// Pop

func TestAtomicShardedMap_Pop__should_delete_and_return_value(t *testing.T) {
	m := newAtomicShardedMap[int, int](0)
	n := 1024

	for i := 0; i < n; i++ {
		m.Set(i, i)
	}

	for i := 0; i < n; i++ {
		v, ok := m.Pop(i)
		require.True(t, ok)
		require.Equal(t, i, v)
	}

	n1 := m.Len()
	assert.Equal(t, 0, n1)
}

func TestAtomicShardedMap_Pop__should_return_false_if_key_not_exists(t *testing.T) {
	m := newAtomicShardedMap[int, int](0)
	n := 1024

	for i := 0; i < n; i++ {
		_, ok := m.Pop(i)
		require.False(t, ok)
	}
}

// Set

func TestAtomicShardedMap_Set__should_set_value(t *testing.T) {
	m := newAtomicShardedMap[int, int](0)
	n := 1024

	for i := 0; i < n; i++ {
		m.Set(i, i)
	}

	for i := 0; i < n; i++ {
		v, ok := m.Get(i)
		require.True(t, ok)
		require.Equal(t, i, v)
	}
}

// Swap

func TestAtomicShardedMap_Swap__should_swap_value(t *testing.T) {
	m := newAtomicShardedMap[int, int](0)
	n := 1024

	for i := 0; i < n; i++ {
		m.Set(i, i)
	}

	for i := 0; i < n; i++ {
		v, ok := m.Swap(i, i*2)
		require.True(t, ok)
		require.Equal(t, i, v)
	}

	for i := 0; i < n; i++ {
		v, ok := m.Get(i)
		require.True(t, ok)
		require.Equal(t, i*2, v)
	}
}

// Range

func TestAtomicShardedMap_Range__should_iterate_over_all_items(t *testing.T) {
	m := newAtomicShardedMap[int, int](0)
	n := 1024
	items := make([]int, 0, n)

	for i := 0; i < n; i++ {
		m.Set(i, i)
		items = append(items, i)
	}

	items1 := make([]int, 0, n)
	m.Range(func(k int, v int) bool {
		items1 = append(items1, v)
		return true
	})

	slices.Sort(items1)
	assert.Equal(t, items, items1)
}
