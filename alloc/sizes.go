package alloc

const (
	minBlockSize = 1 << 10 // 1024
	maxBlockSize = 1 << 24 // 16MB
)

var blockClassSizes []int // powers of two

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
