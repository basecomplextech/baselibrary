package pools

import (
	"sync"
	"testing"
)

func BenchmarkGetPool(b *testing.B) {
	b.ReportAllocs()
	pp := New()

	for i := 0; i < b.N; i++ {
		_ = GetPool[int32, int32](pp)
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkPool_Get(b *testing.B) {
	b.ReportAllocs()
	p := NewPool[int64]()

	for i := 0; i < b.N; i++ {
		_, _ = p.Get()
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkSyncPool_Get(b *testing.B) {
	b.ReportAllocs()
	p := &sync.Pool{}

	for i := 0; i < b.N; i++ {
		_ = p.Get()
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}
