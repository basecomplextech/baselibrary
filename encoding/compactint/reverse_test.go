// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package compactint

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ReverseInt32

func TestReverseInt32(t *testing.T) {
	fn := func(v int32) {
		b := make([]byte, MaxLen32)
		n := PutReverseInt32(b, v)
		off := len(b) - n
		v1, n1 := ReverseInt32(b[off:])

		if v != v1 {
			t.Errorf("ReverseInt32(%d): got %d", v, v1)
		}
		if n != n1 {
			t.Errorf("ReverseInt32(%d): expected n=%d; actual=%d", v, n, n1)
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

// ReverseInt64

func TestReverseInt64(t *testing.T) {
	fn := func(v int64) {
		b := make([]byte, MaxLen64)
		n := PutReverseInt64(b, v)
		off := len(b) - n
		v1, n1 := ReverseInt64(b[off:])

		if v != v1 {
			t.Errorf("ReverseInt64(%d): got %d", v, v1)
		}
		if n != n1 {
			t.Errorf("ReverseInt64(%d): expected n=%d; actual=%d", v, n, n1)
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

// ReverseUint32

func TestReverseUint32(t *testing.T) {
	fn := func(v uint32) {
		b := make([]byte, MaxLen32)
		n := PutReverseUint32(b, v)
		off := len(b) - n
		v1, n1 := ReverseUint32(b[off:])

		if v != v1 {
			t.Errorf("ReverseUint32(%d): got %d", v, v1)
		}
		if n != n1 {
			t.Errorf("ReverseUint32(%d): expected n=%d; actual=%d", v, n, n1)
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

func TestReverseUint32__should_return_minus_one_on_overflow(t *testing.T) {
	b := make([]byte, MaxLen64)
	n := PutReverseUint64(b, math.MaxUint64)
	off := len(b) - n

	v, n := ReverseUint32(b[off:])
	assert.Equal(t, uint32(0), v)
	assert.Equal(t, -1, n)
}

// ReverseUint64

func TestReverseUint64(t *testing.T) {
	fn := func(v uint64) {
		b := make([]byte, MaxLen64)
		n := PutReverseUint64(b, v)
		off := len(b) - n
		v1, n1 := ReverseUint64(b[off:])

		if v != v1 {
			t.Errorf("ReverseUint64(%d): got %d", v, v1)
		}
		if n != n1 {
			t.Errorf("ReverseUint64(%d): expected n=%d; actual=%d", v, n, n1)
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

func TestReverseUint64__should_return_n_zero_on_small_buffer(t *testing.T) {
	b := []byte{}
	v, n := ReverseUint64(b)
	assert.Equal(t, uint64(0), v)
	assert.Equal(t, 0, n)
}
