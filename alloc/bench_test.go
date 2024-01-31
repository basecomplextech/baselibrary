package alloc

import (
	"math"
	"testing"
	"time"
	"unsafe"

	"github.com/basecomplextech/baselibrary/alloc/internal/arena"
)

func Benchmark_AllocInt64(b *testing.B) {
	a := arena.Test()
	size := unsafe.Sizeof(int64(0))

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(int64(size))

	t0 := time.Now()

	var v *int64
	for i := 0; i < b.N; i++ {
		v = Alloc[int64](a)
	}

	*v = math.MaxInt64
	if *v != math.MaxInt64 {
		b.Fatal()
	}

	sec := time.Since(t0).Seconds()
	ops := float64(b.N) / float64(sec)
	capacity := a.Cap() / (1024 * 1024)

	b.ReportMetric(ops/1000_000, "mops")
	b.ReportMetric(float64(size), "size")
	b.ReportMetric(float64(capacity), "cap,mb")
}

func Benchmark_AllocStruct(b *testing.B) {
	type Struct struct {
		Int8  int8
		Int16 int16
		Int32 int32
		Int64 int64
	}

	a := arena.Test()
	size := unsafe.Sizeof(Struct{})

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(int64(size))

	t0 := time.Now()

	var s *Struct
	for i := 0; i < b.N; i++ {
		s = Alloc[Struct](a)
		s.Int64 = math.MaxInt64
		if s.Int64 != math.MaxInt64 {
			b.Fatal()
		}
	}

	sec := time.Since(t0).Seconds()
	ops := float64(b.N) / float64(sec)
	capacity := a.Cap() / (1024 * 1024)

	b.ReportMetric(ops/1000_000, "mops")
	b.ReportMetric(float64(size), "size")
	b.ReportMetric(float64(capacity), "cap,mb")
}

func Benchmark_AllocBytes(b *testing.B) {
	a := arena.Test()
	size := 16

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(int64(size))

	t0 := time.Now()

	var v []byte
	for i := 0; i < b.N; i++ {
		v = Bytes(a, size)
		if len(v) != size {
			b.Fatal()
		}
	}

	sec := time.Since(t0).Seconds()
	ops := float64(b.N) / float64(sec)
	capacity := a.Cap() / (1024 * 1024)

	b.ReportMetric(ops/1000_000, "mops")
	b.ReportMetric(float64(size), "size")
	b.ReportMetric(float64(capacity), "cap,mb")
}

func Benchmark_AllocSlice(b *testing.B) {
	a := arena.Test()
	n := 4
	size := n * 4

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(int64(size))

	t0 := time.Now()

	var v []int32
	for i := 0; i < b.N; i++ {
		v = Slice[[]int32](a, n, n)
		if len(v) != 4 {
			b.Fatal()
		}
	}

	sec := time.Since(t0).Seconds()
	ops := float64(b.N) / float64(sec)
	capacity := a.Cap() / (1024 * 1024)

	b.ReportMetric(ops/1000_000, "mops")
	b.ReportMetric(float64(size), "size")
	b.ReportMetric(float64(capacity), "cap,mb")
}

func Benchmark_Alloc(b *testing.B) {
	a := arena.Test()
	size := 8

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(int64(size))

	t0 := time.Now()

	var v unsafe.Pointer
	for i := 0; i < b.N; i++ {
		v = a.Alloc(size)
		if uintptr(v) == 0 {
			b.Fatal()
		}
	}

	sec := time.Since(t0).Seconds()
	ops := float64(b.N) / float64(sec)
	capacity := a.Cap() / (1024 * 1024)

	b.ReportMetric(ops/1000_000, "mops")
	b.ReportMetric(float64(size), "size")
	b.ReportMetric(float64(capacity), "cap,mb")
}

func Benchmark_FreeList_Get_Put(b *testing.B) {
	a := arena.Test()
	list := NewFreeList[int64](a)
	size := 8

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(int64(size))

	t0 := time.Now()

	for i := 0; i < b.N; i++ {
		v := list.Get()
		list.Put(v)
	}

	sec := time.Since(t0).Seconds()
	ops := float64(b.N) / float64(sec)
	capacity := a.Cap() / (1024 * 1024)

	b.ReportMetric(ops/1000_000, "mops")
	b.ReportMetric(float64(size), "size")
	b.ReportMetric(float64(capacity), "cap,mb")
}

// Heap

func BenchmarkHeap_Alloc_Free(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	t0 := time.Now()

	for i := 0; i < b.N; i++ {
		block := globalHeap.Alloc(0)
		globalHeap.Free(block)
	}

	sec := time.Since(t0).Seconds()
	ops := float64(b.N) / float64(sec)

	b.ReportMetric(ops/1000_000, "mops")
}

// StringFormat

func BenchmarkStringFormat(b *testing.B) {
	a := arena.Test()

	b.ResetTimer()
	b.ReportAllocs()

	t0 := time.Now()

	for i := 0; i < b.N; i++ {
		a.Reset()

		s := StringFormat(a, "hello %s", "world")
		if len(s) == 0 {
			b.Fatal()
		}
	}

	sec := time.Since(t0).Seconds()
	ops := float64(b.N) / float64(sec)

	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkStringFormat_Parallel(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	t0 := time.Now()

	b.RunParallel(func(p *testing.PB) {
		a := arena.Test()

		for p.Next() {
			a.Reset()

			s := StringFormat(a, "hello %s", "world")
			if len(s) == 0 {
				b.Fatal()
			}
		}
	})

	sec := time.Since(t0).Seconds()
	ops := float64(b.N) / float64(sec)

	b.ReportMetric(ops/1000_000, "mops")
}
