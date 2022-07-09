package u128

import "io"

// Hash is the common interface for hash functions which return hashes as U128.
type Hash interface {
	io.Writer

	// SumU128 returns the current hash as u128.
	SumU128() U128

	// Reset resets the hash to its initial state.
	Reset()
}
