package blake3

import (
	"github.com/epochtimeout/baselibrary/bin"
	"github.com/zeebo/blake3"
)

// SumBin128 returns a Blake3 hash truncated from 64 bytes to bin128.
func SumBin128(b []byte) bin.Bin128 {
	out := blake3.Sum512(b)
	sum := bin.Bin128{}
	copy(sum[:], out[:16])
	return sum
}

var _ (bin.Hash128) = (*HashBin128)(nil)

// HashBin128 computes a Blake3 hash and truncates it from 64 bytes to bin128.
type HashBin128 struct {
	h *blake3.Hasher
}

// NewHashBin128 returns a new hasher.
func NewHashBin128() *HashBin128 {
	return &HashBin128{
		h: blake3.New(),
	}
}

// Write adds more data to the running hash.
// It never returns an error.
func (h *HashBin128) Write(p []byte) (int, error) {
	h.h.Write(p)
	return len(p), nil
}

// SumBin128 returns the current hash as bin128.
func (h *HashBin128) SumBin128() bin.Bin128 {
	out := [64]byte{}
	h.h.Sum(out[:0])

	sum := bin.Bin128{}
	copy(sum[:], out[:16])
	return sum
}

// Reset resets the hash to its initial state.
func (h HashBin128) Reset() {
	h.h.Reset()
}
