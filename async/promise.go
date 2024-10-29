// Copyright 2021 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"sync"

	"github.com/basecomplextech/baselibrary/status"
)

// Promise is a future which can be completed.
type Promise[T any] interface {
	Future[T]

	// Complete completes the promise with a status and a result.
	Complete(result T, st status.Status) bool

	// Reject rejects the promise with a status.
	Reject(st status.Status) bool
}

// NewPromise returns a pending promise.
func NewPromise[T any]() Promise[T] {
	return newPromise[T]()
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
		return closedChan
	}
	return p.wait
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

// Reject rejects the promise with a status.
func (p *promise[T]) Reject(st status.Status) bool {
	var zero T
	return p.Complete(zero, st)
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
