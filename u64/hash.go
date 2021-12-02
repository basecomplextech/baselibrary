package u64

import "github.com/zeebo/xxh3"

// Hash returns a XXH3 hash as U64/uint64.
func Hash(b []byte) U64 {
	return xxh3.Hash(b)
}
