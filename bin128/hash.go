package bin128

import "io"

// Hash is the common interface for hash functions which return hashes as B128.
type Hash interface {
	io.Writer

	// SumB128 returns the current hash as bin128.
	SumB128() B128

	// Reset resets the hash to its initial state.
	Reset()
}
