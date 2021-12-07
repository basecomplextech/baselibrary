package run

import (
	"runtime/debug"
	"sync"

	"github.com/baseone-run/library/async"
)

// Thread runs a function in a separate goroutine and returns its result as a future.
type Thread interface {
	async.Future

	// Stop stops a thread and returns its future.
	Stop() async.Future
}

type thread struct {
	async.Promise

	mu     sync.Mutex
	stop   bool
	stopCh chan struct{}
}

func newThread() *thread {
	return &thread{
		Promise: async.NewPromise(),
		stopCh:  make(chan struct{}),
	}
}

// Stop stops a thread and returns a result done channel.
func (th *thread) Stop() async.Future {
	th.mu.Lock()
	defer th.mu.Unlock()

	if th.stop {
		return th.Promise
	}

	th.stop = true
	close(th.stopCh)
	return th.Promise
}

func (th *thread) catch(e interface{}) {
	if e == nil {
		th.complete(nil, nil)
		return
	}

	s := debug.Stack()
	pnc := &Panic{E: e, Stack: s}
	th.complete(nil, pnc)
}

func (th *thread) complete(result interface{}, err error) {
	th.Promise.Complete(result, err)
}
