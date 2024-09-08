// Copyright 2022 Ivan Korobkov. All rights reserved.

package alloc

import (
	"math"
	"testing"
	"unsafe"

	"github.com/basecomplextech/baselibrary/alloc/internal/arena"
	"github.com/basecomplextech/baselibrary/alloc/internal/heap"
)

func Benchmark_AllocInt64(b *testing.B) {
	a := arena.Test()
	size := unsafe.Sizeof(int64(0))
	max := a.Cap() - 16

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(int64(size))

	var v *int64
	for i := 0; i < b.N; i++ {
		if a.Len() >= max {
			a.Reset()
		}

		v = Alloc[int64](a)
	}

	*v = math.MaxInt64
	if *v != math.MaxInt64 {
		b.Fatal()
	}

	sec := b.Elapsed().Seconds()
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
	max := a.Cap() - 16

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(int64(size))

	var s *Struct
	for i := 0; i < b.N; i++ {
		if a.Len() >= max {
			a.Reset()
		}

		s = Alloc[Struct](a)
		s.Int64 = math.MaxInt64
		if s.Int64 != math.MaxInt64 {
			b.Fatal()
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / float64(sec)
	capacity := a.Cap() / (1024 * 1024)

	b.ReportMetric(ops/1000_000, "mops")
	b.ReportMetric(float64(size), "size")
	b.ReportMetric(float64(capacity), "cap,mb")
}

func Benchmark_AllocBytes(b *testing.B) {
	a := arena.Test()
	size := 16
	max := a.Cap() - 16

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(int64(size))

	var v []byte
	for i := 0; i < b.N; i++ {
		if a.Len() >= max {
			a.Reset()
		}

		v = Bytes(a, size)
		if len(v) != size {
			b.Fatal()
		}
	}

	sec := b.Elapsed().Seconds()
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
	max := a.Cap() - 16

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(int64(size))

	var v []int32
	for i := 0; i < b.N; i++ {
		if a.Len() >= max {
			a.Reset()
		}

		v = Slice[[]int32](a, n, n)
		if len(v) != 4 {
			b.Fatal()
		}
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / float64(sec)
	capacity := a.Cap() / (1024 * 1024)

	b.ReportMetric(ops/1000_000, "mops")
	b.ReportMetric(float64(size), "size")
	b.ReportMetric(float64(capacity), "cap,mb")
}

func Benchmark_Alloc(b *testing.B) {
	a := arena.Test()
	size := 8
	max := a.Cap() - 16

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(int64(size))

	var v unsafe.Pointer
	for i := 0; i < b.N; i++ {
		if a.Len() >= max {
			a.Reset()
		}

		v = a.Alloc(size)
		if uintptr(v) == 0 {
			b.Fatal()
		}
	}

	sec := b.Elapsed().Seconds()
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

	for i := 0; i < b.N; i++ {
		v := list.Get()
		list.Put(v)
	}

	sec := b.Elapsed().Seconds()
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

	for i := 0; i < b.N; i++ {
		block := heap.Global.Alloc(0)
		heap.Global.Free(block)
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / float64(sec)

	b.ReportMetric(ops/1000_000, "mops")
}
