// Copyright 2024 Ivan Korobkov. All rights reserved.

package bin

import (
	"crypto/rand"
	"testing"
)

func BenchmarkCryptoRandom64(b *testing.B) {
	b.ReportAllocs()
	buf := [8]byte{}

	for i := 0; i < b.N; i++ {
		_, _ = rand.Read(buf[:])
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkCryptoRandom64_Parallel(b *testing.B) {
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		buf := [8]byte{}

		for pb.Next() {
			_, _ = rand.Read(buf[:])
		}
	})

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

// randomReader

func BenchmarkReader_Random64(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = random.read64()
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkReader_Random64_Parallel(b *testing.B) {
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = random.read64()
		}
	})

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}
