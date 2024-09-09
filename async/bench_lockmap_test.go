// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"sync/atomic"
	"testing"
)

func BenchmarkLockMap(b *testing.B) {
	m := NewLockMap[int]()
	key := 123
	ctx := NoContext()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		lock, st := m.Lock(ctx, key)
		if !st.OK() {
			b.Fatal(st)
		}
		lock.Free()
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkLockMap_Parallel(b *testing.B) {
	m := NewLockMap[int]()
	ctx := NoContext()
	delta := int64(0)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		key := 123 + int(atomic.AddInt64(&delta, 1))

		for pb.Next() {
			lock, st := m.Lock(ctx, key)
			if !st.OK() {
				b.Fatal(st)
			}

			lock.Free()
		}
	})

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}
