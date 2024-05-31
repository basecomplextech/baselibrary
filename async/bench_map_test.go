package async

import (
	"testing"

	"github.com/basecomplextech/baselibrary/collect/slices"
)

func BenchmarkMap_Write(b *testing.B) {
	m := NewMap[int, int]()
	items := benchMapItems(1024)

	for _, item := range items {
		m.Set(item, item)
	}
	b.ResetTimer()

	var j int
	for i := 0; i < b.N; i++ {
		item := items[j]

		m.Set(item, item)
		_, _ = m.Get(item)
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

func BenchmarkMap_Write_Parallel(b *testing.B) {
	m := NewMap[int, int]()
	items := benchMapItems(1024)
	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		items1 := slices.Clone(items)
		slices.Shuffle(items1)

		var j int
		for p.Next() {
			item := items1[j]

			m.Set(item, item)
			_, _ = m.Get(item)
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

// ConcurrentMap

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
		slices.Shuffle(items1)
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

// private

func benchMapItems(n int) []int {
	items := make([]int, n)
	for i := 0; i < n; i++ {
		items[i] = i
	}
	return items
}
