// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package worker

import (
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/basecomplextech/baselibrary/async"
	"github.com/basecomplextech/baselibrary/async/internal/queue"
	"github.com/basecomplextech/baselibrary/async/internal/service"
	"github.com/basecomplextech/baselibrary/status"
)

type Worker[T any] interface {
	service.Service

	// Clear clears the task queue.
	Clear()

	// Push enqueues a task, even if the worker is stopped.
	Push(task T)
}

// Func specifies a worker function.
type Func[T any] func(ctx async.Context, task T) status.Status

// New returns a new worker with the default number of max routines (NumCPU * 2).
func New[T any](fn Func[T]) Worker[T] {
	return newWorker(fn)
}

// NewMax returns a new worker with the specified number of max routines, 0 is unlimited.
func NewMax[T any](fn Func[T], max int) Worker[T] {
	return newWorkerMax(fn, max)
}

// internal

var _ Worker[any] = (*worker[any])(nil)

type worker[T any] struct {
	fn    Func[T]
	max   int
	queue queue.Queue[T]
	service.Service

	// routines can be accessed only when handling is true
	mu       sync.Mutex
	handling atomic.Bool
	routines routines
}

func newWorker[T any](fn Func[T]) *worker[T] {
	max := runtime.NumCPU() * 2
	return newWorkerMax(fn, max)
}

func newWorkerMax[T any](fn Func[T], max int) *worker[T] {
	w := &worker[T]{
		fn:    fn,
		max:   max,
		queue: queue.New[T](),

		routines: newRoutines(),
	}
	w.Service = service.New(w.run)
	return w
}

// Clear clears the task queue.
func (w *worker[T]) Clear() {
	w.queue.Clear()
}

// Push enqueues a task, even if the worker is stopped.
func (w *worker[T]) Push(task T) {
	w.queue.Push(task)

	// Maybe start routine
	w.startRoutine()
}

// run

func (w *worker[T]) run(ctx async.Context) status.Status {
	w.start()
	defer w.stop()

	for {
		// Start routines
		for w.startRoutine() {
		}

		// Await any stopped
		select {
		case <-ctx.Wait():
			return ctx.Status()
		case <-w.routines.wait():
		}

		// Poll stopped routine
		r, ok := w.routines.poll()
		if !ok {
			continue
		}

		// Return if routine failed
		if st := r.Status(); !st.OK() {
			return st
		}
	}
}

func (w *worker[T]) start() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.handling.Store(true)
	w.routines.start()
}

func (w *worker[T]) stop() {
	// Disable handling
	w.mu.Lock()
	w.handling.Store(false)
	w.mu.Unlock()

	// From now on, routines cannot be accessed,
	// and started concurrently by other threads.

	// Stop and await all routines
	w.routines.stop()
}

// routine

func (w *worker[T]) startRoutine() bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Check handling
	if !w.handling.Load() {
		return false
	}

	// Check need more
	num := w.routines.len()
	need := num == 0 || w.queue.Len() > 1
	if !need {
		return false
	}

	// Check max reached
	maxed := (w.max > 0) && (num >= w.max)
	if maxed {
		return false
	}

	// Start routine
	r := async.Go(w.routine)
	w.routines.add(r)
	return true
}

func (w *worker[T]) routine(ctx async.Context) status.Status {
	for {
		// Poll task
		task, ok := w.queue.Poll()
		if !ok {
			return status.OK
		}

		// Process task
		st := w.fn(ctx, task)
		if st.OK() {
			continue
		}

		// Requeue task if cancelled
		if st.Cancelled() {
			w.queue.Push(task)
		}
		return st
	}
}
