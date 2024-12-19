// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package routinepool

import "runtime"

// worker is a gc-collectable object which is stored in a sync pool.
// It uses a finalizer to close the queue channel and exit the goroutine.
type worker struct {
	pool    *pool
	queue   chan task
	running bool
}

type task struct {
	w  *worker
	fn Runner
}

func newWorker(p *pool) *worker {
	w := &worker{
		pool:  p,
		queue: make(chan task, 1),
	}

	runtime.SetFinalizer(w, func(w *worker) {
		close(w.queue)
	})
	return w
}

func (w *worker) run(runner Runner) {
	if !w.running {
		w.running = true

		go run(w.queue)
	}

	task := task{
		w:  w,
		fn: runner,
	}
	w.queue <- task
}

// routine

// run runs the tasks in the queue, and releases the worker when done.
//
// The goroutine exits when the queue is closed.
// The queue is closed when the worker is finalized.
func run(queue chan task) {
	for {
		t, ok := <-queue
		if !ok {
			return
		}

		runTask(t)
	}
}

func runTask(t task) {
	t.fn.Run()
	t.w.pool.release(t.w)
}
