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
}

// Decode reverse

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
}

func BenchmarkPutBigEndianUint64(b *testing.B) {
	buf := make([]byte, 8)
	b.SetBytes(8)

	for i := 0; i < b.N; i++ {
		binary.BigEndian.PutUint64(buf, math.MaxUint64)
	}
}
