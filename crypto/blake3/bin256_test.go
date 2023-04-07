package blake3

import (
	"testing"

	"github.com/complex1tech/baselibrary/bin"
	"github.com/stretchr/testify/assert"
	"github.com/zeebo/blake3"
)

func TestSumBin256__should_compute_blake3_hash(t *testing.T) {
	b := []byte("hello, world")
	h := blake3.Sum256(b)
	u := SumBin256(b)

	assert.Equal(t, h[:], u[:])
	assert.NotEqual(t, bin.Bin256{}, u)
}

func TestHashBin256__should_compute_blake3_hash(t *testing.T) {
	b := []byte("hello, world")
	h := blake3.Sum256(b)

	hash := NewHashBin256()
	hash.Write(b)
	u := hash.SumBin256()

	assert.Equal(t, h[:], u[:])
	assert.NotEqual(t, bin.Bin256{}, u)
}
