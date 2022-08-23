package bin

import "io"

// Hash256 is the common interface for hash functions which return hashes as bin256.
type Hash256 interface {
	io.Writer

	// SumBin256 returns the current hash as bin256.
	SumBin256() Bin256

	// Reset resets the hash to its initial state.
	Reset()
}
