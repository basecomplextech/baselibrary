package refstream

import (
	"testing"

	"github.com/basecomplextech/baselibrary/ref"
	"github.com/zeebo/assert"
)

func TestStream__should_increment_refs(t *testing.T) {
	src := NewSource[*ref.R[int]]()
	src.Filter(func(r *ref.R[int]) bool {
		v := r.Unwrap()
		return v > 10
	}).Subscribe()

	r0 := ref.NewNoFreer(1)
	r1 := ref.NewNoFreer(11)
	r2 := ref.NewNoFreer(100)

	src.Send(r0)
	src.Send(r1)
	src.Send(r2)

	assert.Equal(t, 1, r0.Refcount())
	assert.Equal(t, 2, r1.Refcount())
	assert.Equal(t, 2, r2.Refcount())
}

func TestStreamQueue__should_decrement_refs_on_clear(t *testing.T) {
	src := NewSource[*ref.R[int]]()
	q := src.Subscribe()

	r0 := ref.NewNoFreer(1)
	r1 := ref.NewNoFreer(2)

	src.Send(r0)
	src.Send(r1)

	assert.Equal(t, 2, r0.Refcount())
	assert.Equal(t, 2, r1.Refcount())

	q.Clear()

	assert.Equal(t, 1, r0.Refcount())
	assert.Equal(t, 1, r1.Refcount())
}
