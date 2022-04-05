package compactint

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ReverseUint32

func TestReverseUint32__should_return_minus_one_on_overflow(t *testing.T) {
	b := make([]byte, MaxLen)
	n := PutReverseUint64(b, math.MaxUint64)
	off := len(b) - n

	v, n := ReverseUint32(b[off:])
	assert.Equal(t, uint32(0), v)
	assert.Equal(t, -1, n)
}

// ReverseUint64

func TestReverseUint64__should_read_write_uvarint(t *testing.T) {
	fn := func(v uint64) {
		b := make([]byte, MaxLen)
		n := PutReverseUint64(b, v)
		off := len(b) - n
		v1, n1 := ReverseUint64(b[off:])

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

func TestReverseUint64__should_return_n_zero_on_small_buffer(t *testing.T) {
	b := []byte{}
	v, n := ReverseUint64(b)
	assert.Equal(t, uint64(0), v)
	assert.Equal(t, 0, n)
}
