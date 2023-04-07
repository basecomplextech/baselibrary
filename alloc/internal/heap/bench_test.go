package heap

import (
	"testing"
	"time"
)

func Benchmark_Alloc_Free(b *testing.B) {
	h := New()

	b.ResetTimer()
	b.ReportAllocs()

	t0 := time.Now()
	size := 1024

	for i := 0; i < b.N; i++ {
		block := h.Alloc(size)
		h.Free(block)
	}

	sec := time.Since(t0).Seconds()
	ops := float64(b.N) / float64(sec)
	b.ReportMetric(ops, "ops")
}

func Benchmark_blockPool(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	t0 := time.Now()
	maxSize := 1 << maxIndex

	for i := 0; i < b.N; i++ {
		i := blockPool(maxSize)
		if i != maxIndex {
			b.Fatal()
		}
	}

	sec := time.Since(t0).Seconds()
	ops := float64(b.N) / float64(sec)
	b.ReportMetric(ops, "ops")
}
