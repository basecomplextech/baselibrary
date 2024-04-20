package refqueue

import (
	"testing"

	"github.com/basecomplextech/baselibrary/ref"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Clear

func TestQueue_Clear__should_release_items(t *testing.T) {
	r := ref.NewNoop(123)

	q := newQueue[ref.R[int]]()
	q.Push(r)
	q.Clear()

	assert.Equal(t, int64(1), r.Refcount())
}

// Push

func TestQueue_Push__should_retain_item(t *testing.T) {
	r := ref.NewNoop(123)

	q := newQueue[ref.R[int]]()
	q.Push(r)

	assert.Equal(t, int64(2), r.Refcount())
}

// Pop

func TestQueue_Pop__should_pop_item_noop(t *testing.T) {
	r := ref.NewNoop(123)

	q := newQueue[ref.R[int]]()
	q.Push(r)

	r1, ok := q.Pop()
	require.True(t, ok)

	assert.Equal(t, int64(2), r.Refcount())
	assert.Same(t, r, r1)
}

// Free

func TestQueue_Free__should_release_items_close_queue(t *testing.T) {
	r := ref.NewNoop(123)

	q := newQueue[ref.R[int]]()
	q.Push(r)
	q.Free()

	assert.Equal(t, int64(1), r.Refcount())
}
