// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package xxh3

import (
	"github.com/zeebo/xxh3"
)

// Sum32 returns a XXH3 hash high bits as uint32.
func Sum32(b []byte) uint32 {
	return uint32(xxh3.Hash(b) >> 32)
}

// Sum64 returns a XXH3 hash as uint64.
func Sum64(b []byte) uint64 {
	return xxh3.Hash(b)
}

// Hash64

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

// Sum32 returns the current hash high bits as uint32.
func (h *Hash64) Sum32() uint32 {
	return uint32(h.h.Sum64() >> 32)
}

// Sum64 returns the current hash as uint64.
func (h *Hash64) Sum64() uint64 {
	return h.h.Sum64()
}

// Reset resets the hash to its initial state.
func (h *Hash64) Reset() {
	h.h.Reset()
}
