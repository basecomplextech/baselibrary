package alloc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAllocator_allocBlock__should_allocate_block(t *testing.T) {
	a := newAllocator()

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
		block, cls := a.allocBlock(c.size)
		require.Equal(t, c.blockSize, block.cap())
		require.Equal(t, c.cls, cls)
	}
}

func TestAllocator_freeBlocks__should_free_blocks(t *testing.T) {
	a := newAllocator()
	sizes := []int{1, 1023, 1024, 2048}

	blocks := []*block{}
	for _, size := range sizes {
		block, _ := a.allocBlock(size)
		blocks = append(blocks, block)
	}

	a.freeBlocks(blocks...)
}
