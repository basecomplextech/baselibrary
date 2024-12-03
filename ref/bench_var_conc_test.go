// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"testing"
)

func BenchmarkConcurrentVar(b *testing.B) {
	r := NewNoop(1)
	v := NewConcurrentVar[int]()
	v.SetRetain(r)
	r.Release()

	for i := 0; i < b.N; i++ {
		r, ok := v.Acquire()
		if !ok {
			b.Fatal("acquire failed")
		}
		r.Release()
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkConcurrentVar_Parallel(b *testing.B) {
	r := NewNoop(1)
	v := NewConcurrentVar[int]()
	v.SetRetain(r)
	r.Release()

	b.SetParallelism(10)
	b.ReportAllocs()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			r, ok := v.Acquire()
			if !ok {
				b.Fatal("acquire failed")
			}
			r.Release()
		}
	})

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

// Acquire

func BenchmarkConcurrentVar_Acquire(b *testing.B) {
	r := NewNoop(1)
	v := NewConcurrentVar[int]()
	v.SetRetain(r)
	r.Release()

	for i := 0; i < b.N; i++ {
		r, ok := v.Acquire()
		if !ok {
			b.Fatal("acquire failed")
		}
		_ = r
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkConcurrentVar_Acquire_Parallel(b *testing.B) {
	r := NewNoop(1)
	v := NewConcurrentVar[int]()
	v.SetRetain(r)
	r.Release()

	b.SetParallelism(10)
	b.ReportAllocs()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			r, ok := v.Acquire()
			if !ok {
				b.Fatal("acquire failed")
			}
			_ = r
		}
	})

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

// SetRetain

func BenchmarkConcurrentVar_SetRetain(b *testing.B) {
	v := NewConcurrentVar[int]()

	for i := 0; i < b.N; i++ {
		r := newDummyRef[int]()
		v.SetRetain(r)
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}
