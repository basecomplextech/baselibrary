package freelist

import (
	"math"
	"testing"
	"unsafe"

	"github.com/basecomplextech/baselibrary/alloc/internal/arena"
	"github.com/stretchr/testify/assert"
)

// New

func TestNew__should_allocate_free_list(t *testing.T) {
	a := arena.Test()
	list := New[int64](a)

	v0 := list.Get()
	*v0 = math.MaxInt64
	list.Put(v0)

	v1 := list.Get()
	assert.Zero(t, *v1)
}

func TestNew__should_return_different_lists_for_different_types_with_same_size(t *testing.T) {
	type Value struct {
		V int64
	}

	a := arena.Test()
	list0 := New[int64](a)
	list1 := New[Value](a)

	assert.NotSame(t, list0, list1)
}

// Get

func TestList_Get__should_allocate_new_object(t *testing.T) {
	a := arena.Test()
	list := newList[int64](a)

	v := list.Get()
	*v = math.MaxInt64

	assert.Equal(t, int64(math.MaxInt64), *v)
}

func TestList_Get__should_return_free_object(t *testing.T) {
	a := arena.Test()
	list := newList[int64](a)

	v0 := list.Get()
	list.Put(v0)

	v1 := list.Get()
	assert.Same(t, v0, v1)
}

func TestList_Get__should_consume_free_item(t *testing.T) {
	a := arena.Test()
	list := newList[int64](a)

	v0 := list.Get()
	list.Put(v0)

	list.Get()
	assert.Zero(t, list.free)
}

func TestList_Get__should_swap_free_item_with_previous(t *testing.T) {
	a := arena.Test()
	list := newList[int64](a)

	v0 := list.Get()
	v1 := list.Get()
	list.Put(v0)
	list.Put(v1)

	list.Get()
	ptr0 := (uintptr)(unsafe.Pointer(v0))
	assert.Equal(t, ptr0, list.free)
}

func TestList_Get__should_zero_object(t *testing.T) {
	type Value struct {
		A int64
		B int64
		C int64
	}

	a := arena.Test()
	list := newList[Value](a)

	v := list.Get()
	v.A = 1
	v.B = 2
	v.C = 3

	list.Put(v)

	v1 := list.Get()
	assert.Zero(t, *v1)
}

// Put

func TestList_Put__should_swap_free_item_with_next(t *testing.T) {
	a := arena.Test()
	list := newList[int64](a)

	v0 := list.Get()
	v1 := list.Get()

	list.Put(v0)
	ptr0 := (uintptr)(unsafe.Pointer(v0))
	assert.Equal(t, ptr0, list.free)

	list.Put(v1)
	ptr1 := (uintptr)(unsafe.Pointer(v1))
	assert.Equal(t, ptr1, list.free)
}
