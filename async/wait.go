package async

import (
	"github.com/basecomplextech/baselibrary/status"
)

// Waiter awaits an operation completion.
type Waiter interface {
	// Wait returns a channel which is closed when the operation is complete.
	Wait() <-chan struct{}
}

// WaitAll awaits all operations completion in a group.
func WaitAll[W Waiter](cancel <-chan struct{}, group ...W) status.Status {
	for _, f := range group {
		select {
		case <-f.Wait():
		case <-cancel:
			return status.Cancelled
		}
	}
	return status.OK
}
