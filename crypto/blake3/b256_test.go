package blake3

import (
	"testing"

	"github.com/epochtimeout/baselibrary/bin256"
	"github.com/stretchr/testify/assert"
	"github.com/zeebo/blake3"
)

func TestSumB256__should_compute_blake3_hash(t *testing.T) {
	b := []byte("hello, world")
	h := blake3.Sum256(b)
	u := SumB256(b)

	assert.Equal(t, h[:], u[:])
	assert.NotEqual(t, bin256.B256{}, u)
}

func TestHashB256__should_compute_blake3_hash(t *testing.T) {
	b := []byte("hello, world")
	h := blake3.Sum256(b)

	hash := NewHashB256()
	hash.Write(b)
	u := hash.SumB256()

	assert.Equal(t, h[:], u[:])
	assert.NotEqual(t, bin256.B256{}, u)
}
