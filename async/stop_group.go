package async

import (
	"sync"
	"time"
)

// StopGroup stops all operations in the group and awaits their completion.
type StopGroup struct {
	mu   sync.Mutex
	done bool

	services []Service
	timers   []*time.Timer
	tickers  []*time.Ticker
	waiters  []StopWaiter
}

// NewStopGroup creates a new stop group.
func NewStopGroup() *StopGroup {
	return &StopGroup{}
}

// Add adds a routine to the group, or immediately stops it if the group is stopped.
func (g *StopGroup) Add(c StopWaiter) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.done {
		c.Stop()
		return
	}

	g.waiters = append(g.waiters, c)
}

// AddMany adds multiple routines to the group, or immediatelly stops them if the group is stopped.
func (g *StopGroup) AddMany(c ...StopWaiter) {
	for _, c := range c {
		g.Add(c)
	}
}

// AddService adds a service to the group, or immediately stops it if the group is stopped.
func (g *StopGroup) AddService(s Service) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.done {
		s.Stop()
		return
	}

	g.services = append(g.services, s)
}

// AddTicker adds a ticker to the group, or immediately stops it if the group is stopped.
func (g *StopGroup) AddTicker(t *time.Ticker) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.done {
		t.Stop()
		return
	}

	g.tickers = append(g.tickers, t)
}

// AddTimer adds a timer to the group, or immediately stops it if the group is stopped.
func (g *StopGroup) AddTimer(t *time.Timer) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.done {
		t.Stop()
		return
	}

	g.timers = append(g.timers, t)
}

// Cancel stops all operations in the group.
func (g *StopGroup) Cancel() {
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
		w.Stop()
	}
}

// StopWait stops all operations in the group and awaits them.
func (g *StopGroup) StopWait() {
	g.Cancel()
	g.Wait()
}

// Wait awaits all operations in the group.
func (g *StopGroup) Wait() {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, s := range g.services {
		<-s.Wait()
	}

	for _, w := range g.waiters {
		<-w.Wait()
	}
}
