package async

import (
	"slices"
	"testing"

	"github.com/basecomplextech/baselibrary/collect/slices2"
)

// Read

func BenchmarkConcurrentMap_Read(b *testing.B) {
	m := NewConcurrentMap[int, int]()
	items := benchMapItems(1024)
	for _, item := range items {
		m.Store(item, item)
	}
	b.ResetTimer()

	var j int
	for i := 0; i < b.N; i++ {
		item := items[j]

		_, ok := m.Load(item)
		if !ok {
			b.Fatal("item not found")
		}

		j++
		if j >= len(items) {
			j = 0
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkConcurrentMap_Read_Parallel(b *testing.B) {
	m := NewConcurrentMap[int, int]()
	items := benchMapItems(1024)
	for _, item := range items {
		m.Store(item, item)
	}
	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		items1 := slices.Clone(items)
		slices2.Shuffle(items1)

		var j int
		for p.Next() {
			item := items1[j]

			_, ok := m.Load(item)
			if !ok {
				b.Fatal("item not found")
			}

			j++
			if j >= len(items) {
				j = 0
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
	items := benchMapItems(1024)

	b.ResetTimer()

	var j int
	for i := 0; i < b.N; i++ {
		item := items[j]

		m.Store(item, item)
		_, _ = m.Load(item)
		m.Delete(item)

		j++
		if j >= len(items) {
			j = 0
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkConcurrentMap_Write_Parallel(b *testing.B) {
	m := NewConcurrentMap[int, int]()
	items := benchMapItems(1024)
	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		items1 := slices.Clone(items)
		slices2.Shuffle(items1)
		var j int

		for p.Next() {
			item := items[j]

			m.Store(item, item)
			_, _ = m.Load(item)
			m.Delete(item)

			j++
			if j >= len(items) {
				j = 0
			}
		}
	})

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}
