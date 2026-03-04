// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package refmap

import (
	"testing"
)

const benchTableSize = 100_000

func BenchmarkMap_SetRetain(b *testing.B) {
	btree := testBtree(b)
	items := testItemsN(benchTableSize)

	b.ResetTimer()
	b.ReportAllocs()

	j := 0
	for i := 0; i < b.N; i++ {
		if j > len(items)-1 {
			j = 0
		}

		item := items[j]
		btree.SetRetain(item.Key, item.Value)

		j++
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkMap_Clone(b *testing.B) {
	btree := testBtree(b)
	btree.Freeze()
	items := testItemsN(benchTableSize)

	b.ResetTimer()
	b.ReportAllocs()

	j := 0
	for i := 0; i < b.N; i++ {
		if j > len(items)-1 {
			j = 0
		}

		clone := btree.Clone()
		clone.Freeze()

		prev := btree
		btree = testUnwrap(clone)
		prev.Free()
		j++
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkMap_Clone_Set(b *testing.B) {
	btree := testBtree(b)
	btree.Freeze()
	items := testItemsN(benchTableSize)

	b.ResetTimer()
	b.ReportAllocs()

	j := 0
	for i := 0; i < b.N; i++ {
		if j > len(items)-1 {
			j = 0
		}

		item := items[j]
		clone := btree.Clone()
		clone.SetRetain(item.Key, item.Value)
		clone.Freeze()

		prev := btree
		btree = testUnwrap(clone)
		prev.Free()
		j++
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkMap_Clone_Set__small(b *testing.B) {
	btree := testBtree(b)
	btree.Freeze()
	items := testItemsN(16)

	b.ResetTimer()
	b.ReportAllocs()

	j := 0
	for i := 0; i < b.N; i++ {
		if j > len(items)-1 {
			j = 0
		}

		item := items[j]
		clone := btree.Clone()
		clone.SetRetain(item.Key, item.Value)
		clone.Freeze()
		clone.Free()

		j++
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkMap_Iterator(b *testing.B) {
	items := testItemsN(benchTableSize)
	btree := testBtree(b, items...)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		it := btree.Iterator()
		it.SeekToStart()
		it.Next()
		it.Free()
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkIterator(b *testing.B) {
	items := testItemsN(benchTableSize)
	btree := testBtree(b, items...)
	it := btree.Iterator()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _, ok := it.Next()
		if !ok {
			it.SeekToStart()
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}
