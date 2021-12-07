package run

import (
	"sync"

	"github.com/baseone-run/library/async"
	"github.com/baseone-run/library/errs"
)

var _ Group = (*group)(nil)

// Group groups multiple threads in one single group.
// Group implements the thread interface and can be used to compose complex thread trees.
type Group interface {
	Thread

	// Group

	// Add adds threads to the group.
	// Add stops (but does not wait) the threads if the group is completed.
	Add(th ...Thread)

	// Remove removes a thread from the group.
	// Remove does not stop the thread.
	Remove(th Thread)
}

type group struct {
	async.Promise

	mu      sync.Mutex
	stop    bool
	threads []Thread
}

// NewGroup copies the threads and returns a new thread group.
func NewGroup(th ...Thread) Group {
	g := &group{
		Promise: async.NewPromise(),
		threads: make([]Thread, len(th)),
	}
	copy(g.threads, th)
	return g
}

// Add adds threads to the group.
// Add stops (but does not wait) the threads if the group has stopped.
func (g *group) Add(th ...Thread) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.threads = append(g.threads, th...)
	if !g.stop {
		return
	}

	for _, th1 := range th {
		th1.Stop()
	}
}

// Remove removes a thread from the group.
// Remove does not stop the thread.
func (g *group) Remove(th Thread) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.stop {
		return
	}

	threads := make([]Thread, 0, len(g.threads))
	for _, th1 := range g.threads {
		if th1 == th {
			continue
		}
		threads = append(threads, th1)
	}
	g.threads = threads
}

// Stop stops the threads and starts an internal goroutine which awaits their stop.
func (g *group) Stop() async.Future {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.stop {
		return g.Promise
	}
	g.stop = true

	for _, th := range g.threads {
		th.Stop()
	}

	go g.await()
	return g.Promise
}

func (g *group) await() {
	threads := g.cloneThreads()
	errs := make([]error, 0, len(threads))

	for _, th := range threads {
		<-th.Done()

		if err := th.Err(); err != nil {
			errs = append(errs, err)
		}
	}

	g.complete(errs)
}

func (g *group) complete(ee []error) {
	err := errs.Combine(ee...) // err is nil when errs are empty.
	g.Promise.Complete(nil, err)
}

func (g *group) cloneThreads() []Thread {
	g.mu.Lock()
	defer g.mu.Unlock()

	threads := make([]Thread, len(g.threads))
	copy(threads, g.threads)
	return threads
}
