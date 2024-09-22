// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the Business Source License (BSL 1.1)
// that can be found in the LICENSE file.

package arena

import "testing"

func BenchmarkFreeList_GetPut(b *testing.B) {
	arena := Test()
	list := NewFreeList[int64](arena)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v := list.Get()
		list.Put(v)
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}
