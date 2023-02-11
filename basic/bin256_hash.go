package basic

import "io"

// HashBin256 is a common interface for hash functions which return hashes as bin256.
type HashBin256 interface {
	io.Writer

	// SumBin256 returns the current hash as bin256.
	SumBin256() Bin256

	// Reset resets the hash to its initial state.
	Reset()
}
