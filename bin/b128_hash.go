package bin

import "io"

// Hash128 is the common interface for hash functions which return hashes as bin128.
type Hash128 interface {
	io.Writer

	// SumBin128 returns the current hash as bin128.
	SumBin128() Bin128

	// Reset resets the hash to its initial state.
	Reset()
}
