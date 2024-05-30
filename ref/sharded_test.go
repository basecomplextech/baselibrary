package ref

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NewSharded

func TestNewSharded__should_set_shard_refs(t *testing.T) {
	v := new(int)
	*v = 123

	r := NewNoop(v)
	s := newSharded(r)

	for i := range s.shards {
		sh := &s.shards[i]
		v1 := sh.ref.Unwrap()

		require.NotNil(t, sh.ref)
		require.Equal(t, int64(1), sh.ref.Refcount())
		require.Same(t, v, v1)
	}
}

// Get

func TestSharded_Get__should_return_shard_ref(t *testing.T) {
	v := new(int)
	*v = 123

	r := NewNoop(v)
	s := newSharded(r)

	r1, ok := s.Get()
	require.True(t, ok)

	v1 := r1.Unwrap()
	require.Same(t, v, v1)
}

func TestSharded_Get__should_retain_shard_ref(t *testing.T) {
	v := new(int)
	*v = 123

	r := NewNoop(v)
	s := newSharded(r)

	r1, ok := s.Get()
	require.True(t, ok)
	assert.Equal(t, int64(2), r1.Refcount())
}

// Set

func TestSharded_Set__should_set_shard_refs(t *testing.T) {
	v := new(int)
	*v = 123
	r := NewNoop(v)

	s := newSharded[*int](nil)
	s.Set(r)

	for i := range s.shards {
		sh := &s.shards[i]
		v1 := sh.ref.Unwrap()

		require.NotNil(t, sh.ref)
		require.Equal(t, int64(1), sh.ref.Refcount())
		require.Same(t, v, v1)
	}
}

func TestSharded_Set__should_retain_next_ref(t *testing.T) {
	v := new(int)
	*v = 123
	r := NewNoop(v)

	s := newSharded[*int](nil)
	s.Set(r)

	assert.Equal(t, int64(9), r.Refcount())
}

func TestSharded_Set__should_release_previous_ref(t *testing.T) {
	v0 := new(int)
	v1 := new(int)
	*v0 = 123
	*v1 = 456

	r0 := NewNoop(v0)
	r1 := NewNoop(v1)

	s := newSharded(r0)
	s.Set(r1)

	assert.Equal(t, int64(1), r0.Refcount())
	assert.Equal(t, int64(9), r1.Refcount())
}

// Clear

func TestSharded_Clear__should_clear_shard_refs(t *testing.T) {
	v := new(int)
	*v = 123

	r := NewNoop(v)
	s := newSharded(r)
	s.Clear()

	for i := range s.shards {
		sh := &s.shards[i]

		require.Nil(t, sh.ref)
	}
}

func TestSharded_Clear__should_release_shard_refs(t *testing.T) {
	v := new(int)
	*v = 123

	r := NewNoop(v)
	s := newSharded(r)
	s.Clear()

	assert.Equal(t, int64(1), r.Refcount())
}
