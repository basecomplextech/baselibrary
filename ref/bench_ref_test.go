// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"testing"
)

func BenchmarkR(b *testing.B) {
	r := NewNoop(1)
	b.SetParallelism(10)

	for i := 0; i < b.N; i++ {
		r.Retain()
		r.Release()
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkR_Parallel(b *testing.B) {
	r := NewNoop(1)
	b.SetParallelism(10)

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			r.Retain()
			r.Release()
		}
	})

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}
