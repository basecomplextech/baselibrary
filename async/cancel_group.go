package async

import (
	"sync"
	"time"
)

// CancelGroup cancels all operations in the group and awaits their completion.
type CancelGroup struct {
	mu   sync.Mutex
	done bool

	services []Service
	timers   []*time.Timer
	tickers  []*time.Ticker
	waiters  []CancelWaiter
}

// NewCancelGroup creates a new cancel group.
func NewCancelGroup() *CancelGroup {
	return &CancelGroup{}
}

// Add adds a canceller to the group, or immediately cancels it if the group is cancelled.
func (g *CancelGroup) Add(c CancelWaiter) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.done {
		c.Cancel()
		return
	}

	g.waiters = append(g.waiters, c)
}

// AddMany adds multiple cancellers to the group, or immediatelly cancels them if the group is cancelled.
func (g *CancelGroup) AddMany(c ...CancelWaiter) {
	for _, c := range c {
		g.Add(c)
	}
}

// AddService adds a service to the group, or immediately stops it if the group is cancelled.
func (g *CancelGroup) AddService(s Service) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.done {
		s.Stop()
		return
	}

	g.services = append(g.services, s)
}

// AddTicker adds a ticker to the group, or immediately stops it if the group is cancelled.
func (g *CancelGroup) AddTicker(t *time.Ticker) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.done {
		t.Stop()
		return
	}

	g.tickers = append(g.tickers, t)
}

// AddTimer adds a timer to the group, or immediately cancels it if the group is cancelled.
func (g *CancelGroup) AddTimer(t *time.Timer) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.done {
		t.Stop()
		return
	}

	g.timers = append(g.timers, t)
}

// Cancel cancels all operations in the group.
func (g *CancelGroup) Cancel() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.done {
		return
	}
	g.done = true

	for _, s := range g.services {
		s.Stop()
	}
	for _, t := range g.tickers {
		t.Stop()
	}
	for _, t := range g.timers {
		t.Stop()
	}
	for _, w := range g.waiters {
		w.Cancel()
	}
}

// CancelWait cancels all operations in the group and awaits them.
func (g *CancelGroup) CancelWait() {
	g.Cancel()
	g.Wait()
}

// Wait awaits all operations in the group.
func (g *CancelGroup) Wait() {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, s := range g.services {
		<-s.Routine().Wait()
	}

	for _, w := range g.waiters {
		<-w.Wait()
	}
}
