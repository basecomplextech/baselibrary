// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"math/rand/v2"
	"testing"
)

// Read

func BenchmarkConcurrentMap_Read(b *testing.B) {
	m := NewConcurrentMap[int, int]()
	for i := 0; i < benchMapNum; i++ {
		m.Store(i, i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := rand.IntN(benchMapNum)

		_, ok := m.Load(key)
		if !ok {
			b.Fatal("item not found")
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkConcurrentMap_Read_Parallel(b *testing.B) {
	m := NewConcurrentMap[int, int]()
	for i := 0; i < benchMapNum; i++ {
		m.Store(i, i)
	}
	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			key := rand.IntN(benchMapNum)

			_, ok := m.Load(key)
			if !ok {
				b.Fatal("item not found")
			}
		}
	})

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

// Write

func BenchmarkConcurrentMap_Write(b *testing.B) {
	m := NewConcurrentMap[int, int]()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := rand.IntN(benchMapNum)

		m.Store(key, key)
		_, _ = m.Load(key)
		m.Delete(key)
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkConcurrentMap_Write_Parallel(b *testing.B) {
	m := NewConcurrentMap[int, int]()
	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			key := rand.IntN(benchMapNum)

			m.Store(key, key)
			_, _ = m.Load(key)
			m.Delete(key)
		}
	})

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}
