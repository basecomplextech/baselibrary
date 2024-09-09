// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package compactint

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Int32

func TestInt32(t *testing.T) {
	fn := func(v int32) {
		buf := make([]byte, MaxLen)
		n := PutInt32(buf, v)
		v1, n1 := Int32(buf)
		if v != v1 {
			t.Errorf("Int32(%d): got %d", v, v1)
		}
		if n != n1 {
			t.Errorf("Int32(%d): expected n=%d; n=%d", v, n, n1)
		}
	}

	tests := []int32{
		0,
		1,
		2,
		10,
		20,
		63,
		64,
		65,
		127,
		128,
		129,
		255,
		256,
		257,
		math.MaxInt16 - 1,
		math.MaxInt16,
		math.MaxInt32 - 1,
		math.MaxInt32,

		-1,
		-2,
		-255,
		-256,
		-257,
		math.MinInt16 + 1,
		math.MinInt16,
		math.MinInt32 + 1,
		math.MinInt32,
	}

	for _, v := range tests {
		fn(v)
	}
}

// Int64

func TestInt64(t *testing.T) {
	fn := func(v int64) {
		buf := make([]byte, MaxLen)
		n := PutInt64(buf, v)
		v1, n1 := Int64(buf)
		if v != v1 {
			t.Errorf("Int64(%d): got %d", v, v1)
		}
		if n != n1 {
			t.Errorf("Int64(%d): expected n=%d; n=%d", v, n, n1)
		}
	}

	tests := []int64{
		0,
		1,
		2,
		10,
		20,
		63,
		64,
		65,
		127,
		128,
		129,
		255,
		256,
		257,
		math.MaxInt16 - 1,
		math.MaxInt16,
		math.MaxInt32 - 1,
		math.MaxInt32,
		math.MaxInt64 - 1,
		math.MaxInt64,

		-1,
		-2,
		-255,
		-256,
		-257,
		math.MinInt16 + 1,
		math.MinInt16,
		math.MinInt32 + 1,
		math.MinInt32,
		math.MinInt64 + 1,
		math.MinInt64,
	}

	for _, v := range tests {
		fn(v)
	}
}

// Uint32

func TestUint32(t *testing.T) {
	fn := func(v uint32) {
		buf := make([]byte, MaxLen)
		n := PutUint32(buf, v)
		v1, n1 := Uint32(buf)

		if v != v1 {
			t.Errorf("Uint32(%d): got %d", v, v1)
		}
		if n != n1 {
			t.Errorf("Uint32(%d): expected n=%d; actual=%d", v, n, n1)
		}
	}

	tests := []uint32{
		0,
		1,
		2,
		10,
		20,
		63,
		64,
		65,
		127,
		128,
		129,
		255,
		256,
		257,
		math.MaxUint16 - 1,
		math.MaxUint16,
		math.MaxUint32 - 1,
		math.MaxUint32,
	}

	for _, v := range tests {
		fn(v)
	}
}

func TestUint32__should_return_minus_one_on_overflow(t *testing.T) {
	b := make([]byte, MaxLen)
	n := PutUint64(b, math.MaxUint64)
	b = b[:n]

	v, n := Uint32(b)
	assert.Equal(t, uint32(0), v)
	assert.Equal(t, -1, n)
}

// Uint64

func TestUint64(t *testing.T) {
	fn := func(v uint64) {
		buf := make([]byte, MaxLen)
		n := PutUint64(buf, v)
		v1, n1 := Uint64(buf)

		if v != v1 {
			t.Errorf("Uint64(%d): got %d", v, v1)
		}
		if n != n1 {
			t.Errorf("Uint64(%d): expected n=%d; actual=%d", v, n, n1)
		}
	}

	tests := []uint64{
		0,
		1,
		2,
		10,
		20,
		63,
		64,
		65,
		127,
		128,
		129,
		255,
		256,
		257,
		math.MaxUint16 - 1,
		math.MaxUint16,
		math.MaxUint32 - 1,
		math.MaxUint32,
		math.MaxUint64 - 1,
		math.MaxUint64,
	}

	for _, v := range tests {
		fn(v)
	}
}

func TestUint64__should_return_n_zero_on_small_buffer(t *testing.T) {
	b := []byte{}
	v, n := Uint64(b)
	assert.Equal(t, uint64(0), v)
	assert.Equal(t, 0, n)
}
