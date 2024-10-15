// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestConcurrentVar__should_have_cache_line_size(t *testing.T) {
	s := concurrentSlot[int]{}
	size := unsafe.Sizeof(s)

	assert.Equal(t, 256, int(size))
}

// Acquire

func TestConcurrentVar_Acquire__should_acquire_chained_reference(t *testing.T) {
	r := NewNoop(1)

	v := NewConcurrentVar[int]()
	v.SwapRetain(r)
	r.Release()

	r1, ok := v.Acquire()
	if !ok {
		t.Fatal(ok)
	}
	assert.Equal(t, int64(2), r1.Refcount())
	assert.Equal(t, int64(concurrentNum), r.Refcount())

	v.Clear()
	assert.Equal(t, int64(1), r1.Refcount())
	assert.Equal(t, int64(1), r.Refcount())

	r1.Release()
	assert.Equal(t, int64(0), r1.Refcount())
	assert.Equal(t, int64(0), r.Refcount())
}

// SwapRetain

func TestConcurrentVar_SwapRetain__should_retain_new_reference(t *testing.T) {
	r := NewNoop(1)

	v := NewConcurrentVar[int]()
	v.SwapRetain(r)

	r.Release()
	assert.Equal(t, int64(concurrentNum), r.Refcount())
}

func TestConcurrentVar_SwapRetain__should_release_previous_reference(t *testing.T) {
	r0 := NewNoop(1)
	r1 := NewNoop(2)

	v := NewConcurrentVar[int]()
	v.SwapRetain(r0)
	v.SwapRetain(r1)

	r0.Release()
	r1.Release()

	assert.Equal(t, int64(0), r0.Refcount())
	assert.Equal(t, int64(concurrentNum), r1.Refcount())
}
