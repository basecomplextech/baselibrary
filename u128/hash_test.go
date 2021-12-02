package u128

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeebo/blake3"
)

func TestHash__should_calc_blake3_hash(t *testing.T) {
	b := []byte("hello, world")
	h := blake3.Sum256(b)
	u := Hash(b)

	assert.Equal(t, h[:16], u[:])
	assert.NotEqual(t, U128{}, u)
}

func TestHasher__should_calc_blake3_hash(t *testing.T) {
	b := []byte("hello, world")
	h := blake3.Sum256(b)

	hasher := NewHasher()
	hasher.Write(b)
	u := hasher.Sum()

	assert.Equal(t, h[:16], u[:])
	assert.NotEqual(t, U128{}, u)
}
