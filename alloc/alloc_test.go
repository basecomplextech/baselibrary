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

func TestAllocBytes__should_allocate_byte_slice(t *testing.T) {
	a := newArena()
	buf := AllocBytes(a, 16)

	for i := 0; i < len(buf); i++ {
		buf[i] = byte(i)
	}

	assert.Equal(t, 16, len(buf))
	assert.Equal(t, 16, cap(buf))
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

// AllocString

func TestAllocString__should_return_string_copy(t *testing.T) {
	a := newArena()
	s0 := "hello, world"
	s1 := AllocString(a, s0)

	assert.Equal(t, s0, s1)
	assert.NotSame(t, s0, s1)
}
