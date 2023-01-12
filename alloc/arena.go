package alloc

import (
	"runtime"
	"sync/atomic"
	"unsafe"

	"github.com/complex1tech/baselibrary/ref"
)

// Arena allocates objects in internal byte blocks.
// It is goroutine-safe, but operations on a free arena panic.
type Arena interface {
	// Cap returns the total memory size in bytes.
	Cap() int64

	// Len returns the allocated memory size in bytes.
	Len() int64

	// Bytes allocates a byte slice with a given capacity in the arena.
	Bytes(cap int) Bytes

	// Internal

	// Free frees the arena and releases its memory.
	Free()
}

// NewArena returns a new arena with the global allocator.
func NewArena() *ref.R[Arena] {
	a := newArena(global)
	return ref.Wrap[Arena](a)
}

// ArenaAlloc allocates a new object and returns a pointer to it.
//
// Usage:
//
//	var foo *float64
//	var bar *MyStruct
//	foo = ArenaAlloc[float64](arena)
//	bar = ArenaAlloc[MyStruct](arena)
func ArenaAlloc[T any](a Arena) *T {
	var zero T
	size := int(unsafe.Sizeof(zero))

	arena := a.(*arena)
	ptr := arena.alloc(size)
	return (*T)(ptr)
}

// ArenaBytes allocates a new byte slice.
func ArenaBytes(a Arena, cap int) Bytes {
	if cap == 0 {
		return nil
	}

	arena := a.(*arena)
	ptr := arena.alloc(cap)
	return unsafe.Slice((*byte)(ptr), cap)
}

// ArenaSlice allocates a new slice of a generic type.
//
// Usage:
//
//	var foo []MyStruct
//	foo = ArenaSlice[MyStruct](arena, 16)
func ArenaSlice[T any](a Arena, cap int) []T {
	if cap == 0 {
		return nil
	}

	var zero T
	elem := int(unsafe.Sizeof(zero))
	size := elem * cap

	arena := a.(*arena)
	ptr := arena.alloc(size)
	return unsafe.Slice((*T)(ptr), cap)
}

// ArenaCopy allocates a new slice and copies items from src into it.
// The slice capacity is len(src).
func ArenaCopy[T any](a Arena, src []T) []T {
	dst := ArenaSlice[T](a, len(src))
	copy(dst, src)
	return dst
}

// ArenaCopyBytes allocates a new byte slice and copies items from src into it.
// The slice capacity is len(src).
func ArenaCopyBytes(a Arena, src []byte) Bytes {
	dst := ArenaBytes(a, len(src))
	copy(dst, src)
	return dst
}

// ArenaString allocates a new string and copies data from src into it.
func ArenaString(a Arena, src string) String {
	if len(src) == 0 {
		return ""
	}

	dst := ArenaBytes(a, len(src))
	copy(dst, src)
	return *(*String)(unsafe.Pointer(&dst))
}

// ArenaFreeList returns a new free list which allocates objects in the given arena.
func ArenaFreeList[T any](a Arena) FreeList[T] {
	arena := a.(*arena)
	return newArenaList[T](arena)
}

// internal

type arena struct {
	a *allocator

	spinlock int32
	free     bool
	cap      int64 // total allocated size

	blocks []*block
}

// newArena returns a new arena with the given allocator.
func newArena(a *allocator) *arena {
	return &arena{a: a}
}

// Cap returns the total memory size in bytes.
func (a *arena) Cap() int64 {
	a.lock()
	defer a.unlock()

	return a.cap
}

// Len returns the allocated memory size in bytes.
func (a *arena) Len() int64 {
	a.lock()
	defer a.unlock()

	total := int64(0)
	for _, block := range a.blocks {
		total += int64(block.len())
	}
	return total
}

// Bytes allocates a byte slice with a `size` capacity in the arena.
func (a *arena) Bytes(cap int) Bytes {
	if cap == 0 {
		return nil
	}

	ptr := a.alloc(cap)
	return unsafe.Slice((*byte)(ptr), cap)
}

// Internal

// Free frees the arena and releases its memory.
func (a *arena) Free() {
	a.lock()
	defer a.unlock()

	if a.free {
		return
	}

	a.free = true
	a.cap = 0

	blocks := a.blocks
	a.blocks = nil

	a.a.freeBlocks(blocks...)
}

// alloc

const arenaAlignment = 8

// alloc allocates data and returns a pointer to it.
func (a *arena) alloc(n int) unsafe.Pointer {
	a.lock()
	defer a.unlock()

	if a.free {
		panic("operation on a free arena")
	}

	// maybe allocate block
	last := a.last()
	if last == nil || last.free() < n {
		last = a.allocBlock(n)
	}

	// grow buffer, add end padding for alignement
	start := len(last.buf)
	end := start + n
	if end > last.cap() {
		panic("allocation out of block range") // unreachable
	}

	end += (arenaAlignment - (end % arenaAlignment)) % arenaAlignment
	if end > last.cap() {
		end = last.cap()
	}

	last.buf = last.buf[:end]

	// slice buffer
	p := last.buf[start:end:end] // start:end:max, cap=max-start
	ptr := unsafe.Pointer(&p[0])
	return ptr
}

func (a *arena) allocBlock(n int) *block {
	// double last block size
	// limit it to maxBlockSize
	size := 0
	last := a.last()

	if last != nil {
		size = last.cap() * 2
	}
	if size > maxBlockSize {
		size = maxBlockSize
	}
	if n > size {
		size = n
	}

	block, _ := a.a.allocBlock(size)
	a.blocks = append(a.blocks, block)
	a.cap += int64(block.cap())
	return block
}

func (a *arena) last() *block {
	if len(a.blocks) == 0 {
		return nil
	}
	return a.blocks[len(a.blocks)-1]
}

// lock

func (a *arena) lock() {
	for {
		if atomic.CompareAndSwapInt32(&a.spinlock, 0, 1) {
			return
		}
		runtime.Gosched()
	}
}

func (a *arena) unlock() {
	atomic.StoreInt32(&a.spinlock, 0)
}
