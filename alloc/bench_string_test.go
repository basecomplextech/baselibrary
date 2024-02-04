package alloc

import (
	"testing"
	"time"

	"github.com/basecomplextech/baselibrary/alloc/internal/arena"
)

// String

func BenchmarkString(b *testing.B) {
	a := arena.Test()
	s := "hello, world"
	max := a.Cap() - int64(len(s)) - 8

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if a.Len() >= max {
			a.Reset()
		}

		s1 := String(a, s)
		if len(s1) == 0 {
			b.Fatal()
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / float64(sec)

	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkStringBytes(b *testing.B) {
	a := arena.Test()
	s := []byte("hello, world")
	max := a.Cap() - int64(len(s)) - 8

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if a.Len() >= max {
			a.Reset()
		}

		s1 := StringBytes(a, s)
		if len(s1) == 0 {
			b.Fatal()
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / float64(sec)

	b.ReportMetric(ops/1000_000, "mops")
}

// StringFormat

func BenchmarkStringFormat(b *testing.B) {
	a := arena.Test()

	b.ResetTimer()
	b.ReportAllocs()

	t0 := time.Now()

	for i := 0; i < b.N; i++ {
		a.Reset()

		s := StringFormat(a, "hello %s", "world")
		if len(s) == 0 {
			b.Fatal()
		}
	}

	sec := time.Since(t0).Seconds()
	ops := float64(b.N) / float64(sec)

	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkStringFormat_Parallel(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	t0 := time.Now()

	b.RunParallel(func(p *testing.PB) {
		a := arena.Test()

		for p.Next() {
			a.Reset()

			s := StringFormat(a, "hello %s", "world")
			if len(s) == 0 {
				b.Fatal()
			}
		}
	})

	sec := time.Since(t0).Seconds()
	ops := float64(b.N) / float64(sec)

	b.ReportMetric(ops/1000_000, "mops")
}
