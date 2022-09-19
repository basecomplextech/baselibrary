package async

import (
	"sync"

	"github.com/epochtimeout/baselibrary/status"
)

// Result is an abstract async result.
type Result[T any] interface {
	// Wait returns a channel which is closed when the result is available.
	Wait() <-chan struct{}

	// Status returns a status.
	Status() status.Status

	// Result returns a value and a status.
	Result() (T, status.Status)
}

// internal

type result[T any] struct {
	mu     sync.Mutex
	done   bool
	result T
	status status.Status
	wait   chan struct{}
}

func newResult[T any]() result[T] {
	return result[T]{
		wait: make(chan struct{}),
	}
}

// Wait returns a channel which is closed when the result is available.
func (r *result[T]) Wait() <-chan struct{} {
	return r.wait
}

// Status returns a status.
func (r *result[T]) Status() status.Status {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.status
}

// Result returns a value and a status.
func (r *result[T]) Result() (T, status.Status) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.result, r.status
}

func (r *result[T]) complete(result T, status status.Status) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.done {
		return false
	}

	r.result = result
	r.status = status
	r.done = true

	close(r.wait)
	return true
}
