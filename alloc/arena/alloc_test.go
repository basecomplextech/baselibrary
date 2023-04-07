package arena

import (
	"math"
	"testing"

	"github.com/complex1tech/baselibrary/alloc/internal/heap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testArena returns a new arena with a new heap, for tests only.
func testArena() *arena {
	h := heap.New()
	return newArena(h)
}

// Free

func TestArena_Free__should_release_blocks(t *testing.T) {
	a := testArena()
	a.alloc(1)

	last := a.lastBlock()
	require.Equal(t, 1, last.Len())

	a.Free()
	assert.Len(t, a.blocks, 0)
	assert.Equal(t, 0, last.Len())
}

// Len

func TestArena_Len__should_return_allocated_memory(t *testing.T) {
	a := testArena()
	a.alloc(8)
	a.alloc(32)

	ln := a.Len()
	assert.Equal(t, int64(40), ln)
}

// alloc

func TestArena_alloc__should_allocate_data(t *testing.T) {
	a := testArena()
	b := a.alloc(8)

	v := (*int64)(b)
	*v = math.MaxInt64

	assert.Equal(t, int64(math.MaxInt64), *v)
}

func TestArena_alloc__should_align_allocations(t *testing.T) {
	a := testArena()
	a.alloc(3)

	block := a.lastBlock()
	assert.Equal(t, 3, block.Len())

	a.alloc(9)
	assert.Equal(t, 8+9, block.Len())
}

func TestArena_alloc__should_not_add_padding_when_already_aligned(t *testing.T) {
	a := testArena()
	a.alloc(8)

	block := a.lastBlock()
	assert.Equal(t, 8, block.Len())
}

func TestArena_alloc__should_allocate_next_block_when_not_enough_space(t *testing.T) {
	a := testArena()
	a.alloc(1)

	n := a.lastBlock().Cap()
	a.alloc(n)

	last := a.lastBlock()
	assert.Len(t, a.blocks, 2)
	assert.Equal(t, last.Len(), n)
}

// allocBlock

func TestArena_allocBlock__should_increment_size(t *testing.T) {
	a := testArena()
	a.alloc(1)
	size := a.lastBlock().Cap()
	assert.Equal(t, int64(size), a.cap)

	a.allocBlock(1)
	size += a.lastBlock().Cap()
	assert.Equal(t, int64(size), a.cap)
}
