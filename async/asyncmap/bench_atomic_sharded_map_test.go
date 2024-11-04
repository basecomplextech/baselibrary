// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

import (
	"math/rand/v2"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

// Read

func BenchmarkAtomicShardedMap_Read(b *testing.B) {
	m := newAtomicShardedMap[int, int](benchMapNum)
	for i := 0; i < benchMapNum; i++ {
		m.Set(i, i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := rand.IntN(benchMapNum)

		val, ok := m.Get(key)
		if !ok {
			b.Fatal("item not found")
		}
		if val != key {
			b.Fatal("invalid value", key, val)
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkAtomicShardedMap_Read_Parallel(b *testing.B) {
	m := newAtomicShardedMap[int, int](benchMapNum)
	for i := 0; i < benchMapNum; i++ {
		m.Set(i, i)
	}
	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			key := rand.IntN(benchMapNum)

			val, ok := m.Get(key)
			if !ok {
				b.Fatal("item not found")
			}
			if val != key {
				b.Fatal("invalid value", key, val)
			}
		}
	})

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

// Write

func BenchmarkAtomicShardedMap_Write(b *testing.B) {
	m := newAtomicShardedMap[int, int](benchMapNum)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := rand.IntN(benchMapNum)
		m.Set(key, key)
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkAtomicShardedMap_Write_Parallel(b *testing.B) {
	m := newAtomicShardedMap[int, int](benchMapNum)
	b.ResetTimer()
	b.SetParallelism(10)

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			key := rand.IntN(benchMapNum)
			m.Set(key, key)
		}
	})

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

// Read/Write

func BenchmarkAtomicShardedMap_Read_Write_Parallel(b *testing.B) {
	m := newAtomicShardedMap[int, int](benchMapNum)
	b.ResetTimer()
	b.SetParallelism(10)

	done := make(chan struct{})
	stop := &atomic.Bool{}
	reads := 0

	go func() {
		defer close(done)

		for !stop.Load() {
			key := rand.IntN(benchMapNum)
			_, _ = m.Get(key)
			reads++
		}
	}()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			key := rand.IntN(benchMapNum)
			m.Set(key, key)
		}
	})

	stop.Store(true)
	<-done

	sec := b.Elapsed().Seconds()
	rops := float64(reads) / sec
	wops := float64(b.N) / sec

	b.ReportMetric(rops/1000_000, "rmops")
	b.ReportMetric(wops/1000_000, "wmops")
}

func BenchmarkAtomicShardedMap_Read_Parallel_Write_Parallel(b *testing.B) {
	m := newAtomicShardedMap[int, int](benchMapNum)
	b.ResetTimer()
	b.SetParallelism(10)

	cpus := runtime.NumCPU()
	stop := &atomic.Bool{}
	reads := int64(0)

	wg := &sync.WaitGroup{}
	wg.Add(cpus)

	for i := 0; i < cpus; i++ {
		go func() {
			defer wg.Done()

			q := 0
			for !stop.Load() {
				key := rand.IntN(benchMapNum)
				_, _ = m.Get(key)
				q++
			}

			atomic.AddInt64(&reads, int64(q))
		}()
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			key := rand.IntN(benchMapNum)

			m.Set(key, key)
		}
	})

	stop.Store(true)
	wg.Wait()

	sec := b.Elapsed().Seconds()
	rops := float64(reads) / sec
	wops := float64(b.N) / sec

	b.ReportMetric(rops/1000_000, "rmops")
	b.ReportMetric(wops/1000_000, "wmops")
}
