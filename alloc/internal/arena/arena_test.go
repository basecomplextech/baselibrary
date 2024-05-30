package arena

import (
	"math"
	"testing"

	"github.com/basecomplextech/baselibrary/alloc/internal/heap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testArena() *arena {
	h := heap.New()
	return newArena(h)
}

// Acquire

func TestAcquireArena__should_return_pooled_arena(t *testing.T) {
	a := AcquireArena().(*arena)
	assert.NotNil(t, a.state)

	a.Free()

}

// Free

func TestArena_Free__should_free_arena(t *testing.T) {
	a := testArena()
	a.Alloc(1)

	a.Free()
	assert.Nil(t, a.state)
}

func TestArena_Free__should_reset_first_block_other_release_blocks(t *testing.T) {
	a := testArena()
	a.Alloc(1)
	a.Alloc(1024)
	require.Len(t, a.blocks, 2)

	b := a.blocks[0]
	require.Equal(t, 1, b.Len())

	s := a.state
	a.Free()
	assert.Len(t, s.blocks, 1)
	assert.Equal(t, 0, b.Len())
}

// Len

func TestArena_Len__should_return_allocated_memory(t *testing.T) {
	a := testArena()
	a.Alloc(8)
	a.Alloc(32)

	ln := a.Len()
	assert.Equal(t, int64(40), ln)
}

// Reset

func TestArena_Reset__should_reset_first_free_other_blocks(t *testing.T) {
	a := testArena()

	a.Alloc(16)
	a.Alloc(1024)
	a.Alloc(4096)
	require.Len(t, a.blocks, 3)

	b := a.blocks[0]
	a.Reset()

	assert.Equal(t, b.Cap(), int(a.cap))
	assert.Len(t, a.blocks, 1)
	assert.Equal(t, 0, len(b.Bytes()))
}

func TestArena_Reset__should_free_blocks_except_for_first_when_small(t *testing.T) {
	a := testArena()

	a.Alloc(1024)
	a.Alloc(1)
	assert.Len(t, a.blocks, 2)

	a.Reset()
	assert.Equal(t, int64(1024), a.cap)
	assert.Len(t, a.blocks, 1)
}

// alloc

func TestArena_alloc__should_allocate_data(t *testing.T) {
	a := testArena()
	b := a.Alloc(8)

	v := (*int64)(b)
	*v = math.MaxInt64

	assert.Equal(t, int64(math.MaxInt64), *v)
}

func TestArena_alloc__should_align_allocations(t *testing.T) {
	a := testArena()
	a.Alloc(3)

	b := a.blocks[0]
	assert.Equal(t, 3, b.Len())

	a.Alloc(9)
	assert.Equal(t, 8+9, b.Len())
}

func TestArena_alloc__should_not_add_padding_when_already_aligned(t *testing.T) {
	a := testArena()
	a.Alloc(8)

	b := a.blocks[0]
	assert.Equal(t, 8, b.Len())
}

func TestArena_alloc__should_allocate_next_block_when_not_enough_space(t *testing.T) {
	a := testArena()
	a.Alloc(1)

	n := a.blocks[0].Cap()
	a.Alloc(n)

	b1 := a.blocks[1]
	assert.Len(t, a.blocks, 2)
	assert.Equal(t, n, b1.Len())
}

// allocBlock

func TestArena_allocBlock__should_increment_capacity(t *testing.T) {
	a := testArena()
	a.Alloc(1)
	cp := a.blocks[0].Cap()
	assert.Equal(t, int64(cp), a.cap)

	a.allocBlock(1)
	cp += a.blocks[1].Cap()
	assert.Equal(t, int64(cp), a.cap)
}
