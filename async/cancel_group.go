package async

import (
	"sync"
	"time"
)

// CancelGroup cancels all operations in the group and awaits their completion.
type CancelGroup struct {
	mu   sync.Mutex
	done bool

	cancels  []CancelWaiter
	services []Service
	tickers  []*time.Ticker
	timers   []*time.Timer
}

// NewCancelGroup creates a new cancel group.
func NewCancelGroup() *CancelGroup {
	return &CancelGroup{}
}

// Add adds a routine to the group, or immediately cancels it if the group is cancelled.
func (g *CancelGroup) Add(c CancelWaiter) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.done {
		c.Cancel()
		return
	}

	g.cancels = append(g.cancels, c)
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

// Cancel cancels all operations in the group.
func (g *CancelGroup) Cancel() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.done {
		return
	}
	g.done = true

	for _, c := range g.cancels {
		c.Cancel()
	}
	for _, s := range g.services {
		s.Stop()
	}
	for _, t := range g.tickers {
		t.Stop()
	}
	for _, t := range g.timers {
		t.Stop()
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

	for _, c := range g.cancels {
		<-c.Wait()
	}

	for _, s := range g.services {
		<-s.Stopped()
	}
}
