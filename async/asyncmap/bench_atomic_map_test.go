// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

import (
	"math/rand/v2"
	"testing"
)

// Read

func BenchmarkAtomicMap_Read(b *testing.B) {
	m := NewAtomicMap[int, int]()
	for i := 0; i < benchMapNum; i++ {
		m.Set(i, i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := rand.IntN(benchMapNum)

		_, ok := m.Get(key)
		if !ok {
			b.Fatal("item not found")
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkAtomicMap_Read_Parallel(b *testing.B) {
	m := NewAtomicMap[int, int]()
	for i := 0; i < benchMapNum; i++ {
		m.Set(i, i)
	}
	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			key := rand.IntN(benchMapNum)

			_, ok := m.Get(key)
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

func BenchmarkAtomicMap_Write(b *testing.B) {
	m := NewAtomicMap[int, int]()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := rand.IntN(benchMapNum)

		m.Set(key, key)
		_, _ = m.Get(key)
		m.Delete(key)
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkAtomicMap_Write_Parallel(b *testing.B) {
	m := NewAtomicMap[int, int]()
	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			key := rand.IntN(benchMapNum)

			m.Set(key, key)
			_, _ = m.Get(key)
			m.Delete(key)
		}
	})

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}
