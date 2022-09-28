package blake3

import (
	"testing"

	"github.com/complex1tech/baselibrary/types"
	"github.com/stretchr/testify/assert"
	"github.com/zeebo/blake3"
)

func TestSumBin128__should_compute_blake3_hash(t *testing.T) {
	b := []byte("hello, world")
	h := blake3.Sum256(b)
	u := SumBin128(b)

	assert.Equal(t, h[:16], u[:])
	assert.NotEqual(t, types.Bin128{}, u)
}

func TestHashBin128__should_compute_blake3_hash(t *testing.T) {
	b := []byte("hello, world")
	h := blake3.Sum256(b)

	hash := NewHashBin128()
	hash.Write(b)
	u := hash.SumBin128()

	assert.Equal(t, h[:16], u[:])
	assert.NotEqual(t, types.Bin128{}, u)
}
