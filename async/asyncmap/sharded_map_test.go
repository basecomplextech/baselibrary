// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

import (
	"testing"
	"unsafe"
)

func TestShardedMap__shard_size_should_equal_cache_line(t *testing.T) {
	s := unsafe.Sizeof(shardedMapShard[int, int]{})
	if s != 256 {
		t.Fatal(s, 256-s)
	}
}
