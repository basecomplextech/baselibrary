// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package hashing

type Hasher interface {
	// Hash32 returns a 32-bit hash of the object.
	Hash32() uint32
}
