package u256

import (
	"crypto/sha512"
	"hash"
	"io"
)

// Hasher computes a U256 id.
type Hasher interface {
	io.Writer

	// Sum returns the current hash.
	Sum() U256
}

// NewHasher returns a new hasher which computes a U256 id.
func NewHasher() Hasher {
	return newHasher()
}

var _ (Hasher) = (*hasher)(nil)

type hasher struct {
	sha512_256 hash.Hash
}

func newHasher() *hasher {
	return &hasher{
		sha512_256: sha512.New512_256(),
	}
}

func (h *hasher) Write(p []byte) (n int, err error) {
	return h.sha512_256.Write(p)
}

// Sum returns the current hash.
func (h *hasher) Sum() U256 {
	u := U256{}
	b := u[:]
	h.sha512_256.Sum(b[:0])
	return u
}
