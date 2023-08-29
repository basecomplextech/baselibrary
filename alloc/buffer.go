package alloc

import (
	"github.com/basecomplextech/baselibrary/alloc/internal/buf"
)

// Buffer is a buffer allocated by an allocator.
// The buffer must be freed after usage.
type Buffer = buf.Buffer

// NewBuffer allocates a buffer.
func NewBuffer() *Buffer {
	return buf.New(globalHeap)
}

// NewBuffer allocates a buffer of a preallocated capacity.
func NewBufferSize(size int) *Buffer {
	return buf.NewSize(globalHeap, size)
}
