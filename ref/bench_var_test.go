// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import "testing"

func BenchmarkVar_Parallel(b *testing.B) {
	v := newVar[int]()
	v.Set(123)

	b.SetParallelism(10)
	b.ReportAllocs()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			ref, ok := v.Acquire()
			if !ok {
				b.Fatal("no value")
			}
			ref.Release()
		}
	})

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec

	b.ReportMetric(ops/1000_000, "mops")
}
