package alloc

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Used

func TestArena_Used__should_return_allocated_memory(t *testing.T) {
	a := newArena()
	a.alloc(8)
	a.alloc(32)

	used := a.Used()
	assert.Equal(t, int64(40), used)
}

// Free

func TestArena_Free__should_release_blocks(t *testing.T) {
	a := newArena()
	a.alloc(1)

	last := a.last()
	require.Equal(t, 8, len(last.buf))

	a.Free()
	assert.Len(t, a.blocks, 0)
	assert.Equal(t, 0, len(last.buf))
}

// alloc

func TestArena_alloc__should_allocate_data(t *testing.T) {
	a := newArena()
	b := a.alloc(8)

	v := (*int64)(b)
	*v = math.MaxInt64

	assert.Equal(t, int64(math.MaxInt64), *v)
}

func TestArena_alloc__should_add_padding_for_alignment(t *testing.T) {
	a := newArena()
	a.alloc(3)

	block := a.last()
	assert.Equal(t, 8, len(block.buf))

	a.alloc(9)
	assert.Equal(t, 24, len(block.buf))
}

func TestArena_alloc__should_not_add_padding_when_already_aligned(t *testing.T) {
	a := newArena()
	a.alloc(8)

	block := a.last()
	assert.Equal(t, 8, len(block.buf))
}

func TestArena_alloc__should_allocate_next_block_when_not_enough_space(t *testing.T) {
	a := newArena()
	a.alloc(1)

	n := a.last().cap()
	a.alloc(n)

	last := a.last()
	assert.Len(t, a.blocks, 2)
	assert.Equal(t, len(last.buf), n)
}

// allocBlock

func TestArena_allocBlock__should_increment_size(t *testing.T) {
	a := newArena()
	a.alloc(1)
	size := a.last().cap()
	assert.Equal(t, int64(size), a.size)

	a.allocBlock(1)
	size += a.last().cap()
	assert.Equal(t, int64(size), a.size)
}
