// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package stop

import (
	"sync"
	"time"

	"github.com/basecomplextech/baselibrary/async/internal/routine"
	"github.com/basecomplextech/baselibrary/async/internal/service"
)

// StopGroup stops all operations in the group and awaits their completion.
type StopGroup struct {
	mu   sync.Mutex
	done bool

	routines []routine.RoutineDyn
	services []service.Service
	timers   []*time.Timer
	tickers  []*time.Ticker
}

// NewStopGroup creates a new stop group.
func NewStopGroup() *StopGroup {
	return &StopGroup{}
}

// Add adds a routine to the group, or immediately stops it if the group is stopped.
func (g *StopGroup) Add(r routine.RoutineDyn) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.done {
		r.Stop()
		return
	}

	g.routines = append(g.routines, r)
}

// AddMany adds multiple routines to the group, or immediatelly stops them if the group is stopped.
func (g *StopGroup) AddMany(routines ...routine.RoutineDyn) {
	for _, r := range routines {
		g.Add(r)
	}
}

// AddService adds a service to the group, or immediately stops it if the group is stopped.
func (g *StopGroup) AddService(s service.Service) {
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

// Stop stops all operations in the group.
func (g *StopGroup) Stop() {
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
	for _, w := range g.routines {
		w.Stop()
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

	for _, s := range g.services {
		<-s.Wait()
	}

	for _, w := range g.routines {
		<-w.Wait()
	}
}
