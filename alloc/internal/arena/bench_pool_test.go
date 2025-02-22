// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package arena

import "testing"

func BenchmarkPool(b *testing.B) {
	a := Test()
	p := NewPool[int64](a)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v, _ := p.Get()
		p.Put(v)
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}
