package async

import (
	"reflect"

	"github.com/basecomplextech/baselibrary/status"
)

// Routines combines multiple routines into a slice.
type Routines[T any] []Routine[T]

// Cancel cancels all routines.
func (rr Routines[T]) Cancel() {
	for _, r := range rr {
		r.Cancel()
	}
}

// CancelWait cancels all routines and waits for them to complete.
func (rr Routines[T]) CancelWait() {
	rr.Cancel()
	rr.Wait()
}

// Wait waits for all routines to complete.
func (rr Routines[T]) Wait() {
	for _, r := range rr {
		<-r.Wait()
	}
}

// WaitAny waits for any routine to complete and returns it and its index or -1 when no more routines.
// WaitAny returns OK when any routine completes, or Cancelled when the cancel channel is closed.
func (rr Routines[T]) WaitAny(cancel <-chan struct{}) (int, Routine[T], status.Status) {
	if len(rr) == 0 {
		return -1, nil, status.OK
	}

	// Make cancel case
	cases := make([]reflect.SelectCase, 0, len(rr)+1)
	stop_ := reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(cancel),
	}
	cases = append(cases, stop_)

	// Make wait cases
	for _, r := range rr {
		wait := r.Wait()
		case_ := reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(wait),
		}
		cases = append(cases, case_)
	}

	// Select case
	i, _, _ := reflect.Select(cases)
	if i == 0 {
		return -1, nil, status.Cancelled
	}

	// Return routine
	index := i - 1
	r := rr[index]
	return index, r, status.OK
}
