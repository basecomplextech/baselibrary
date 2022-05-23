package blake3

import (
	"testing"

	"github.com/epochtimeout/library/u128"
	"github.com/stretchr/testify/assert"
	"github.com/zeebo/blake3"
)

func TestSumU128__should_compute_blake3_hash(t *testing.T) {
	b := []byte("hello, world")
	h := blake3.Sum256(b)
	u := SumU128(b)

	assert.Equal(t, h[:16], u[:])
	assert.NotEqual(t, u128.U128{}, u)
}

func TestHashU128__should_compute_blake3_hash(t *testing.T) {
	b := []byte("hello, world")
	h := blake3.Sum256(b)

	hash := NewHashU128()
	hash.Write(b)
	u := hash.SumU128()

	assert.Equal(t, h[:16], u[:])
	assert.NotEqual(t, u128.U128{}, u)
}
