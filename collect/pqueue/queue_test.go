package pqueue

import (
	"testing"

	"github.com/complex1tech/baselibrary/slices"
	"github.com/stretchr/testify/assert"
)

func TestQueue_Init_Pop__should_init_queue_and_pop_items_in_order(t *testing.T) {
	items := []Item[string, int]{
		{Priority: 1, Value: "a"},
		{Priority: 2, Value: "b"},
		{Priority: 3, Value: "c"},
		{Priority: 4, Value: "d"},
		{Priority: 5, Value: "e"},
	}
	values := make([]string, 0, len(items))
	for _, item := range items {
		values = append(values, item.Value)
	}

	items1 := slices.Clone(items)
	slices.Shuffle(items1)

	q := Ordered(items1...)
	values1 := make([]string, 0, len(items))
	for q.Len() > 0 {
		value, ok := q.Pop()
		if !ok {
			t.Fatal("failed to pop")
		}
		values1 = append(values1, value)
	}

	assert.Equal(t, values, values1)
}

func TestQueue_Push_Pop__should_push_and_pop_items_in_order(t *testing.T) {
	items := []Item[string, int]{
		{Priority: 1, Value: "a"},
		{Priority: 2, Value: "b"},
		{Priority: 3, Value: "c"},
		{Priority: 4, Value: "d"},
		{Priority: 5, Value: "e"},
	}
	values := make([]string, 0, len(items))
	for _, item := range items {
		values = append(values, item.Value)
	}

	items1 := slices.Clone(items)
	slices.Shuffle(items1)

	q := Ordered[string, int]()
	for _, item := range items1 {
		q.Push(item.Value, item.Priority)
	}

	values1 := make([]string, 0, len(items))
	for q.Len() > 0 {
		value, ok := q.Pop()
		if !ok {
			t.Fatal("failed to pop")
		}
		values1 = append(values1, value)
	}

	assert.Equal(t, values, values1)
}

func TestQueue_Push__should_support_duplicate_priority_items(t *testing.T) {
	q := Ordered[string, int]()
	q.Push("a", 1)
	q.Push("b", 1)
	q.Push("c", 1)

	v0, _ := q.Pop()
	v1, _ := q.Pop()
	v2, _ := q.Pop()

	values := []string{"a", "c", "b"}
	values1 := []string{v0, v1, v2}
	assert.Equal(t, values, values1)
}
