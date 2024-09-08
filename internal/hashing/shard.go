// Copyright 2024 Ivan Korobkov. All rights reserved.

package hashing

// Shard returns a shard index for a key, panics if the key type is not supported.
// Add more types as needed.
func Shard[K any](key K, shards int) int {
	h := Hash(key)
	return int(h) % shards
}
