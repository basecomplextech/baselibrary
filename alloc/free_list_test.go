package alloc

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// Alloc

func TestFreeList_Alloc__should_allocate_new_object(t *testing.T) {
	a := newArena()
	list := newFreeList[int64](a)

	v := list.Alloc()
	*v = math.MaxInt64

	assert.Equal(t, int64(math.MaxInt64), *v)
}

func TestFreeList_Alloc__should_return_free_object(t *testing.T) {
	a := newArena()
	list := newFreeList[int64](a)

	v0 := list.Alloc()
	list.Dealloc(v0)

	v1 := list.Alloc()
	assert.Same(t, v0, v1)
}

func TestFreeList_Alloc__should_consume_free_item(t *testing.T) {
	a := newArena()
	list := newFreeList[int64](a)

	v0 := list.Alloc()
	list.Dealloc(v0)

	list.Alloc()
	assert.Zero(t, list.free)
}

func TestFreeList_Alloc__should_swap_free_item_with_previous(t *testing.T) {
	a := newArena()
	list := newFreeList[int64](a)

	v0 := list.Alloc()
	v1 := list.Alloc()
	list.Dealloc(v0)
	list.Dealloc(v1)

	list.Alloc()
	ptr0 := (uintptr)(unsafe.Pointer(v0))
	assert.Equal(t, ptr0, list.free)
}

func TestFreeList_Alloc__should_zero_object(t *testing.T) {
	type Value struct {
		A int64
		B int64
		C int64
	}

	a := newArena()
	list := newFreeList[Value](a)

	v := list.Alloc()
	v.A = 1
	v.B = 2
	v.C = 3

	list.Dealloc(v)

	v1 := list.Alloc()
	assert.Zero(t, *v1)
}

// Dealloc

func TestFreeList_Dealloc__should_swap_free_item_with_next(t *testing.T) {
	a := newArena()
	list := newFreeList[int64](a)

	v0 := list.Alloc()
	v1 := list.Alloc()

	list.Dealloc(v0)
	ptr0 := (uintptr)(unsafe.Pointer(v0))
	assert.Equal(t, ptr0, list.free)

	list.Dealloc(v1)
	ptr1 := (uintptr)(unsafe.Pointer(v1))
	assert.Equal(t, ptr1, list.free)
}
