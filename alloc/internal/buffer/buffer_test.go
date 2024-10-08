// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package buffer

import (
	"testing"

	"github.com/basecomplextech/baselibrary/alloc/internal/heap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testBuffer() *bufferImpl {
	h := heap.New()
	return newBuffer(h)
}

func testBufferSize(size int) *bufferImpl {
	h := heap.New()
	return newBufferSize(h, size)
}

// Acquire

func TestAcquire__should_acquire_buffer(t *testing.T) {
	b := Acquire().(*bufferImpl)
	assert.NotNil(t, b.state)

	b.Free()
}

// Free

func TestBuffer_Free__should_free_buffer(t *testing.T) {
	b := testBuffer()
	b.Free()

	assert.Nil(t, b.state)
}

// Bytes

func TestBuffer_Bytes__should_return_bytes(t *testing.T) {
	b := testBuffer()
	data := []byte("hello, world")
	b.Write(data)

	data1 := b.Bytes()
	assert.Equal(t, data, data1)
}

func TestBuffer_Bytes__should_return_empty_buffer_when_empty(t *testing.T) {
	b := testBuffer()
	data := b.Bytes()

	assert.Equal(t, []byte{}, data)
}

func TestBuffer_Bytes__should_merge_multiple_blocks(t *testing.T) {
	b := testBuffer()
	data := []byte("hello, world")

	for i, ch := range data {
		if i != 0 {
			b.allocBlock(1)
		}
		b.Write([]byte{ch})
	}
	require.Len(t, b.blocks, len(data))

	data1 := b.Bytes()
	assert.Len(t, b.blocks, 1)
	assert.Equal(t, data, data1)
}

// Grow

func TestBuffer_Grow__should_grow_buffer(t *testing.T) {
	b := testBuffer()
	data := b.Grow(10)

	assert.Len(t, data, 10)
	assert.Equal(t, 10, b.len)
}

func TestBuffer_Grow__should_alloc_next_block(t *testing.T) {
	b := testBuffer()
	size := 1024

	b.Grow(size)
	b.Grow(1)

	assert.Len(t, b.blocks, 2)
	assert.Equal(t, size+1, b.len)
}

// Write

func TestBuffer_Write__should_append_bytes(t *testing.T) {
	b := testBuffer()
	data := []byte("hello, world")

	n, err := b.Write(data)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(data), n)
}

// Reset

func TestBuffer_Reset__should_free_blocks(t *testing.T) {
	b := testBuffer()
	data := []byte("hello, world")

	b.Write(data)
	block := b.blocks[0]
	b.Reset()

	assert.Equal(t, 0, b.len)
	assert.Len(t, b.blocks, 1)
	assert.Equal(t, 0, len(block.Bytes()))
}

func TestBuffer_Reset__should_free_blocks_except_for_first_when_capacity_matches(t *testing.T) {
	b := testBufferSize(128)

	b.Grow(1024)
	b.Grow(1)
	assert.Len(t, b.blocks, 2)

	b.Reset()
	assert.Equal(t, 0, b.len)
	assert.Len(t, b.blocks, 1)
}
