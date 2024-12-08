// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Acquire

func TestShardedVar_Acquire__should_acquire_current_reference(t *testing.T) {
	r := NewNoop(1)

	v := newShardedVar[int]()
	v.SetRetain(r)
	r.Release()

	r1, ok := v.Acquire()
	if !ok {
		t.Fatal(ok)
	}
	assert.Equal(t, int64(2), r1.Refcount())
	assert.Equal(t, int64(1), r.Refcount())

	v.Unset()
	assert.Equal(t, int64(1), r1.Refcount())
	assert.Equal(t, int64(1), r.Refcount())

	r1.Release()
	assert.Equal(t, int64(0), r1.Refcount())
	assert.Equal(t, int64(0), r.Refcount())
}

func TestShardedVar_Acquire__should_return_false_when_unset(t *testing.T) {
	v := newShardedVar[int]()

	r, ok := v.Acquire()
	assert.False(t, ok)
	assert.Nil(t, r)
}

// SetRetain

func TestShardedVar_SetRetain__should_retain_new_reference(t *testing.T) {
	r := NewNoop(1)

	v := newShardedVar[int]()
	v.SetRetain(r)
	assert.Equal(t, int64(2), r.Refcount())

	r.Release()
	assert.Equal(t, int64(1), r.Refcount())
}

func TestShardedVar_SetRetain__should_release_previous_reference(t *testing.T) {
	r0 := NewNoop(1)
	r1 := NewNoop(2)

	v := newShardedVar[int]()
	v.SetRetain(r0)
	v.SetRetain(r1)

	r0.Release()
	r1.Release()

	assert.Equal(t, int64(0), r0.Refcount())
	assert.Equal(t, int64(1), r1.Refcount())

	r10, ok := v.Acquire()
	require.True(t, ok)

	r11, ok := v.Acquire()
	require.True(t, ok)

	v.Unset()
	assert.Equal(t, int64(1), r1.Refcount())

	r10.Release()
	assert.Equal(t, int64(1), r1.Refcount())

	r11.Release()
	assert.Equal(t, int64(0), r1.Refcount())
}

// Shard

func TestShardedVarShard__should_have_cache_line_size(t *testing.T) {
	s := unsafe.Sizeof(shardedVarShard[any]{})

	assert.Equal(t, 256, int(s))
}

func TestShardedVarRefCount__should_have_cache_line_size(t *testing.T) {
	s := unsafe.Sizeof(varRefCount{})

	assert.Equal(t, 256, int(s))
}
