// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package compactint

import (
	"encoding/binary"
	"math"
	"testing"
)

// Decode

func BenchmarkUint32(b *testing.B) {
	buf := make([]byte, MaxLen)
	v := uint32(math.MaxUint32)
	n := PutUint32(buf, v)

	b.SetBytes(int64(n))

	for i := 0; i < b.N; i++ {
		v1, n1 := Uint32(buf)
		if v1 != v {
			b.Fatal()
		}
		if n1 != n {
			b.Fatal()
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(float64(ops)/1000_000, "mops")
}

func BenchmarkUint64(b *testing.B) {
	buf := make([]byte, MaxLen)
	v := uint64(math.MaxUint64)
	n := PutUint64(buf, v)

	b.SetBytes(int64(n))

	for i := 0; i < b.N; i++ {
		v1, n1 := Uint64(buf)
		if v1 != v {
			b.Fatal()
		}
		if n1 != n {
			b.Fatal()
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(float64(ops)/1000_000, "mops")
}

// Decode reverse

func BenchmarkReverseInt32(b *testing.B) {
	buf := make([]byte, MaxLen)
	v := int32(math.MaxInt32)
	n := PutReverseInt32(buf, v)
	buf = buf[len(buf)-n:]

	b.SetBytes(int64(n))

	for i := 0; i < b.N; i++ {
		v1, n1 := ReverseInt32(buf)
		if v1 != v {
			b.Fatal()
		}
		if n1 != n {
			b.Fatal()
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(float64(ops)/1000_000, "mops")
}

func BenchmarkReverseInt64(b *testing.B) {
	buf := make([]byte, MaxLen)
	v := int64(math.MaxInt64)
	n := PutReverseInt64(buf, v)
	buf = buf[len(buf)-n:]

	b.SetBytes(int64(n))

	for i := 0; i < b.N; i++ {
		v1, n1 := ReverseInt64(buf)
		if v1 != v {
			b.Fatal()
		}
		if n1 != n {
			b.Fatal()
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(float64(ops)/1000_000, "mops")
}

func BenchmarkReverseUint32(b *testing.B) {
	buf := make([]byte, MaxLen)
	v := uint32(math.MaxUint32)
	n := PutReverseUint32(buf, v)
	buf = buf[len(buf)-n:]

	b.SetBytes(int64(n))

	for i := 0; i < b.N; i++ {
		v1, n1 := ReverseUint32(buf)
		if v1 != v {
			b.Fatal()
		}
		if n1 != n {
			b.Fatal()
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(float64(ops)/1000_000, "mops")
}

func BenchmarkReverseUint64(b *testing.B) {
	buf := make([]byte, MaxLen)
	v := uint64(math.MaxUint64)
	n := PutReverseUint64(buf, v)
	buf = buf[len(buf)-n:]

	b.SetBytes(int64(n))

	for i := 0; i < b.N; i++ {
		v1, n1 := ReverseUint64(buf)
		if v1 != v {
			b.Fatal()
		}
		if n1 != n {
			b.Fatal()
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(float64(ops)/1000_000, "mops")
}

// Encode reverse

func BenchmarkPutReverseUint32(b *testing.B) {
	buf := make([]byte, MaxLen)
	b.SetBytes(4)

	for i := 0; i < b.N; i++ {
		n := PutReverseUint64(buf, math.MaxUint32)
		if n <= 0 {
			b.Fatal()
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(float64(ops)/1000_000, "mops")
}

func BenchmarkPutReverseUint64(b *testing.B) {
	buf := make([]byte, MaxLen)
	b.SetBytes(8)

	for i := 0; i < b.N; i++ {
		n := PutReverseUint64(buf, math.MaxUint64)
		if n <= 0 {
			b.Fatal()
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(float64(ops)/1000_000, "mops")
}

// Encode

func BenchmarkPutUint32(b *testing.B) {
	buf := make([]byte, MaxLen)
	b.SetBytes(4)

	for i := 0; i < b.N; i++ {
		n := PutUint64(buf, math.MaxUint32)
		if n <= 0 {
			b.Fatal()
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(float64(ops)/1000_000, "mops")
}

func BenchmarkPutUint64(b *testing.B) {
	buf := make([]byte, MaxLen)
	b.SetBytes(8)

	for i := 0; i < b.N; i++ {
		n := PutUint64(buf, math.MaxUint64)
		if n <= 0 {
			b.Fatal()
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(float64(ops)/1000_000, "mops")
}

// Standard varint

func BenchmarkUvarint64(b *testing.B) {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, math.MaxUint64)
	b.SetBytes(8)

	for i := 0; i < b.N; i++ {
		v, n := binary.Uvarint(buf)
		if v != math.MaxUint64 {
			b.Fatal()
		}
		if n != binary.MaxVarintLen64 {
			b.Fatal()
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(float64(ops)/1000_000, "mops")
}

func BenchmarkPutUvarint64(b *testing.B) {
	buf := make([]byte, binary.MaxVarintLen64)
	b.SetBytes(8)

	for i := 0; i < b.N; i++ {
		n := binary.PutUvarint(buf, math.MaxUint64)
		if n != binary.MaxVarintLen64 {
			b.Fatal()
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(float64(ops)/1000_000, "mops")
}

// Standard big endian

func BenchmarkBigEndianUint64(b *testing.B) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, math.MaxUint64)
	b.SetBytes(8)

	for i := 0; i < b.N; i++ {
		v := binary.BigEndian.Uint64(buf)
		if v != math.MaxUint64 {
			b.Fatal()
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(float64(ops)/1000_000, "mops")
}

func BenchmarkPutBigEndianUint64(b *testing.B) {
	buf := make([]byte, 8)
	b.SetBytes(8)

	for i := 0; i < b.N; i++ {
		binary.BigEndian.PutUint64(buf, math.MaxUint64)
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(float64(ops)/1000_000, "mops")
}
