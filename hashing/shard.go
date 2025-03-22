// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package hashing

// Shard returns a shard index for a key, panics if the key type is not supported.
// Add more types as needed.
func Shard[K any](key K, shards int) int {
	h := Hash(key)
	return int(h) % shards
}
