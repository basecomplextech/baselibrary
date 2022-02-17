package u256

import "io"

// Hash is the common interface for hash functions which return hashes as U256.
type Hash interface {
	io.Writer

	// SumU256 returns the current hash as u256.
	SumU256() U256

	// Reset resets the hash to its initial state.
	Reset()
}
