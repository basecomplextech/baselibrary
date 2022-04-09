package xxh3

import (
	"github.com/zeebo/xxh3"
)

// Sum64 returns a XXH3 hash as uint64.
func Sum64(b []byte) uint64 {
	return xxh3.Hash(b)
}

// Hash64 computes an XXH3 hash and returns it as uint64.
type Hash64 struct {
	h xxh3.Hasher
}

// Write adds more data to the running hash.
// It never returns an error.
func (h *Hash64) Write(p []byte) (int, error) {
	h.h.Write(p)
	return len(p), nil
}

// Sum64 returns the current hash as uint64.
func (h *Hash64) Sum64() uint64 {
	return h.h.Sum64()
}

// Reset resets the hash to its initial state.
func (h *Hash64) Reset() {
	h.h.Reset()
}
