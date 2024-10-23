// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

import (
	"math/rand/v2"
	"testing"

	"github.com/basecomplextech/baselibrary/async"
)

func BenchmarkLockMap_Get(b *testing.B) {
	m := NewLockMap[int]()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := rand.IntN(1024)

		lock := m.Get(key)
		lock.Free()
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkLockMap_Get_Parallel(b *testing.B) {
	m := NewLockMap[int]()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := rand.IntN(1024)

			lock := m.Get(key)
			lock.Free()
		}
	})

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

// Lock

func BenchmarkLockMap_Lock(b *testing.B) {
	m := NewLockMap[int]()
	ctx := async.NoContext()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := rand.IntN(1024)

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

func BenchmarkLockMap_Lock_Parallel(b *testing.B) {
	m := NewLockMap[int]()
	ctx := async.NoContext()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := rand.IntN(1024)

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
