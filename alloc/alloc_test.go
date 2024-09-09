// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

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
