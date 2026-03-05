// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package refvar

import (
	"sync/atomic"
	"testing"

	"github.com/basecomplextech/baselibrary/ref"
)

func BenchmarkSharded(b *testing.B) {
	r := ref.NewNoop(1)
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

func BenchmarkSharded_Parallel(b *testing.B) {
	r := ref.NewNoop(1)
	v := newShardedVar[int]()
	v.SetRetain(r)
	r.Release()

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

func BenchmarkSharded_AcquireSet_Parallel(b *testing.B) {
	r := ref.NewNoop(1)
	v := newShardedVar[int]()
	v.SetRetain(r)
	r.Release()

	b.ReportAllocs()

	stop := make(chan struct{})
	var sets atomic.Int64
	go func() {
		r := ref.NewNoop(1)

		for {
			select {
			case <-stop:
				return
			default:
			}

			v.SetRetain(r)
			sets.Add(1)
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
	b.ReportMetric(float64(sets.Load()), "sets")
}

// Acquire

func BenchmarkSharded_Acquire(b *testing.B) {
	r := ref.NewNoop(1)
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

func BenchmarkSharded_Acquire_Parallel(b *testing.B) {
	r := ref.NewNoop(1)
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

func BenchmarkSharded_SetRetain(b *testing.B) {
	v := newShardedVar[int]()
	r := ref.NewNoop(1)

	for i := 0; i < b.N; i++ {
		v.SetRetain(r)
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}
