package blake3

import (
	"github.com/epochtimeout/baselibrary/bin128"
	"github.com/zeebo/blake3"
)

// SumB128 returns a Blake3 hash truncated from 64 bytes to B128.
func SumB128(b []byte) bin128.B128 {
	out := blake3.Sum512(b)
	sum := bin128.B128{}
	copy(sum[:], out[:16])
	return sum
}

var _ (bin128.Hash) = (*HashB128)(nil)

// HashB128 computes a Blake3 hash and truncates it from 64 bytes to B128.
type HashB128 struct {
	h *blake3.Hasher
}

// NewHashB128 returns a new hasher.
func NewHashB128() *HashB128 {
	return &HashB128{
		h: blake3.New(),
	}
}

// Write adds more data to the running hash.
// It never returns an error.
func (h *HashB128) Write(p []byte) (int, error) {
	h.h.Write(p)
	return len(p), nil
}

// SumB128 returns the current hash as bin128.
func (h *HashB128) SumB128() bin128.B128 {
	out := [64]byte{}
	h.h.Sum(out[:0])

	sum := bin128.B128{}
	copy(sum[:], out[:16])
	return sum
}

// Reset resets the hash to its initial state.
func (h HashB128) Reset() {
	h.h.Reset()
}
