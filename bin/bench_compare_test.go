// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package bin

import "testing"

func BenchmarkCompare64(b *testing.B) {
	v := Random64()
	v1 := Random64()

	b.ReportAllocs()
	b.ResetTimer()

	var c int
	for i := 0; i < b.N; i++ {
		c = v.Compare(v1)
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec

	b.ReportMetric(ops/1000_000, "mops")
	_ = c
}

func BenchmarkCompare256(b *testing.B) {
	v := Random256()
	v1 := Random256()

	b.ReportAllocs()
	b.ResetTimer()

	var c int
	for i := 0; i < b.N; i++ {
		c = v.Compare(v1)
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec

	b.ReportMetric(ops/1000_000, "mops")
	_ = c
}
