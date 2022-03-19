package alloc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHeap_allocBlock__should_allocate_block(t *testing.T) {
	h := newHeap()

	cases := []struct {
		size      int
		blockSize int
		cls       int
	}{
		{0, 1024, 0},
		{1, 1024, 0},
		{1023, 1024, 0},
		{1024, 1024, 0},
		{2047, 2048, 1},
		{2048, 2048, 1},
		{1<<24 + 1, 1<<24 + 1, -1},
	}

	for _, c := range cases {
		block, cls := h.allocBlock(c.size)
		require.Equal(t, c.blockSize, block.size())
		require.Equal(t, c.cls, cls)
	}
}

func TestHeap_freeBlocks__should_free_blocks(t *testing.T) {
	h := newHeap()
	sizes := []int{1, 1023, 1024, 2048}

	blocks := []*block{}
	for _, size := range sizes {
		block, _ := h.allocBlock(size)
		blocks = append(blocks, block)
	}

	h.freeBlocks(blocks...)
}
