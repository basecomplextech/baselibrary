// Copyright 2021 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package promise

import (
	"sync"

	"github.com/basecomplextech/baselibrary/async/internal/future"
	"github.com/basecomplextech/baselibrary/collect/chans"
	"github.com/basecomplextech/baselibrary/status"
)

// Promise is a future which can be completed.
type Promise[T any] interface {
	future.Future[T]

	// Reject rejects the promise with a status.
	Reject(st status.Status) bool

	// Resolve completes the promise with a result and ok.
	Resolve(result T) bool

	// Complete completes the promise with a status and a result.
	Complete(result T, st status.Status) bool
}

// New returns a pending promise.
func New[T any]() Promise[T] {
	return newPromise[T]()
}

// Resolved

// Resolved returns a successful future.
func Resolved[T any](result T) future.Future[T] {
	p := newPromise[T]()
	p.Complete(result, status.OK)
	return p
}

// Rejected returns a rejected future.
func Rejected[T any](st status.Status) future.Future[T] {
	var zero T
	p := newPromise[T]()
	p.Complete(zero, st)
	return p
}

// Completed returns a completed future.
func Completed[T any](result T, st status.Status) future.Future[T] {
	p := newPromise[T]()
	p.Complete(result, st)
	return p
}

// internal

var _ Promise[any] = (*promise[any])(nil)

type promise[T any] struct {
	mu   sync.Mutex
	wait chan struct{}

	st     status.Status
	done   bool
	result T
}

func newPromise[T any]() *promise[T] {
	return &promise[T]{
		wait: make(chan struct{}),
	}
}

func newPromiseEmbedded[T any]() promise[T] {
	return promise[T]{
		wait: make(chan struct{}),
	}
}

// Done returns true if the future is complete.
func (p *promise[T]) Done() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.done
}

// Wait returns a channel which is closed when the result is available.
func (p *promise[T]) Wait() <-chan struct{} {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.done {
		return chans.Closed()
	}
	return p.wait
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
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.done {
		return false
	}

	p.st = st
	p.done = true
	p.result = result
	close(p.wait)
	return true
}

// Result returns a value and a status.
func (p *promise[T]) Result() (T, status.Status) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.result, p.st
}

// Status returns a status.
func (p *promise[T]) Status() status.Status {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.st
}
