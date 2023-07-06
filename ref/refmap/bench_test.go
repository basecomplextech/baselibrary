package refmap

import (
	"testing"
	"time"
)

const benchTableSize = 100_000

func BenchmarkTree_Put(b *testing.B) {
	btree := testBtree(b)
	items := testItemsN(benchTableSize)

	b.ResetTimer()
	b.ReportAllocs()

	t0 := time.Now()

	j := 0
	for i := 0; i < b.N; i++ {
		if j > len(items)-1 {
			j = 0
		}

		item := items[j]
		btree.Put(item.Key, item.Value)

		j++
	}

	sec := time.Since(t0).Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkTree_Clone_Put(b *testing.B) {
	btree := testBtree(b)
	btree.Freeze()
	items := testItemsN(benchTableSize)

	b.ResetTimer()
	b.ReportAllocs()

	t0 := time.Now()

	j := 0
	for i := 0; i < b.N; i++ {
		if j > len(items)-1 {
			j = 0
		}

		item := items[j]
		clone := btree.Clone()
		clone.Put(item.Key, item.Value)
		clone.Freeze()

		btree = testUnwrap(clone)
		j++
	}

	sec := time.Since(t0).Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}
