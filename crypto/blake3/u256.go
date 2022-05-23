package blake3

import (
	"github.com/epochtimeout/basekit/u256"
	"github.com/zeebo/blake3"
)

// SumU256 returns a Blake3 hash truncated from 64 bytes to U256.
func SumU256(b []byte) u256.U256 {
	out := blake3.Sum256(b)
	return (u256.U256)(out)
}

var _ (u256.Hash) = (*HashU256)(nil)

// HashU256 computes a Blake3 hash and truncates it from 64 bytes to U256.
type HashU256 struct {
	h *blake3.Hasher
}

// NewHashU256 returns a new U256 hash.
func NewHashU256() *HashU256 {
	return &HashU256{
		h: blake3.New(),
	}
}

// Write adds more data to the running hash.
// It never returns an error.
func (h *HashU256) Write(p []byte) (int, error) {
	h.h.Write(p)
	return len(p), nil
}

// SumU256 returns a hash as U256.
func (h *HashU256) SumU256() u256.U256 {
	out := [64]byte{}
	h.h.Sum(out[:0])

	u := u256.U256{}
	copy(u[:], out[:])
	return u
}

// Reset resets the hash to its initial state.
func (h *HashU256) Reset() {
	h.h.Reset()
}
