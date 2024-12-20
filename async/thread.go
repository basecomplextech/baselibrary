// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"sync"

	"github.com/basecomplextech/baselibrary/status"
)

// Thread is a routine thread which has to be started manually.
// The thread won't start if already stopped.
type Thread[T any] interface {
	Routine[T]

	// Start start the thread, but only if not started or stopped yet.
	Start()
}

// ThreadDyn is a thread interface without generics, i.e. Thread[?].
type ThreadDyn interface {
	RoutineDyn

	// Start start the thread, but only if not started or stopped yet.
	Start()
}

// NewThread returns a new thread which has to be started manually.
func NewThread[T any](fn func(ctx Context) (T, status.Status)) Thread[T] {
	return newThread(fn)
}

// NewThread returns a new thread which has to be started manually.
func NewThreadDyn(fn func(ctx Context) status.Status) ThreadDyn {
	fn1 := func(ctx Context) (struct{}, status.Status) {
		return struct{}{}, fn(ctx)
	}
	return newThread(fn1)
}

// internal

var _ Thread[any] = (*thread[any])(nil)

type thread[T any] struct {
	mu      sync.Mutex
	routine routine[T]

	started bool
	stopped bool
}

func newThread[T any](fn func(ctx Context) (T, status.Status)) *thread[T] {
	return &thread[T]{
		routine: newRoutineEmbedded(fn),
	}
}

// Done returns true if the future is complete.
func (t *thread[T]) Done() bool {
	return t.routine.result.Done()
}

// Wait returns a channel which is closed when the future is complete.
func (t *thread[T]) Wait() <-chan struct{} {
	return t.routine.result.Wait()
}

// Result returns a value and a status.
func (t *thread[T]) Result() (T, status.Status) {
	return t.routine.result.Result()
}

// Status returns a status or none.
func (t *thread[T]) Status() status.Status {
	return t.routine.result.Status()
}

// Start/Stop

// Start start the thread, but only if not started or stopped yet.
func (t *thread[T]) Start() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.stopped || t.started {
		return
	}
	t.started = true

	go t.routine.Run()
}

// Stop requests the routine to stop and returns a wait channel.
func (t *thread[T]) Stop() <-chan struct{} {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.stopped {
		return t.routine.result.Wait()
	}

	if t.started {
		t.routine.Stop()
	} else {
		t.routine.ctx.Free()
		t.routine.result.Reject(status.Cancelled)
	}

	t.stopped = true
	return t.routine.result.Wait()
}
