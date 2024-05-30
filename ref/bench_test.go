package ref

import (
	"sync"
	"testing"
)

// Sharded

func BenchmarkSharded_Get(b *testing.B) {
	v := new(int)
	*v = 123

	r := NewNoop(v)
	s := newSharded(r)

	for i := 0; i < b.N; i++ {
		r1, ok := s.Get()
		if !ok {
			b.Fail()
		}
		r1.Release()
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkSharded_Get_Parallel(b *testing.B) {
	v := new(int)
	*v = 123

	r := NewNoop(v)
	s := newSharded(r)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r1, ok := s.Get()
			if !ok {
				b.Fail()
			}
			r1.Release()
		}
	})

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

// NonSharded

func BenchmarkNonSharded_Get(b *testing.B) {
	v := new(int)
	*v = 123

	r := NewNoop(v)
	non := newNonSharded(r)

	for i := 0; i < b.N; i++ {
		r1, ok := non.Get()
		if !ok {
			b.Fail()
		}
		r1.Release()
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

func BenchmarkNonSharded_Get_Parallel(b *testing.B) {
	v := new(int)
	*v = 123

	r := NewNoop(v)
	non := newNonSharded(r)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r1, ok := non.Get()
			if !ok {
				b.Fail()
			}
			r1.Release()
		}
	})

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops/1000_000, "mops")
}

// private

type nonsharded[T any] struct {
	mu  sync.Mutex
	ref R[T]
}

func newNonSharded[T any](ref R[T]) *nonsharded[T] {
	s := &nonsharded[T]{}
	if ref != nil {
		s.ref = Retain(ref)
	}
	return s
}

func (s *nonsharded[T]) Get() (R[T], bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ref == nil {
		return nil, false
	}

	s.ref.Retain()
	return s.ref, true
}
