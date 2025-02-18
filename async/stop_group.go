// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"sync"
	"time"
)

// StopGroup stops all operations in the group and awaits their completion.
type StopGroup struct {
	mu   sync.Mutex
	done bool

	stoppers []Stopper
	timers   []*time.Timer
	tickers  []*time.Ticker
}

// NewStopGroup creates a new stop group.
func NewStopGroup() *StopGroup {
	return &StopGroup{}
}

// Add adds a stopper to the group, or immediately stops it if the group is stopped.
func (g *StopGroup) Add(s Stopper) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.done {
		s.Stop()
		return
	}

	g.stoppers = append(g.stoppers, s)
}

// AddMany adds multiple stoppers to the group, or immediatelly stops them if the group is stopped.
func (g *StopGroup) AddMany(stoppers ...Stopper) {
	for _, s := range stoppers {
		g.Add(s)
	}
}

// AddService adds a service to the group, or immediately stops it if the group is stopped.
func (g *StopGroup) AddService(s Service) {
	g.Add(s)
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

// Stop stops all operations in the group.
func (g *StopGroup) Stop() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.done {
		return
	}
	g.done = true

	for _, w := range g.stoppers {
		w.Stop()
	}
	for _, t := range g.tickers {
		t.Stop()
	}
	for _, t := range g.timers {
		t.Stop()
	}
}

// StopWait stops all operations in the group and awaits them.
func (g *StopGroup) StopWait() {
	g.Stop()
	g.Wait()
}

// Wait awaits all operations in the group.
func (g *StopGroup) Wait() {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, w := range g.stoppers {
		<-w.Wait()
	}
}
