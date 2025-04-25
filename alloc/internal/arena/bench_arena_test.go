// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package arena

import "testing"

func BenchmarkArena_Alloc(b *testing.B) {
	a := Test()
	num := 10_000
	size := 8

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		a.Alloc(size)

		if i%num == 0 {
			a.Reset()
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec

	b.SetBytes(int64(size))
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkMutexArena_Alloc(b *testing.B) {
	a := TestMutex()
	num := 10_000
	size := 8

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		a.Alloc(size)

		if i%num == 0 {
			a.Reset()
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec

	b.SetBytes(int64(size))
	b.ReportMetric(ops/1000_000, "mops")
}
