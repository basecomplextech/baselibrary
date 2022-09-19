package async

import (
	"sync"

	"github.com/epochtimeout/baselibrary/status"
)

// Thread executes a function in a goroutine, recovers on panics, and returns a future result.
// Thread can be restarted multiple times.
type Thread[T any] interface {
	Result[T]

	// Flags

	// Running returns a channel which is closed when the internal goroutine is running.
	Running() <-chan struct{}

	// Stopped returns a channel which is closed when the thread is stopped.
	Stopped() <-chan struct{}

	// Start/stop

	// Start starts/restarts the thread if not already running.
	Start()

	// Stop stops the thread and returns a channel which is closed when the thread stops.
	Stop() <-chan struct{}
}

// NewThread returns a new thread.
func NewThread[T any](fn func(stop <-chan struct{}) (T, status.Status)) Thread[T] {
	return newThread(fn)
}

// NewVoidThread returns a thread which does not a result.
func NewVoidThread(fn func(stop <-chan struct{}) status.Status) Thread[struct{}] {
	return newThread(func(stop <-chan struct{}) (struct{}, status.Status) {
		return struct{}{}, fn(stop)
	})
}

// internal

type thread[T any] struct {
	fn func(stop <-chan struct{}) (T, status.Status)

	running *Flag
	stopped *Flag

	mu      sync.Mutex
	result  result[T]
	routine Future[T]
}

func newThread[T any](fn func(stop <-chan struct{}) (T, status.Status)) *thread[T] {
	return &thread[T]{
		fn: fn,

		running: NewFlag(),
		stopped: SetFlag(),

		result: newResult[T](),
	}
}

// Wait returns a channel which is closed when the result is available.
func (th *thread[T]) Wait() <-chan struct{} {
	th.mu.Lock()
	defer th.mu.Unlock()

	return th.result.Wait()
}

// Status returns a status.
func (th *thread[T]) Status() status.Status {
	th.mu.Lock()
	defer th.mu.Unlock()

	return th.result.Status()
}

// Result returns a value and a status.
func (th *thread[T]) Result() (T, status.Status) {
	th.mu.Lock()
	defer th.mu.Unlock()

	return th.result.Result()
}

// Flags

// Running returns a channel which is closed when the internal goroutine is running.
func (th *thread[T]) Running() <-chan struct{} {
	return th.running.Wait()
}

// Stopped returns a channel which is closed when the thread is stopped.
func (th *thread[T]) Stopped() <-chan struct{} {
	return th.stopped.Wait()
}

// Start/stop

// Start starts/restarts the thread if not already running.
func (th *thread[T]) Start() {
	th.mu.Lock()
	defer th.mu.Unlock()

	if th.routine != nil {
		return
	}

	if th.result.done {
		th.result = newResult[T]()
	}

	th.routine = Execute(th.run)
	th.running.Signal()
	th.stopped.Reset()
}

// Stop stops the thread and returns a channel which is closed when the thread stops.
func (th *thread[T]) Stop() <-chan struct{} {
	th.mu.Lock()
	defer th.mu.Unlock()

	if th.routine == nil {
		var zero T
		th.result.complete(zero, status.Cancelled)
		return th.result.Wait()
	}

	return th.routine.Cancel()
}

// private

func (th *thread[T]) run(stop <-chan struct{}) (result T, st status.Status) {
	defer func() {
		th.mu.Lock()
		defer th.mu.Unlock()

		th.routine = nil
		th.result.complete(result, st)

		th.running.Reset()
		th.stopped.Signal()
	}()

	result, st = th.fn(stop)
	return
}
