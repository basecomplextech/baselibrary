// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAtomicMap_Len__should_return_number_of_items(t *testing.T) {
	m := newAtomicMap[int, int](0)
	n := 128

	for i := 0; i < n; i++ {
		m.Set(i, i)
	}

	n1 := m.Len()
	assert.Equal(t, n, n1)
}

// Clear

func TestAtomicMap_Clear__should_delete_all_items(t *testing.T) {
	m := newAtomicMap[int, int](0)
	n := 128

	for i := 0; i < n; i++ {
		m.Set(i, i)
	}

	m.Clear()

	n1 := m.Len()
	assert.Equal(t, 0, n1)
}

// Contains

func TestAtomicMap_Contains__should_return_true_if_key_exists(t *testing.T) {
	m := newAtomicMap[int, int](0)
	n := 128

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

func TestAtomicMap_Get__should_return_value_by_key(t *testing.T) {
	m := newAtomicMap[int, int](0)
	n := 128

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

func TestAtomicMap_Get__should_reduce_entry_refs_when_they_reach_1000_000_000(t *testing.T) {
	m := newAtomicMap[int, int](0)
	m.Set(1, 1)

	state := m.state.Load()
	b := &state.buckets[1]
	b.ref.Add(1000_000_000 - 2)

	_, ext := unpackAtomicMapEntryRef(b.ref.Load())
	require.Equal(t, 999_999_999, int(ext))

	e := b.entry.Load()
	int0 := e.refs.Add(-1000_000_000 + 2)
	require.Equal(t, -999_999_998, int(int0))

	m.Get(1)

	_, ext1 := unpackAtomicMapEntryRef(b.ref.Load())
	require.Equal(t, 1000_000, int(ext1))

	int1 := e.refs.Load()
	require.Equal(t, -999_999, int(int1))
}

// GetOrSet

func TestAtomicMap_GetOrSet__should_return_value(t *testing.T) {
	m := newAtomicMap[int, int](0)
	n := 128

	for i := 0; i < n; i++ {
		m.Set(i, i)
	}

	for i := 0; i < n; i++ {
		v, ok := m.GetOrSet(i, i)
		require.True(t, ok)
		require.Equal(t, i, v)
	}
}

func TestAtomicMap_GetOrSet__should_set_value_if_not_exists(t *testing.T) {
	m := newAtomicMap[int, int](0)
	n := 128

	for i := 0; i < n; i++ {
		v, ok := m.GetOrSet(i, i)
		require.False(t, ok)
		require.Equal(t, i, v)
	}
}

// Delete

func TestAtomicMap_Delete__should_delete_value(t *testing.T) {
	m := newAtomicMap[int, int](0)
	n := 128

	for i := 0; i < n; i++ {
		m.Set(i, i)
	}

	for i := 0; i < n; i++ {
		m.Delete(i)
	}

	n1 := m.Len()
	assert.Equal(t, 0, n1)
}

func TestAtomicMap_Delete__should_skip_absent_key(t *testing.T) {
	m := newAtomicMap[int, int](0)
	n := 128

	for i := 0; i < n; i++ {
		m.Delete(i)
	}

	n1 := m.Len()
	assert.Equal(t, 0, n1)
}

func TestAtomicMap_Delete__should_delete_and_return_value(t *testing.T) {
	m := newAtomicMap[int, int](0)
	n := 128

	for i := 0; i < n; i++ {
		m.Set(i, i)
	}

	for i := 0; i < n; i++ {
		v, ok := m.Delete(i)
		require.True(t, ok)
		require.Equal(t, i, v)
	}

	n1 := m.Len()
	assert.Equal(t, 0, n1)
}

func TestAtomicMap_Delete__should_return_false_if_key_not_exists(t *testing.T) {
	m := newAtomicMap[int, int](0)
	n := 128

	for i := 0; i < n; i++ {
		_, ok := m.Delete(i)
		require.False(t, ok)
	}
}

// Set

func TestAtomicMap_Set__should_set_value(t *testing.T) {
	m := newAtomicMap[int, int](0)
	n := 128

	for i := 0; i < n; i++ {
		m.Set(i, i)
	}

	for i := 0; i < n; i++ {
		v, ok := m.Get(i)
		require.True(t, ok)
		require.Equal(t, i, v)
	}
}

func TestAtomicMap_Set__should_resize_map_on_threshold(t *testing.T) {
	m := newAtomicMap[int, int](0)

	n := 16
	threshold := int(float64(n) * atomicMapThreshold)

	for i := 0; i < threshold-1; i++ {
		m.Set(i, i)
	}

	size := len(m.state.Load().buckets)
	assert.Equal(t, n, size)

	m.Set(threshold, threshold)

	size1 := len(m.state.Load().buckets)
	assert.Equal(t, n*2, size1)
}

// SetAbsent

func TestAtomicMap_SetAbsent__should_set_value_if_absent(t *testing.T) {
	m := newAtomicMap[int, int](0)
	n := 128

	for i := 0; i < n; i++ {
		m.SetAbsent(i, i)
	}

	for i := 0; i < n; i++ {
		v, ok := m.Get(i)
		require.True(t, ok)
		require.Equal(t, i, v)
	}
}

func TestAtomicMap_SetAbsent__should_not_set_value_if_exists(t *testing.T) {
	m := newAtomicMap[int, int](0)
	n := 128

	for i := 0; i < n; i++ {
		m.Set(i, i)
	}

	for i := 0; i < n; i++ {
		ok := m.SetAbsent(i, i*2)
		assert.False(t, ok)
	}

	for i := 0; i < n; i++ {
		v, ok := m.Get(i)
		require.True(t, ok)
		require.Equal(t, i, v)
	}
}

// Swap

func TestAtomicMap_Swap__should_swap_value(t *testing.T) {
	m := newAtomicMap[int, int](0)
	n := 128

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

func TestAtomicMap_Range__should_iterate_over_all_items(t *testing.T) {
	m := newAtomicMap[int, int](0)
	n := 128
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

	assert.Equal(t, items, items1)
}

// LockMap

func TestAtomicMap_LockMap__should_lock_map(t *testing.T) {
	m := newAtomicMap[int, int](0)
	n := 128
	items := make([]int, 0, n)

	for i := 0; i < n; i++ {
		m.Set(i, i)
		items = append(items, i)
	}

	locked := m.LockMap()
	locked.Free()
}
