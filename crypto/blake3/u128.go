package blake3

import (
	"github.com/epochtimeout/library/u128"
	"github.com/zeebo/blake3"
)

// SumU128 returns a Blake3 hash truncated from 64 bytes to U128.
func SumU128(b []byte) u128.U128 {
	out := blake3.Sum512(b)
	sum := u128.U128{}
	copy(sum[:], out[:16])
	return sum
}

var _ (u128.Hash) = (*HashU128)(nil)

// HashU128 computes a Blake3 hash and truncates it from 64 bytes to U128.
type HashU128 struct {
	h *blake3.Hasher
}

// NewHashU128 returns a new hasher.
func NewHashU128() *HashU128 {
	return &HashU128{
		h: blake3.New(),
	}
}

// Write adds more data to the running hash.
// It never returns an error.
func (h *HashU128) Write(p []byte) (int, error) {
	h.h.Write(p)
	return len(p), nil
}

// SumU128 returns the current hash as u128.
func (h *HashU128) SumU128() u128.U128 {
	out := [64]byte{}
	h.h.Sum(out[:0])

	sum := u128.U128{}
	copy(sum[:], out[:16])
	return sum
}

// Reset resets the hash to its initial state.
func (h HashU128) Reset() {
	h.h.Reset()
}
