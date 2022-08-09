package blake3

import (
	"github.com/epochtimeout/baselibrary/bin256"
	"github.com/zeebo/blake3"
)

// SumB256 returns a Blake3 hash truncated from 64 bytes to B256.
func SumB256(b []byte) bin256.B256 {
	out := blake3.Sum256(b)
	return (bin256.B256)(out)
}

var _ (bin256.Hash) = (*HashB256)(nil)

// HashB256 computes a Blake3 hash and truncates it from 64 bytes to B256.
type HashB256 struct {
	h *blake3.Hasher
}

// NewHashB256 returns a new B256 hash.
func NewHashB256() *HashB256 {
	return &HashB256{
		h: blake3.New(),
	}
}

// Write adds more data to the running hash.
// It never returns an error.
func (h *HashB256) Write(p []byte) (int, error) {
	h.h.Write(p)
	return len(p), nil
}

// SumB256 returns a hash as B256.
func (h *HashB256) SumB256() bin256.B256 {
	out := [64]byte{}
	h.h.Sum(out[:0])

	u := bin256.B256{}
	copy(u[:], out[:])
	return u
}

// Reset resets the hash to its initial state.
func (h *HashB256) Reset() {
	h.h.Reset()
}
