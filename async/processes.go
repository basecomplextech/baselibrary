package async

import (
	"reflect"

	"github.com/complex1tech/baselibrary/status"
)

// Processes combines multiple processes into a slice.
type Processes[T any] []Process[T]

// Cancel cancels all processes.
func (pp Processes[T]) Cancel() {
	for _, p := range pp {
		p.Cancel()
	}
}

// CancelWait cancels all processes and waits for them to complete.
func (pp Processes[T]) CancelWait() {
	pp.Cancel()
	pp.Wait()
}

// Wait waits for all processes to complete.
func (pp Processes[T]) Wait() {
	for _, p := range pp {
		<-p.Wait()
	}
}

// WaitAny waits for any routine to complete and returns it and its index or -1 when no more processes.
// WaitAny returns OK when any routine completes, or Cancelled when the cancel channel is closed.
func (pp Processes[T]) WaitAny(cancel <-chan struct{}) (int, Process[T], status.Status) {
	if len(pp) == 0 {
		return -1, nil, status.OK
	}

	// make cancel case
	cases := make([]reflect.SelectCase, 0, len(pp)+1)
	stop_ := reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(cancel),
	}
	cases = append(cases, stop_)

	// make wait cases
	for _, p := range pp {
		wait := p.Wait()
		case_ := reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(wait),
		}
		cases = append(cases, case_)
	}

	// select case
	i, _, _ := reflect.Select(cases)
	if i == 0 {
		return -1, nil, status.Cancelled
	}

	// return routine
	index := i - 1
	p := pp[index]
	return index, p, status.OK
}
