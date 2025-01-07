// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package collect

import (
	"slices"
	"testing"

	"github.com/basecomplextech/baselibrary/collect/slices2"
	"github.com/stretchr/testify/assert"
)

func TestUniquePriorityQueue_Init_Pop__should_init_queue_and_pop_items_in_order(t *testing.T) {
	items := []PriorityQueueItem[string, int]{
		{Value: "a", Priority: 1},
		{Value: "b", Priority: 2},
		{Value: "c", Priority: 3},
		{Value: "d", Priority: 4},
		{Value: "e", Priority: 5},
	}

	items1 := slices.Clone(items)
	slices2.Shuffle(items1)

	q := NewUniquePriorityQueue(items1...)
	items2 := make([]PriorityQueueItem[string, int], 0, len(items))
	for q.Len() > 0 {
		value, priority, ok := q.Poll()
		if !ok {
			t.Fatal("failed to pop")
		}

		item2 := PriorityQueueItem[string, int]{Value: value, Priority: priority}
		items2 = append(items2, item2)
	}

	assert.Equal(t, items, items2)
}

func TestUniquePriorityQueue_Push_Pop__should_push_and_pop_items_in_order(t *testing.T) {
	items := []PriorityQueueItem[string, int]{
		{Priority: 1, Value: "a"},
		{Priority: 2, Value: "b"},
		{Priority: 3, Value: "c"},
		{Priority: 4, Value: "d"},
		{Priority: 5, Value: "e"},
	}

	items1 := slices.Clone(items)
	slices2.Shuffle(items1)

	q := NewUniquePriorityQueue[string, int]()
	for _, item := range items1 {
		q.Push(item.Value, item.Priority)
	}

	items2 := make([]PriorityQueueItem[string, int], 0, len(items))
	for q.Len() > 0 {
		value, priority, ok := q.Poll()
		if !ok {
			t.Fatal("failed to pop")
		}

		item2 := PriorityQueueItem[string, int]{Value: value, Priority: priority}
		items2 = append(items2, item2)
	}

	assert.Equal(t, items, items2)
}

func TestUniquePriorityQueue_Push__should_support_duplicate_priority_items(t *testing.T) {
	q := NewUniquePriorityQueue[string, int]()
	q.Push("a", 1)
	q.Push("b", 1)
	q.Push("c", 1)

	v0, _, _ := q.Poll()
	v1, _, _ := q.Poll()
	v2, _, _ := q.Poll()

	values := []string{"a", "c", "b"}
	values1 := []string{v0, v1, v2}
	assert.Equal(t, values, values1)
}

func TestUniquePriorityQueue_Push__should_update_existing_priority(t *testing.T) {
	q := NewUniquePriorityQueue[string, int]()
	q.Push("a", 1)
	q.Push("a", 2)
	q.Push("a", 3)

	v, p, ok := q.Poll()
	assert.True(t, ok)
	assert.Equal(t, "a", v)
	assert.Equal(t, 3, p)
}
