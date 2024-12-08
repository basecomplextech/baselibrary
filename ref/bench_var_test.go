// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"sync/atomic"
	"testing"
)

func BenchmarkVar(b *testing.B) {
	v := newVar[int]()
	v.Set(123)

	b.SetParallelism(10)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ref, ok := v.Acquire()
		if !ok {
			b.Fatal("no value")
		}
		ref.Release()
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec

	b.ReportMetric(ops/1000_000, "mops")
}

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

// AcquireSet

func BenchmarkVar_AcquireSet_Parallel(b *testing.B) {
	r := NewNoop(1)
	v := newVar[int]()
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

// SetRetain

func BenchmarkVar_SetRetain(b *testing.B) {
	v := NewVar[int]()

	for i := 0; i < b.N; i++ {
		r := newDummyRef[int]()
		v.SetRetain(r)
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}
