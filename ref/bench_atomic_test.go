// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkAtomic_Parallel(b *testing.B) {
	v := int64(0)

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			atomic.AddInt64(&v, 1)
		}
	})

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkMutex_Parallel(b *testing.B) {
	v := int64(0)
	mu := &sync.Mutex{}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			mu.Lock()
			v++
			mu.Unlock()
		}
	})

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}
