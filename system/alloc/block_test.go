package alloc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSizeClass__should_return_block_class_for_allocation_size(t *testing.T) {
	cls := 0
	for size := 0; size <= maxBlockSize; size += 1 {
		actual := getBlockClass(size)
		if cls != actual {
			t.Fatalf("size=%d, expected class=%d, actual class=%d", size, cls, actual)
		}

		if size == blockClassSizes[cls] {
			cls++
		}
	}
}

func TestGetSizeClass__should_return_minus_one_when_no_block_class_for_size(t *testing.T) {
	cls := getBlockClass(maxBlockSize + 1)
	assert.Equal(t, -1, cls)
}
