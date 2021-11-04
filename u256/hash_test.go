package u256

import (
	"crypto/sha512"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSum__should_calc_sha512_256_hash(t *testing.T) {
	b := []byte("hello, world")
	h := sha512.Sum512_256(b)
	u := Sum(b)

	assert.Equal(t, h[:], u[:])
	assert.NotEqual(t, U256{}, u)
}
