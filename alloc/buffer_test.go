package alloc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testBuffer() *Buffer {
	return &Buffer{a: newAllocator()}
}

// Bytes

func TestBuffer_Bytes__should_return_bytes(t *testing.T) {
	b := testBuffer()
	data := []byte("hello, world")
	b.Write(data)

	data1 := b.Bytes()
	assert.Equal(t, data, data1)
}

func TestBuffer_Bytes__should_return_nil_when_empty(t *testing.T) {
	b := testBuffer()
	data := b.Bytes()

	assert.Nil(t, data)
}

func TestBuffer_Bytes__should_merge_multiple_blocks(t *testing.T) {
	b := testBuffer()
	data := []byte("hello, world")

	for _, ch := range data {
		b.allocBlock(1)
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
	size := blockClassSizes[0]

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

func TestBuffer_Reset__should_clear_buffer(t *testing.T) {
	b := testBuffer()
	data := []byte("hello, world")

	b.Write(data)
	b.Reset()

	assert.Equal(t, 0, b.len)
	assert.Len(t, b.blocks, 0)
}

func TestBuffer_Reset__should_clear_and_free_blocks(t *testing.T) {
	b := testBuffer()
	data := []byte("hello, world")

	b.Write(data)
	block := b.blocks[0]
	b.Reset()

	assert.Equal(t, 0, len(block.buf))
}
