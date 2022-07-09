package alloc

import (
	"sync"
)

var globalHeap *heap

type heap struct {
	pools []*sync.Pool // match blockClassSizes
}

func newHeap() *heap {
	pools := make([]*sync.Pool, 0, len(blockClassSizes))
	for _, size := range blockClassSizes {
		pool := makeHeapPool(size)
		pools = append(pools, pool)
	}
	return &heap{pools: pools}
}

// allocBlock allocates a block in the heap and returns it and its class.
func (h *heap) allocBlock(size int) (*block, int) {
	cls := getBlockClass(size)
	if cls < 0 {
		return newBlock(size), -1
	}

	pool := h.pools[cls]
	block := pool.Get().(*block)
	return block, cls
}

// freeBlock frees blocks.
func (h *heap) freeBlocks(blocks ...*block) {
	for _, block := range blocks {
		// get block class
		cls := getBlockClass(block.cap())
		if cls < 0 {
			continue
		}

		// skip blocks of nonstandard sizes
		size := blockClassSizes[cls]
		if block.cap() != size {
			continue
		}

		// reset block and release it to pool
		block.reset()
		pool := h.pools[cls]
		pool.Put(block)
	}
}

func makeHeapPool(blockSize int) *sync.Pool {
	return &sync.Pool{
		New: func() any {
			return newBlock(blockSize)
		},
	}
}

func initGlobalHeap() {
	globalHeap = newHeap()
}
