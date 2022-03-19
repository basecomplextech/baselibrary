package alloc

import (
	"fmt"
	"unsafe"
)

const alignment = 8

type block struct {
	data []byte
}

func newBlock(size int) *block {
	return &block{
		data: make([]byte, 0, size),
	}
}

// allocated returns the number of allocated bytes in the block.
func (b *block) allocated() int {
	return len(b.data)
}

// free returns a free space in bytes.
func (b *block) free() int {
	return cap(b.data) - len(b.data)
}

// size returns the block size in bytes.
func (b *block) size() int {
	return cap(b.data)
}

// alloc returns a pointer to memory of size.
func (b *block) alloc(size int) unsafe.Pointer {
	free := cap(b.data) - len(b.data)
	if free < size {
		panic(fmt.Sprintf("block out of memory, free=%d, requested size=%d", free, size))
	}

	// calc range
	start := len(b.data)
	end := start + size

	// grow buffer, add padding
	pad := (alignment - (end % alignment)) % alignment
	ln := end + pad
	if ln > cap(b.data) {
		ln = cap(b.data)
	}
	b.data = b.data[:ln]

	// slice buffer
	out := b.data[start:end:end] // start:end:max (cap=max-start)

	// zero out
	for i := range out {
		out[i] = 0
	}
	return unsafe.Pointer(&out[0])
}

// reset resets the block.
func (b *block) reset() {
	b.data = b.data[:0]
}
