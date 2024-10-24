// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"testing"

	"github.com/basecomplextech/baselibrary/bin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRef(t *testing.T) {
	freed := false

	obj := FreeFunc(func() {
		freed = true
	})
	ref := New(obj)
	ref.Retain()
	require.Equal(t, int64(2), ref.Refcount())

	ref.Release()
	require.False(t, freed)

	ref.Release()
	require.True(t, freed)
}

// Freer

func TestRefFreer(t *testing.T) {
	freed := false
	freer := FreeFunc(func() {
		freed = true
	})

	ref := NewFreer(10, freer)
	ref.Retain()

	v := ref.Unwrap()
	require.Equal(t, 10, v)

	ref.Release()
	require.False(t, freed)

	ref.Release()
	require.True(t, freed)
}

func TestRefFreerPooled(t *testing.T) {
	freed := false
	freer := FreeFunc(func() {
		freed = true
	})

	v := bin.Random256()
	ref := NewFreer(v, freer)
	ref.Retain()

	v1 := ref.Unwrap()
	require.Equal(t, v, v1)

	ref.Release()
	require.False(t, freed)

	ref.Release()
	require.True(t, freed)
}

// Next

func TestRefNext(t *testing.T) {
	ref := NewNoop(10)
	ref1 := NextRetain(20, ref)
	require.Equal(t, int64(2), ref.Refcount())
	require.Equal(t, int64(1), ref1.Refcount())

	v := ref1.Unwrap()
	require.Equal(t, 20, v)

	ref1.Retain()
	require.Equal(t, int64(2), ref.Refcount())
	require.Equal(t, int64(2), ref1.Refcount())

	ref1.Release()
	ref1.Release()
	assert.Equal(t, int64(1), ref.Refcount())
}

func TestRefNextPooled(t *testing.T) {
	ref := NewNoop(10)

	v := bin.Random256()
	ref1 := NextRetain(v, ref)
	require.Equal(t, int64(2), ref.Refcount())
	require.Equal(t, int64(1), ref1.Refcount())

	v1 := ref1.Unwrap()
	require.Equal(t, v1, v)

	ref1.Retain()
	require.Equal(t, int64(2), ref.Refcount())
	require.Equal(t, int64(2), ref1.Refcount())

	ref1.Release()
	ref1.Release()
	assert.Equal(t, int64(1), ref.Refcount())
}

// Noop

func TestRefNoop(t *testing.T) {
	ref := NewNoop(10)
	ref.Retain()

	v := ref.Unwrap()
	require.Equal(t, 10, v)

	ref.Release()
}
