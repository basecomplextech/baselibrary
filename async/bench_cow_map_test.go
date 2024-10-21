// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"slices"
	"testing"

	"github.com/basecomplextech/baselibrary/collect/slices2"
)

func BenchmarkCowMap_Read(b *testing.B) {
	m := NewCopyOnWriteMap[int, int]()
	items := benchMapItems(1024)
	for _, item := range items {
		m.Set(item, item)
	}
	b.ResetTimer()

	var j int
	for i := 0; i < b.N; i++ {
		item := items[j]

		_, ok := m.Get(item)
		if !ok {
			b.Fatal("item not found")
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkCowMap_Read_Parallel(b *testing.B) {
	m := NewCopyOnWriteMap[int, int]()
	items := benchMapItems(1024)
	for _, item := range items {
		m.Set(item, item)
	}
	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		items1 := slices.Clone(items)
		slices2.Shuffle(items1)

		var j int
		for p.Next() {
			item := items1[j]

			_, ok := m.Get(item)
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
