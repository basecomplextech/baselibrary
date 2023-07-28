package alloc

import (
	"math"
	"testing"

	"github.com/basecomplextech/baselibrary/alloc/internal/arena"
	"github.com/stretchr/testify/assert"
)

// Alloc

func TestAlloc__should_allocate_pointer(t *testing.T) {
	a := arena.Test()

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

	a := arena.Test()
	s := Alloc[Struct](a)

	s.Int8 = math.MaxInt8
	s.Int16 = math.MaxInt16
	s.Int32 = math.MaxInt32
	s.Int64 = math.MaxInt64
}

// Slice

func TestSlice__should_allocate_slice(t *testing.T) {
	a := arena.Test()

	s := Slice[int](a, 16)
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

// Copy

func TestCopy__should_copy_existing_slice_into_arena(t *testing.T) {
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

	a := arena.Test()
	v1 := Copy(a, v0)

	assert.Equal(t, v0, v1)
}

// Bytes

func TestBytes__should_allocate_bytes(t *testing.T) {
	a := arena.Test()
	buf := Bytes(a, 16)

	for i := 0; i < len(buf); i++ {
		buf[i] = byte(i)
	}

	assert.Equal(t, 16, len(buf))
	assert.Equal(t, 16, cap(buf))
}

func TestCopyBytes__should_allocate_bytes_copy(t *testing.T) {
	a := arena.Test()
	b := []byte("hello, world")
	buf := CopyBytes(a, b)

	assert.Equal(t, b, buf)
}

// String

func TestString__should_return_string_copy(t *testing.T) {
	a := arena.Test()
	s0 := "hello, world"
	s1 := String(a, s0)

	assert.Equal(t, s0, s1)
	assert.NotSame(t, s0, s1)
}
