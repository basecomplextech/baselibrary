// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncrc

import (
	"sync"
	"sync/atomic"

	"github.com/basecomplextech/baselibrary/async"
	"github.com/basecomplextech/baselibrary/collect/chans"
	"github.com/basecomplextech/baselibrary/pools"
	"github.com/basecomplextech/baselibrary/ref"
	"github.com/basecomplextech/baselibrary/status"
)

// Promise is a reference counted promise.
type Promise[T any] interface {
	async.Promise[T]

	// Refcount returns the current reference count.
	Refcount() int64

	// Retain increments the reference count, panics if the promise is already released.
	Retain()

	// Release decrements the reference count, frees the future if it reaches zero.
	Release()
}

// NewPromise returns a pending reference counted promise.
func NewPromise[T any]() Promise[T] {
	return newPromise[T]()
}

// internal

var _ Promise[any] = (*promise[any])(nil)

type promise[T any] struct {
	refs  ref.Atomic64
	state atomic.Pointer[promiseState[T]]
}

type promiseState[T any] struct {
	pool pools.Pool[*promiseState[T]]

	mu   sync.Mutex
	wait chan struct{}

	st     status.Status
	done   bool
	result T
}

func newPromise[T any]() *promise[T] {
	s := acquirePromiseState[T]()

	p := &promise[T]{}
	p.refs.Init(1)
	p.state.Store(s)
	return p
}

// Done returns true if the future is complete.
func (p *promise[T]) Done() bool {
	s, ok := p.acquire()
	if !ok {
		return true
	}
	defer p.release()

	s.mu.Lock()
	defer s.mu.Unlock()

	return s.done
}

// Wait returns a channel which is closed when the result is available.
func (p *promise[T]) Wait() <-chan struct{} {
	s, ok := p.acquire()
	if !ok {
		return chans.Closed()
	}
	defer p.release()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.done {
		return chans.Closed()
	}
	return s.wait
}

// Reject rejects the promise with a status.
func (p *promise[T]) Reject(st status.Status) bool {
	var zero T
	return p.Complete(zero, st)
}

// Resolve completes the promise with a result and ok.
func (p *promise[T]) Resolve(result T) bool {
	return p.Complete(result, status.OK)
}

// Complete completes the promise with a status and a result.
func (p *promise[T]) Complete(result T, st status.Status) bool {
	s, ok := p.acquire()
	if !ok {
		panic("promise already released")
	}
	defer p.release()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.done {
		return false
	}

	s.st = st
	s.done = true
	s.result = result

notify:
	for {
		select {
		case s.wait <- struct{}{}:
		default:
			break notify
		}
	}
	return true
}

// Result returns a value and a status.
func (p *promise[T]) Result() (T, status.Status) {
	s, ok := p.acquire()
	if !ok {
		panic("promise already released")
	}
	defer p.release()

	s.mu.Lock()
	defer s.mu.Unlock()

	return s.result, s.st
}

// Status returns a status.
func (p *promise[T]) Status() status.Status {
	s, ok := p.acquire()
	if !ok {
		panic("promise already released")
	}
	defer p.release()

	s.mu.Lock()
	defer s.mu.Unlock()

	return s.st
}

// Retain/Release

// Refcount returns the current reference count.
func (p *promise[T]) Refcount() int64 {
	return p.refs.Refcount()
}

// Retain increments the reference count, panics if the promise is already released.
func (p *promise[T]) Retain() {
	_, ok := p.acquire()
	if !ok {
		panic("retain of released promise already")
	}
}

// Release decrements the reference count, frees the future if it reaches zero.
func (p *promise[T]) Release() {
	p.release()
}

// private

// acquire increments refs and returns the state, or immediately releases it if released.
func (p *promise[T]) acquire() (*promiseState[T], bool) {
	acquired := p.refs.Acquire()
	if acquired {
		s := p.state.Load()
		return s, true
	}

	// Release immediately
	p.release()
	return nil, false
}

// release decrements refs and returns the state to the pool if refs reach zero.
func (p *promise[T]) release() {
	released := p.refs.Release()
	if !released {
		return
	}

	// Release state
	s := p.state.Swap(nil)
	releasePromiseState(s)
}

// pool

var promiseStatePools = pools.NewPools()

func acquirePromiseState[T any]() *promiseState[T] {
	s, ok, pool := pools.Acquire1[*promiseState[T]](promiseStatePools)
	if ok {
		return s
	}

	s = &promiseState[T]{
		pool: pool,
		wait: make(chan struct{}, 1),
	}
	return s
}

func releasePromiseState[T any](s *promiseState[T]) {
	pool := s.pool
	s.reset()

	pool.Put(s)
}

func (s *promiseState[T]) reset() {
	pool := s.pool
	wait := s.wait

notify:
	for {
		select {
		case wait <- struct{}{}:
		default:
			break notify
		}
	}

drain:
	for {
		select {
		case <-wait:
		default:
			break drain
		}
	}

	*s = promiseState[T]{}
	s.pool = pool
	s.wait = wait
}
