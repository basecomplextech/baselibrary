// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package bin

import "io"

// Hash128 is a common interface for hash functions which return hashes as bin128.
type Hash128 interface {
	io.Writer

	// SumBin128 returns the current hash as bin128.
	SumBin128() Bin128

	// Reset resets the hash to its initial state.
	Reset()
}

// Hash256 is a common interface for hash functions which return hashes as bin256.
type Hash256 interface {
	io.Writer

	// SumBin256 returns the current hash as bin256.
	SumBin256() Bin256

	// Reset resets the hash to its initial state.
	Reset()
}
