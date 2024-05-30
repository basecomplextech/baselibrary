package ref

import "sync"

type Sharded[T any] interface {
	// Get returns a reference, increments the reference count.
	Get() (R[T], bool)

	// Set sets a reference.
	Set(next R[T])

	// Clear clears a reference.
	Clear()
}

// NewSharded returns a new sharded reference.
func NewSharded[T any](ref R[T]) Sharded[T] {
	return newSharded[T](ref)
}

// EmptySharded returns an empty sharded reference.
func EmptySharded[T any]() Sharded[T] {
	return newSharded[T](nil)
}

// internal

var _ Sharded[any] = &sharded[any]{}

type sharded[T any] struct {
	shards [8]shard[T]
	wmu    sync.Mutex
}

func newSharded[T any](ref R[T]) *sharded[T] {
	s := &sharded[T]{}
	if ref != nil {
		s.Set(ref)
	}
	return s
}

type shard[T any] struct {
	mu  sync.Mutex // 8 bytes
	ref R[T]       // 16 bytes

	_ [232]byte // cache line padding upto 256 bytes
}

// Get returns a reference.
func (s *sharded[T]) Get() (R[T], bool) {
	n := uint32(len(s.shards))
	i := runtimeFastRandn(n)

	sh := &s.shards[i]
	return sh.get()
}

// Set sets a reference.
func (s *sharded[T]) Set(next R[T]) {
	obj := next.Unwrap()

	s.wmu.Lock()
	defer s.wmu.Unlock()

	for i := range s.shards {
		ref := NextRetain(obj, next)

		sh := &s.shards[i]
		sh.set(ref)
	}
}

// Clear clears a reference.
func (s *sharded[T]) Clear() {
	s.wmu.Lock()
	defer s.wmu.Unlock()

	for i := range s.shards {
		sh := &s.shards[i]
		sh.clear()
	}
}

// private

func (s *shard[T]) get() (R[T], bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ref := s.ref
	if ref == nil {
		return nil, false
	}

	ref.Retain()
	return ref, true
}

func (s *shard[T]) set(next R[T]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ref == nil {
		s.ref = next
	} else {
		s.ref = SwapNoRetain(s.ref, next)
	}
}

func (s *shard[T]) clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	ref := s.ref
	s.ref = nil

	if ref != nil {
		ref.Release()
	}
}
