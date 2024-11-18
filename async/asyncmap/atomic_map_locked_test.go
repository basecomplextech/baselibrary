// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtomicLockedMap_Range__should_iterate_over_items(t *testing.T) {
	m := newAtomicMap[int, int](0)
	n := 128
	items := make([]int, 0, n)

	for i := 0; i < n; i++ {
		m.Set(i, i)
		items = append(items, i)
	}

	locked := m.LockMap()
	defer locked.Free()

	items1 := make([]int, 0, n)
	locked.Range(func(k int, v int) bool {
		items1 = append(items1, v)
		return true
	})

	assert.Equal(t, items, items1)
}

func TestAtomicLockedMap_Clear__should_clear_items(t *testing.T) {
	m := newAtomicMap[int, int](0)
	n := 128
	items := make([]int, 0, n)

	for i := 0; i < n; i++ {
		m.Set(i, i)
		items = append(items, i)
	}

	locked := m.LockMap()
	defer locked.Free()

	locked.Clear()

	items1 := make([]int, 0, n)
	locked.Range(func(k int, v int) bool {
		items1 = append(items1, v)
		return true
	})

	assert.Empty(t, items1)
}
