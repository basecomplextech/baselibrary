package async

import (
	"reflect"

	"github.com/epochtimeout/baselibrary/status"
)

// Result is a generic future result interface.
type Result[T any] interface {
	// Wait awaits the result.
	Wait() <-chan struct{}

	// Result returns a value and a status or zero.
	Result() (T, status.Status)

	// Status returns a status.
	Status() status.Status
}

// Any awaits for any result and returns its index or -1 when no more results.
func Any[T any](stop <-chan struct{}, results ...Result[T]) (index int, st status.Status) {
	if len(results) == 0 {
		return -1, status.OK
	}

	// make stop case
	cases := make([]reflect.SelectCase, 0, len(results)+1)
	stop_ := reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(stop),
	}
	cases = append(cases, stop_)

	// make result cases
	for _, result := range results {
		wait := result.Wait()
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
