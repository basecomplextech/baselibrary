package alloc

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlloc__should_allocate_value(t *testing.T) {
	a := newArena()

	i := Alloc[int64](a)
	*i = math.MaxInt64

	assert.Equal(t, int64(math.MaxInt64), *i)
}

func TestAlloc__should_allocate_struct(t *testing.T) {
	type Struct struct {
		Int8  int8
		Int16 int16
		Int32 int32
		Int64 int64
	}

	a := newArena()
	s := Alloc[Struct](a)

	s.Int8 = math.MaxInt8
	s.Int16 = math.MaxInt16
	s.Int32 = math.MaxInt32
	s.Int64 = math.MaxInt64
}

// AllocBytes

func TestAllocBytes__should_allocate_bytes(t *testing.T) {
	a := newArena()
	buf := AllocBytes(a, 16)

	for i := 0; i < len(buf); i++ {
		buf[i] = byte(i)
	}

	assert.Equal(t, 16, len(buf))
	assert.Equal(t, 16, cap(buf))
}

func TestCopyBytes__should_allocate_bytes_copy(t *testing.T) {
	a := newArena()
	b := []byte("hello, world")
	buf := CopyBytes(a, b)

	assert.Equal(t, b, buf)
}

// AllocSlice

func TestAllocSlice__should_allocate_slice(t *testing.T) {
	a := newArena()

	s := AllocSlice[int](a, 16)
	s[0] = 1
	s[1] = 2
	s[2] = 3
	s[3] = 4
	s[4] = 5

	assert.Equal(t, 1, s[0])
	assert.Equal(t, 2, s[1])
	assert.Equal(t, 3, s[2])
	assert.Equal(t, 4, s[3])
	assert.Equal(t, 5, s[4])

	assert.Equal(t, 16, len(s))
	assert.Equal(t, 16, cap(s))
}

// CopySlice

func TestCopySlice__should_copy_existing_slice_into_arena(t *testing.T) {
	type Value struct {
		A int64
		B int64
		C int64
	}

	v0 := []Value{
		{1, 2, 3},
		{10, 20, 30},
		{100, 200, 300},
	}

	a := newArena()
	v1 := CopySlice(a, v0)

	assert.Equal(t, v0, v1)
}

// AllocString

func TestAllocString__should_return_string_copy(t *testing.T) {
	a := newArena()
	s0 := "hello, world"
	s1 := AllocString(a, s0)

	assert.Equal(t, s0, s1)
	assert.NotSame(t, s0, s1)
}

// AllocFreeList

func TestAllocFreeList__should_allocate_free_list(t *testing.T) {
	a := newArena()
	list := AllocFreeList[int64](a)

	v0 := list.Get()
	*v0 = math.MaxInt64
	list.Put(v0)

	v1 := list.Get()
	assert.Zero(t, *v1)
}

func TestAllocFreeList__should_return_same_free_list_for_same_type(t *testing.T) {
	a := newArena()
	list0 := AllocFreeList[int64](a)
	list1 := allocFreeList[int64](a)

	assert.Same(t, list0, list1)
}

func TestAllocFreeList__should_return_different_lists_for_different_types_with_same_size(t *testing.T) {
	type Value struct {
		V int64
	}

	a := newArena()
	list0 := AllocFreeList[int64](a)
	list1 := allocFreeList[Value](a)

	assert.NotSame(t, list0, list1)
}
