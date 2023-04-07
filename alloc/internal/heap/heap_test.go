package heap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHeap_Alloc_Free__should_allocate_and_free_block(t *testing.T) {
	h := New()

	for size := 1; size <= (maxSize << 2); size *= 2 {
		b := h.Alloc(size)
		require.True(t, cap(b.buf) >= size, "size=%d, cap=%d", size, cap(b.buf))
		h.Free(b)

		b = h.Alloc(size + 1)
		require.True(t, cap(b.buf) >= size, "size=%d, cap=%d", size, cap(b.buf))
		h.Free(b)
	}
}

func TestAllocator_Alloc__should_allocate_block(t *testing.T) {
	h := New()

	cases := []struct {
		size      int
		poolSize  int
		poolIndex int
	}{
		{0, 1024, 10},
		{1, 1024, 10},
		{1023, 1024, 10},
		{1024, 1024, 10},
		{2047, 2048, 11},
		{2048, 2048, 11},
		{1<<24 + 1, 1 << 25, 25},
	}

	for _, c := range cases {
		b := h.Alloc(c.size)
		i := blockPool(cap(b.buf))
		require.Equal(t, c.poolSize, cap(b.buf))
		require.Equal(t, c.poolIndex, i, "size=%d", c.size)
	}
}

func TestHeap_FreeMany__should_free_blocks(t *testing.T) {
	h := New()
	sizes := []int{1, 1023, 1024, 2048}

	blocks := []*Block{}
	for _, size := range sizes {
		b := h.Alloc(size)
		blocks = append(blocks, b)
	}

	h.FreeMany(blocks...)
}
