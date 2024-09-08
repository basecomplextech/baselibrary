// Copyright 2024 Ivan Korobkov. All rights reserved.

package pools

import (
	"sync"
	"testing"
)

func BenchmarkGetPool(b *testing.B) {
	b.ReportAllocs()
	pp := NewPools()

	for i := 0; i < b.N; i++ {
		_ = GetPool[int64](pp)
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkGetPool_Parallel(b *testing.B) {
	b.ReportAllocs()
	pp := NewPools()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_ = GetPool[int64](pp)
		}
	})

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

// Acquire

func BenchmarkAcquire(b *testing.B) {
	b.ReportAllocs()
	pp := NewPools()

	for i := 0; i < b.N; i++ {
		v, _, pool := Acquire1[int64](pp)
		pool.Put(v)
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkAcquire_Parallel(b *testing.B) {
	b.ReportAllocs()
	pp := NewPools()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			v, _, pool := Acquire1[int64](pp)
			pool.Put(v)
		}
	})

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

// Get

func BenchmarkPool_Get(b *testing.B) {
	b.ReportAllocs()
	p := NewPool[int64]()

	for i := 0; i < b.N; i++ {
		v, _ := p.Get()
		p.Put(v)
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkPool_Get_Parallel(b *testing.B) {
	b.ReportAllocs()
	p := NewPool[int64]()

	b.RunParallel(func(b *testing.PB) {
		for b.Next() {
			v, _ := p.Get()
			p.Put(v)
		}
	})

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

// SyncPool

func BenchmarkSyncPool_Get(b *testing.B) {
	b.ReportAllocs()
	p := &sync.Pool{}

	for i := 0; i < b.N; i++ {
		v := p.Get()
		if v == nil {
			v = new(int64)
		}
		p.Put(v)
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkSyncPool_Get_Parallel(b *testing.B) {
	b.ReportAllocs()
	p := &sync.Pool{}

	b.RunParallel(func(b *testing.PB) {
		for b.Next() {
			v := p.Get()
			if v == nil {
				v = new(int64)
			}
			p.Put(v)
		}
	})

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}
