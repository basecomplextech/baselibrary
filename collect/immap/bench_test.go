// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package immap

import (
	"testing"
)

const benchTableSize = 100_000

func BenchmarkMap_Set(b *testing.B) {
	m := testMap(b)
	items := testItemsN(benchTableSize)

	b.ResetTimer()
	b.ReportAllocs()

	j := 0
	for i := 0; i < b.N; i++ {
		if j > len(items)-1 {
			j = 0
		}

		item := items[j]
		m.Set(item.Key, item.Value)

		j++
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkMap_Clone(b *testing.B) {
	m := testMap(b)
	m.Freeze()
	items := testItemsN(benchTableSize)

	b.ResetTimer()
	b.ReportAllocs()

	j := 0
	for i := 0; i < b.N; i++ {
		if j > len(items)-1 {
			j = 0
		}

		clone := m.Clone()
		clone.Freeze()

		prev := m
		m = testUnwrap(clone)
		prev.Free()
		j++
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkMap_Clone_Set(b *testing.B) {
	m := testMap(b)
	m.Freeze()
	items := testItemsN(benchTableSize)

	b.ResetTimer()
	b.ReportAllocs()

	j := 0
	for i := 0; i < b.N; i++ {
		if j > len(items)-1 {
			j = 0
		}

		item := items[j]
		clone := m.Clone()
		clone.Set(item.Key, item.Value)
		clone.Freeze()

		prev := m
		m = testUnwrap(clone)
		prev.Free()
		j++
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkMap_Clone_Set__small(b *testing.B) {
	m := testMap(b)
	m.Freeze()
	items := testItemsN(16)

	b.ResetTimer()
	b.ReportAllocs()

	j := 0
	for i := 0; i < b.N; i++ {
		if j > len(items)-1 {
			j = 0
		}

		item := items[j]
		clone := m.Clone()
		clone.Set(item.Key, item.Value)
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
	m := testMap(b, items...)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		it := m.Iterator()
		it.Next()
		it.Free()
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkIterator(b *testing.B) {
	items := testItemsN(benchTableSize)
	m := testMap(b, items...)

	it := m.Iterator()
	defer it.Free()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _, ok := it.Next()
		if !ok {
			it.SeekToStart()
			continue
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}
