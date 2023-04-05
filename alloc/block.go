package alloc

const (
	minBlockSize = 1 << 10 // 1024
	maxBlockSize = 1 << 26 // 64MB
)

var blockClassSizes []int // powers of two

type block struct {
	buf []byte
}

func newBlock(size int) *block {
	return &block{
		buf: make([]byte, 0, size),
	}
}

// len returns the number of used bytes in the block.
func (b *block) len() int {
	return len(b.buf)
}

// cap returns the block size in bytes.
func (b *block) cap() int {
	return cap(b.buf)
}

// free returns a free space in bytes.
func (b *block) free() int {
	return cap(b.buf) - len(b.buf)
}

// reset resets the block.
func (b *block) reset() {
	// Zero out
	for i := range b.buf {
		b.buf[i] = 0
	}

	b.buf = b.buf[:0]
}

// getBlockClass returns an index in blockClassSizes, or -1.
func getBlockClass(size int) int {
	for cls, blockSize := range blockClassSizes {
		if size > blockSize {
			continue
		}
		return cls
	}
	return -1
}

func initBlockClasses() {
	for size := minBlockSize; size <= maxBlockSize; size *= 2 {
		blockClassSizes = append(blockClassSizes, size)
	}
}
