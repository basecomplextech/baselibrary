package async

import (
	"reflect"

	"github.com/epochtimeout/baselibrary/status"
)

// Future represents an optionally cancellable concurrent computation with a future result.
type Future[T any] interface {
	Result[T]

	// Cancel tries to cancel the future and returns the wait channel.
	Cancel() <-chan struct{}
}

// Resolved returns a resolved future.
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

// Await

// AwaitAny waits for any future result and returns its index or -1 when no more results.
// AwaitAny returns OK when any future completes, or Cancelled when the stop channel is closed.
func AwaitAny[T any](stop <-chan struct{}, futures ...Future[T]) (index int, st status.Status) {
	if len(futures) == 0 {
		return -1, status.OK
	}

	// make stop case
	cases := make([]reflect.SelectCase, 0, len(futures)+1)
	stop_ := reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(stop),
	}
	cases = append(cases, stop_)

	// make future cases
	for _, f := range futures {
		wait := f.Wait()
		case_ := reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(wait),
		}
		cases = append(cases, case_)
	}

	// select case
	i, _, _ := reflect.Select(cases)
	if i == 0 {
		return 0, status.Cancelled
	}

	// return index
	index = i - 1
	return index, status.OK
}

// AwaitAll waits for all future results.
// AwaitAll returns OK when all future complete, or Cancelled when the stop channel is closed.
func AwaitAll[T any](stop <-chan struct{}, futures ...Future[T]) status.Status {
	if len(futures) == 0 {
		return status.OK
	}

	for _, f := range futures {
		select {
		case <-stop:
			return status.Cancelled
		case <-f.Wait():
		}
	}

	return status.OK
}

// Cancel

// CancelAll cancels all futures.
func CancelAll[T any](futures ...Future[T]) {
	for _, f := range futures {
		f.Cancel()
	}
}

// CancelWait cancels a future and awaits its result.
func CancelWait[T any](f Future[T]) (T, status.Status) {
	<-f.Cancel()
	return f.Result()
}
