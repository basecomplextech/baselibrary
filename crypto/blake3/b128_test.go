package blake3

import (
	"testing"

	"github.com/epochtimeout/baselibrary/bin128"
	"github.com/stretchr/testify/assert"
	"github.com/zeebo/blake3"
)

func TestSumB128__should_compute_blake3_hash(t *testing.T) {
	b := []byte("hello, world")
	h := blake3.Sum256(b)
	u := SumB128(b)

	assert.Equal(t, h[:16], u[:])
	assert.NotEqual(t, bin128.B128{}, u)
}

func TestHashB128__should_compute_blake3_hash(t *testing.T) {
	b := []byte("hello, world")
	h := blake3.Sum256(b)

	hash := NewHashB128()
	hash.Write(b)
	u := hash.SumB128()

	assert.Equal(t, h[:16], u[:])
	assert.NotEqual(t, bin128.B128{}, u)
}
