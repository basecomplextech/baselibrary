package heap

import (
	"math/bits"
	"sync"
)

const (
	minIndex = 10 // 1024
	maxIndex = 27 // 128MB
)

type pools [maxIndex + 1]*sync.Pool

func newPools() (pools pools) {
	for j := minIndex; j <= maxIndex; j++ {
		size := 1 << j
		pool := newPool(size)
		pools[j] = pool
	}
	return
}

func newPool(size int) *sync.Pool {
	return &sync.Pool{
		New: func() any {
			return newBlock(size)
		},
	}
}

func blockPool(size int) int {
	if size == 0 {
		return 0
	}
	if isPowerOfTwo(size) {
		return poolIndex(size)
	}
	return poolIndex(size) + 1
}

func poolIndex(size int) int {
	return bits.Len(uint(size)) - 1
}

func isPowerOfTwo(size int) bool {
	return (size & (-size)) == size
}
