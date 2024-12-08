// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"sync/atomic"
	"testing"
)

func BenchmarkShardedVar(b *testing.B) {
	r := NewNoop(1)
	v := newShardedVar[int]()
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

func BenchmarkShardedVar_Parallel(b *testing.B) {
	r := NewNoop(1)
	v := newShardedVar[int]()
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

func BenchmarkShardedVar_AcquireSet_Parallel(b *testing.B) {
	r := NewNoop(1)
	v := newShardedVar[int]()
	v.SetRetain(r)
	r.Release()

	b.ReportAllocs()

	stop := make(chan struct{})
	var writes atomic.Int64
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
			}

			r := newDummyRef[int]()
			v.SetRetain(r)

			writes.Add(1)
		}
	}()
	defer close(stop)

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
	b.ReportMetric(float64(writes.Load()), "writes")
}

// Acquire

func BenchmarkShardedVar_Acquire(b *testing.B) {
	r := NewNoop(1)
	v := newShardedVar[int]()
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

func BenchmarkShardedVar_Acquire_Parallel(b *testing.B) {
	r := NewNoop(1)
	v := newShardedVar[int]()
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

func BenchmarkShardedVar_SetRetain(b *testing.B) {
	v := newShardedVar[int]()

	for i := 0; i < b.N; i++ {
		r := newDummyRef[int]()
		v.SetRetain(r)
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}
