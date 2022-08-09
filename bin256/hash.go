package bin256

import "io"

// Hash is the common interface for hash functions which return hashes as B256.
type Hash interface {
	io.Writer

	// SumB256 returns the current hash as bin256.
	SumB256() B256

	// Reset resets the hash to its initial state.
	Reset()
}
