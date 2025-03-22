// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package hashing

import (
	"testing"

	"github.com/basecomplextech/baselibrary/bin"
)

func BenchmarkShard(b *testing.B) {
	key := bin.Int128(0, 123)

	for b.Loop() {
		i := Shard(key, 10)
		if i < 0 || i >= 10 {
			b.Fatal(i)
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}
