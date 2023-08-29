package heap

var global = New()

type Heap struct {
	pools pools
}

// New returns a new heap.
func New() *Heap {
	return &Heap{
		pools: newPools(),
	}
}

// Alloc allocates a new block.
func (h *Heap) Alloc(size int) *Block {
	i := blockPool(size)
	if i > maxIndex {
		return newBlock(size)
	}

	if i < minIndex {
		i = minIndex
	}

	pool := h.pools[i]
	block := pool.Get().(*Block)
	return block
}

// Free frees a block.
func (h *Heap) Free(b *Block) {
	cp := cap(b.buf)
	if !isPowerOfTwo(cp) {
		return
	}

	i := blockPool(cp)
	if i < minIndex || i > maxIndex {
		return
	}

	b.reset()
	pool := h.pools[i]
	pool.Put(b)
}

// FreeMany frees multiple blocks.
func (h *Heap) FreeMany(blocks ...*Block) {
	for _, block := range blocks {
		h.Free(block)
	}
}
