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

func TestPriorityQueue_Init_Pop__should_init_queue_and_pop_items_in_order(t *testing.T) {
	items := []PriorityQueueItem[string, int]{
		{Value: "a", Priority: 1},
		{Value: "b", Priority: 2},
		{Value: "c", Priority: 3},
		{Value: "d", Priority: 4},
		{Value: "e", Priority: 5},
	}
	values := make([]string, 0, len(items))
	for _, item := range items {
		values = append(values, item.Value)
	}

	items1 := slices.Clone(items)
	slices2.Shuffle(items1)
	compare := func(a, b int) int { return a - b }

	q := NewPriorityQueueCompare(compare, items1...)
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

func TestPriorityQueue_Push_Pop__should_push_and_pop_items_in_order(t *testing.T) {
	items := []PriorityQueueItem[string, int]{
		{Value: "a", Priority: 1},
		{Value: "b", Priority: 2},
		{Value: "c", Priority: 3},
		{Value: "d", Priority: 4},
		{Value: "e", Priority: 5},
	}

	items1 := slices.Clone(items)
	slices2.Shuffle(items1)
	compare := func(a, b int) int { return a - b }

	q := NewPriorityQueueCompare[string, int](compare)
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

func TestPriorityQueue_Push__should_support_duplicate_priority_items(t *testing.T) {
	compare := func(a, b int) int { return a - b }

	q := NewPriorityQueueCompare[string, int](compare)
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
