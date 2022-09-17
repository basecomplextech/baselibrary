package bin

import "io"

// HashBin128 is a common interface for hash functions which return hashes as bin128.
type HashBin128 interface {
	io.Writer

	// SumBin128 returns the current hash as bin128.
	SumBin128() Bin128

	// Reset resets the hash to its initial state.
	Reset()
}
