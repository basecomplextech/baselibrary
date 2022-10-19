package async

import (
	"github.com/complex1tech/baselibrary/status"
)

var _ Waiter = (Future[any])(nil)

// Future represents a result available in the future.
type Future[T any] interface {
	// Wait returns a channel which is closed when the future is complete.
	Wait() <-chan struct{}

	// Result returns a value and a status.
	Result() (T, status.Status)

	// Status returns a status or none.
	Status() status.Status
}

// Resolved returns a successful future.
func Resolved[T any](result T) Future[T] {
	p := newPromise[T]()
	p.Complete(result, status.OK)
	return p
}

// Rejected returns a rejected future.
func Rejected[T any](st status.Status) Future[T] {
	var zero T
	p := newPromise[T]()
	p.Complete(zero, st)
	return p
}

// Completed returns a completed future.
func Completed[T any](result T, st status.Status) Future[T] {
	p := newPromise[T]()
	p.Complete(result, st)
	return p
}
