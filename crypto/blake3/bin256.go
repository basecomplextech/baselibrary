package blake3

import (
	"github.com/complex1tech/baselibrary/basic"
	"github.com/zeebo/blake3"
)

// SumBin256 returns a Blake3 hash truncated from 64 bytes to Bin256.
func SumBin256(b []byte) basic.Bin256 {
	out := blake3.Sum256(b)
	return (basic.Bin256)(out)
}

var _ (basic.HashBin256) = (*HashBin256)(nil)

// HashBin256 computes a Blake3 hash and truncates it from 64 bytes to Bin256.
type HashBin256 struct {
	h *blake3.Hasher
}

// NewHashBin256 returns a new Bin256 hash.
func NewHashBin256() *HashBin256 {
	return &HashBin256{
		h: blake3.New(),
	}
}

// Write adds more data to the running hash.
// It never returns an error.
func (h *HashBin256) Write(p []byte) (int, error) {
	h.h.Write(p)
	return len(p), nil
}

// SumBin256 returns a hash as Bin256.
func (h *HashBin256) SumBin256() basic.Bin256 {
	out := [64]byte{}
	h.h.Sum(out[:0])

	u := basic.Bin256{}
	copy(u[:], out[:])
	return u
}

// Reset resets the hash to its initial state.
func (h *HashBin256) Reset() {
	h.h.Reset()
}
