package alloc

import (
	"sync"
)

type Allocator interface {
	// Arena allocates a new arena.
	Arena() Arena

	// Buffer allocates a new buffer with a preallocated memory storage.
	Buffer(size int) *Buffer
}

// New returns a new allocator.
func New() Allocator {
	return newAllocator()
}

// Global returns the global allocator.
func Global() Allocator {
	return global
}

// internal

var global *allocator

type allocator struct {
	pools []*sync.Pool // match blockClassSizes
}

func newAllocator() *allocator {
	pools := make([]*sync.Pool, 0, len(blockClassSizes))
	for _, size := range blockClassSizes {
		pool := makeHeapPool(size)
		pools = append(pools, pool)
	}
	return &allocator{pools: pools}
}

// Arena allocates a new arena.
func (a *allocator) Arena() Arena {
	return newArena(a)
}

// Buffer allocates a new buffer with a preallocated memory storage.
func (a *allocator) Buffer(size int) *Buffer {
	return newBuffer(a)
}

// internal

// allocBlock allocates a block in the allocator and returns it and its class.
func (h *allocator) allocBlock(size int) (*block, int) {
	cls := getBlockClass(size)
	if cls < 0 {
		return newBlock(size), -1
	}

	pool := h.pools[cls]
	block := pool.Get().(*block)
	return block, cls
}

// freeBlock frees blocks.
func (h *allocator) freeBlocks(blocks ...*block) {
	for _, block := range blocks {
		// Get block class
		cls := getBlockClass(block.cap())
		if cls < 0 {
			continue
		}

		// Skip blocks of nonstandard sizes
		size := blockClassSizes[cls]
		if block.cap() != size {
			continue
		}

		// Reset block and release it to pool
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
