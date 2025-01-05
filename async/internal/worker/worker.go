// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package worker

import (
	"runtime"
	"sync"

	"github.com/basecomplextech/baselibrary/async"
	"github.com/basecomplextech/baselibrary/async/internal/queue"
	"github.com/basecomplextech/baselibrary/collect/chans"
	"github.com/basecomplextech/baselibrary/opt"
	"github.com/basecomplextech/baselibrary/status"
)

type Worker[T any] interface {
	// Status returns the exit status or none.
	Status() status.Status

	// Flags

	// Running indicates that the worker is running.
	Running() async.Flag

	// Stopped indicates that the worker is stopped.
	Stopped() async.Flag

	// Lifecycle

	// Start starts the worker if not running.
	Start()

	// Stop requests the worker to stop and returns a stopped channel.
	Stop() <-chan struct{}

	// Wait returns a channel which is closed when the worker is stopped.
	Wait() <-chan struct{}

	// Tasks

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

	// flags
	stopped async.MutFlag

	// main routine
	mainMu sync.Mutex
	main   opt.Opt[async.RoutineDyn]

	// state can be accessed only when running is true
	runMu    sync.Mutex
	running  async.MutFlag
	routines routineGroup
}

func newWorker[T any](fn Func[T]) *worker[T] {
	max := runtime.NumCPU() * 2
	return newWorkerMax(fn, max)
}

func newWorkerMax[T any](fn Func[T], max int) *worker[T] {
	return &worker[T]{
		fn:    fn,
		max:   max,
		queue: queue.New[T](),

		stopped: async.UnsetFlag(),

		running:  async.UnsetFlag(),
		routines: newRoutineGroup(),
	}
}

// Status returns the exit status or none.
func (w *worker[T]) Status() status.Status {
	w.mainMu.Lock()
	defer w.mainMu.Unlock()

	r, ok := w.main.Unwrap()
	if !ok {
		return status.Unavailable("worker is stopped")
	}
	return r.Status()
}

// Flags

// Running indicates that the worker is running.
func (w *worker[T]) Running() async.Flag {
	return w.running
}

// Stopped indicates that the worker is stopped.
func (w *worker[T]) Stopped() async.Flag {
	return w.stopped
}

// Lifecycle

// Start starts the worker if not running.
func (w *worker[T]) Start() {
	w.mainMu.Lock()
	defer w.mainMu.Unlock()

	// Check routine
	main, ok := w.main.Unwrap()
	if ok && !main.Done() {
		return
	}

	// Start main
	main = async.Go(w.run)
	w.main.Set(main)
	return
}

// Stop requests the worker to stop and returns a stopped channel.
func (w *worker[T]) Stop() <-chan struct{} {
	w.mainMu.Lock()
	defer w.mainMu.Unlock()

	main, ok := w.main.Unwrap()
	if !ok {
		return chans.Closed()
	}
	return main.Stop()
}

// Wait returns a channel which is closed when the worker is stopped.
func (w *worker[T]) Wait() <-chan struct{} {
	w.mainMu.Lock()
	defer w.mainMu.Unlock()

	main, ok := w.main.Unwrap()
	if !ok {
		return chans.Closed()
	}
	return main.Wait()
}

// Tasks

// Clear clears the task queue.
func (w *worker[T]) Clear() {
	w.runMu.Lock()
	defer w.runMu.Unlock()

	w.queue.Clear()
}

// Push enqueues a task, even if the worker is stopped.
func (w *worker[T]) Push(task T) {
	w.queue.Push(task)

	// Maybe start routine
	w.startRoutine()
}

// main

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
	w.runMu.Lock()
	defer w.runMu.Unlock()

	w.running.Set()
	w.stopped.Unset()
	w.routines.start()
}

func (w *worker[T]) stop() {
	defer w.stopped.Set()

	// Disable running
	w.runMu.Lock()
	w.running.Unset()
	w.runMu.Unlock()

	// From now on, routines cannot be accessed,
	// and started concurrently by other threads.

	// Stop and await all routines
	w.routines.stop()
}

// routine

func (w *worker[T]) startRoutine() bool {
	w.runMu.Lock()
	defer w.runMu.Unlock()

	// Check running
	if !w.running.IsSet() {
		return false
	}

	// Check need routine
	num := w.routines.len()
	need := num == 0 || w.queue.Len() > 1
	if !need {
		return false
	}

	// Check max routines
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
