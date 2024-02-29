package arena

import (
	"sync"
	"unsafe"

	"github.com/basecomplextech/baselibrary/alloc/internal/heap"
	"github.com/basecomplextech/baselibrary/collect/sets"
	"github.com/basecomplextech/baselibrary/collect/slices"
)

// Arena is an arena allocator, which internally allocates memory in blocks.
//
// Arena is not thread-safe. If you need to use it in multiple goroutines,
// you must synchronize access or you may consider adding an AtomicArena wrapper.
type Arena interface {
	// Cap returns the arena capacity.
	Cap() int64

	// Len calculates and returns the arena used size.
	Len() int64

	// Methods

	// Alloc allocates a memory block and returns a pointer to it.
	Alloc(size int) unsafe.Pointer

	// Pin pins an external object to the arena.
	// The method is used to prevent the object from being collected by the garbage collector.
	Pin(obj any)

	// Reset resets the arena.
	Reset()

	// Internal

	// Free frees the arena and releases its memory.
	// The method is not thread-safe and must be called only once.
	Free()
}

// New returns an empty arena.
func New() Arena {
	return newArena(heap.Global)
}

// NewSize returns an arena with an initial capacity.
func NewSize(size int) Arena {
	return newArenaSize(heap.Global, size)
}

// AcquireArena returns a pooled arena, which is released to the pool on Free.
//
// The arena must not be used or even referenced after Free.
// Use these method only when arenas do not escape an isolated scope.
//
// Typical usage:
//
//	arena := alloc.AcquireArena()
//	defer arena.Free() // free immediately
func AcquireArena() Arena {
	return acquireArena()
}

// internal

type arena struct {
	*state
}

type state struct {
	heap   *heap.Heap
	pooled bool

	init int   // initial capacity
	cap  int64 // total allocated capacity

	blocks []*heap.Block
	pinned sets.Set[any]
}

func newArena(h *heap.Heap) *arena {
	return newArenaSize(h, heap.MinBlockSize)
}

func newArenaSize(heap *heap.Heap, size int) *arena {
	a := &arena{acquireState()}
	a.heap = heap

	if size > 0 {
		a.init = a.allocBlock(size).Cap()
	}
	return a
}

// Cap returns the arena capacity.
func (a *arena) Cap() int64 {
	return a.cap
}

// Len calculates and returns the arena used size.
func (a *arena) Len() int64 {
	n := int64(0)
	for _, block := range a.blocks {
		n += int64(block.Len())
	}
	return n
}

// Alloc allocates a memory block and returns a pointer to it.
func (a *arena) Alloc(size int) unsafe.Pointer {
	if len(a.blocks) > 0 {
		b := a.blocks[len(a.blocks)-1]

		ptr := b.Alloc(size)
		if ptr != nil {
			return ptr
		}
	}

	b := a.allocBlock(size)
	return b.Alloc(size)
}

// Pin pins an external object to the arena.
// The method is used to prevent the object from being collected by the garbage collector.
func (a *arena) Pin(obj any) {
	if a.pinned == nil {
		a.pinned = sets.New[any]()
	}

	a.pinned.Add(obj)
}

// Reset resets the arena.
func (a *arena) Reset() {
	a.reset()
}

// Internal

// Free frees the arena and releases its memory.
// The method is not thread-safe and must be called only once.
func (a *arena) Free() {
	if a.pooled {
		releaseArena(a)
		return
	}

	a.free()
	s := a.state
	a.state = nil
	releaseState(s)
}

// private

func (a *arena) free() {
	a.heap.FreeMany(a.blocks...)

	a.cap = 0
	a.blocks = slices.Clear(a.blocks)

	if a.pinned != nil {
		clear(a.pinned)
	}
}

func (a *arena) reset() {
	// Clear pinned objects
	if a.pinned != nil {
		clear(a.pinned)
	}

	// Reset blocks
	if len(a.blocks) == 0 {
		return
	}

	// Maybe just reset the first block
	n := 0
	if b := a.blocks[0]; b.Cap() == a.init {
		n = 1

		b.Reset()
		a.cap = int64(b.Cap())

		if len(a.blocks) == 1 {
			return
		}
	}

	// Free other blocks
	a.heap.FreeMany(a.blocks[n:]...)
	slices.Zero(a.blocks[n:]) // for gc

	a.cap = 0
	a.blocks = a.blocks[:n]
}

// blocks

func (a *arena) allocBlock(n int) *heap.Block {
	// Double last block capacity
	size := 0
	if len(a.blocks) > 0 {
		last := a.blocks[len(a.blocks)-1]
		size = last.Cap() * 2
	}
	if n > size {
		size = n
	}

	// Alloc next block
	b := a.heap.Alloc(size)
	a.blocks = append(a.blocks, b)
	a.cap += int64(b.Cap())
	return b
}

// arena pool

var arenaPool = &sync.Pool{
	New: func() any {
		a := newArena(heap.Global)
		a.pooled = true
		return a
	},
}

func acquireArena() *arena {
	return arenaPool.Get().(*arena)
}

func releaseArena(a *arena) {
	a.reset()
	arenaPool.Put(a)
}

// state pool

var statePool = &sync.Pool{}

func acquireState() *state {
	v := statePool.Get()
	if v != nil {
		return v.(*state)
	}
	return &state{}
}

func releaseState(s *state) {
	s.reset()
	statePool.Put(s)
}

func (s *state) reset() {
	blocks := slices.Clear(s.blocks)

	*s = state{}
	s.blocks = blocks
}
