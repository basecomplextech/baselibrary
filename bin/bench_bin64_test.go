// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package bin

import "testing"

func BenchmarkBin64_String(b *testing.B) {
	v := Random64()

	b.ReportAllocs()
	b.ResetTimer()

	var s string
	for i := 0; i < b.N; i++ {
		s = v.String()
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec

	b.ReportMetric(ops/1000_000, "mops")
	_ = s
}
