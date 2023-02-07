package alloc

import (
	"math"
	"testing"

	"github.com/complex1tech/baselibrary/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testArena returns a new arena with a new heap, for tests only.
func testArena() *arena {
	a := newAllocator()
	return newArena(a)
}

// Alloc

func TestArenaAlloc__should_allocate_value(t *testing.T) {
	a := testArena()

	i := ArenaAlloc[int64](a)
	*i = math.MaxInt64

	assert.Equal(t, int64(math.MaxInt64), *i)
}

func TestArenaAlloc__should_allocate_struct(t *testing.T) {
	type Struct struct {
		Int8  int8
		Int16 int16
		Int32 int32
		Int64 int64
	}

	a := testArena()
	s := ArenaAlloc[Struct](a)

	s.Int8 = math.MaxInt8
	s.Int16 = math.MaxInt16
	s.Int32 = math.MaxInt32
	s.Int64 = math.MaxInt64
}

// ArenaSlice

func TestArenaSlice__should_allocate_slice(t *testing.T) {
	a := testArena()

	s := ArenaSlice[int](a, 16)
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

// ArenaCopy

func TestArenaCopy__should_copy_existing_slice_into_arena(t *testing.T) {
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

	a := testArena()
	v1 := ArenaCopy(a, v0)

	assert.Equal(t, v0, v1)
}

// Bytes

func TestArena_Bytes__should_allocate_bytes(t *testing.T) {
	a := testArena()
	buf := a.Bytes(16)

	for i := 0; i < len(buf); i++ {
		buf[i] = byte(i)
	}

	assert.Equal(t, 16, len(buf))
	assert.Equal(t, 16, cap(buf))
}

func TestArenaCopyBytes__should_allocate_bytes_copy(t *testing.T) {
	a := testArena()
	b := []byte("hello, world")
	buf := ArenaCopyBytes(a, b)

	assert.Equal(t, types.BytesView(b), buf)
}

// String

func TestArenaString__should_return_string_copy(t *testing.T) {
	a := testArena()
	s0 := "hello, world"
	s1 := ArenaString(a, s0)

	assert.Equal(t, types.StringView(s0), s1)
	assert.NotSame(t, s0, s1)
}

// Free

func TestArena_Free__should_release_blocks(t *testing.T) {
	a := testArena()
	a.alloc(1)

	last := a.last()
	require.Equal(t, 8, len(last.buf))

	a.Free()
	assert.Len(t, a.blocks, 0)
	assert.Equal(t, 0, len(last.buf))
}

// Len

func TestArena_Len__should_return_allocated_memory(t *testing.T) {
	a := testArena()
	a.alloc(8)
	a.alloc(32)

	ln := a.Len()
	assert.Equal(t, int64(40), ln)
}

// alloc

func TestArena_alloc__should_allocate_data(t *testing.T) {
	a := testArena()
	b := a.alloc(8)

	v := (*int64)(b)
	*v = math.MaxInt64

	assert.Equal(t, int64(math.MaxInt64), *v)
}

func TestArena_alloc__should_add_padding_for_alignment(t *testing.T) {
	a := testArena()
	a.alloc(3)

	block := a.last()
	assert.Equal(t, 8, len(block.buf))

	a.alloc(9)
	assert.Equal(t, 24, len(block.buf))
}

func TestArena_alloc__should_not_add_padding_when_already_aligned(t *testing.T) {
	a := testArena()
	a.alloc(8)

	block := a.last()
	assert.Equal(t, 8, len(block.buf))
}

func TestArena_alloc__should_allocate_next_block_when_not_enough_space(t *testing.T) {
	a := testArena()
	a.alloc(1)

	n := a.last().cap()
	a.alloc(n)

	last := a.last()
	assert.Len(t, a.blocks, 2)
	assert.Equal(t, len(last.buf), n)
}

// allocBlock

func TestArena_allocBlock__should_increment_size(t *testing.T) {
	a := testArena()
	a.alloc(1)
	size := a.last().cap()
	assert.Equal(t, int64(size), a.cap)

	a.allocBlock(1)
	size += a.last().cap()
	assert.Equal(t, int64(size), a.cap)
}
