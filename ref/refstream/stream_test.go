// Copyright 2024 Ivan Korobkov. All rights reserved.

package refstream

import (
	"testing"

	"github.com/basecomplextech/baselibrary/ref"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStream__should_increment_refs(t *testing.T) {
	src := NewSource[int]()
	src.Filter(func(r ref.R[int]) bool {
		v := r.Unwrap()
		return v > 10
	}).Subscribe()

	r0 := ref.NewNoop(1)
	r1 := ref.NewNoop(11)
	r2 := ref.NewNoop(100)

	src.Send(r0)
	src.Send(r1)
	src.Send(r2)

	assert.Equal(t, 1, r0.Refcount())
	assert.Equal(t, 2, r1.Refcount())
	assert.Equal(t, 2, r2.Refcount())
}

func TestStreamQueue__should_decrement_refs_on_clear(t *testing.T) {
	src := NewSource[int]()
	q := src.Subscribe()

	r0 := ref.NewNoop(1)
	r1 := ref.NewNoop(2)

	src.Send(r0)
	src.Send(r1)

	assert.Equal(t, 2, r0.Refcount())
	assert.Equal(t, 2, r1.Refcount())

	q.Clear()

	assert.Equal(t, 1, r0.Refcount())
	assert.Equal(t, 1, r1.Refcount())
}

func TestStreamMap__should_map_refs(t *testing.T) {
	src := NewSource[int]()
	q := Map(src, func(r ref.R[int]) ref.R[int32] {
		v := r.Unwrap()
		v1 := (int32)(v)
		return ref.Map(v1, r)
	}).Subscribe()

	r0 := ref.NewNoop(1)
	r1 := ref.NewNoop(2)

	src.Send(r0)
	src.Send(r1)

	assert.Equal(t, 2, r0.Refcount())
	assert.Equal(t, 2, r1.Refcount())

	m0, ok := q.Pop()
	require.True(t, ok)
	assert.Equal(t, int32(1), m0.Unwrap())

	m1, ok := q.Pop()
	require.True(t, ok)
	assert.Equal(t, int32(2), m1.Unwrap())

	m0.Release()
	m1.Release()

	assert.Equal(t, 1, r0.Refcount())
	assert.Equal(t, 1, r1.Refcount())
}
