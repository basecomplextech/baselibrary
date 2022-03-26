package blake3

import (
	"testing"

	"github.com/baseblck/library/u256"
	"github.com/stretchr/testify/assert"
	"github.com/zeebo/blake3"
)

func TestSumU256__should_compute_blake3_hash(t *testing.T) {
	b := []byte("hello, world")
	h := blake3.Sum256(b)
	u := SumU256(b)

	assert.Equal(t, h[:], u[:])
	assert.NotEqual(t, u256.U256{}, u)
}

func TestHashU256__should_compute_blake3_hash(t *testing.T) {
	b := []byte("hello, world")
	h := blake3.Sum256(b)

	hash := NewHashU256()
	hash.Write(b)
	u := hash.SumU256()

	assert.Equal(t, h[:], u[:])
	assert.NotEqual(t, u256.U256{}, u)
}
