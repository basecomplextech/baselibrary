// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package worker

import (
	"sync"

	"github.com/basecomplextech/baselibrary/async/internal/queue"
	"github.com/basecomplextech/baselibrary/async/internal/routine"
	"github.com/basecomplextech/baselibrary/collect/chans"
	"github.com/basecomplextech/baselibrary/collect/sets"
)

type routines interface {
	// len returns the number of active routines.
	len() int

	// start/stop

	// start starts the group.
	start()

	// stop stops all routines, waits for them to stop and clears the group.
	stop()

	// methods

	// add adds a routine to the group, or panics if the group is stopped.
	add(r routine.RoutineVoid)

	// poll returns the next stopped routine or false.
	poll() (routine.RoutineVoid, bool)

	// wait waits for the next stopped routine.
	wait() <-chan struct{}
}

// internal

var _ routines = (*group)(nil)

type group struct {
	// guards start/stop
	mainMu sync.Mutex

	// state can be accessed only when handling is true
	mu       sync.Mutex
	handling bool
	active   sets.Set[routine.RoutineVoid]
	stopped  queue.Queue[routine.RoutineVoid]
}

func newRoutines() *group {
	return &group{
		active:  sets.New[routine.RoutineVoid](),
		stopped: queue.New[routine.RoutineVoid](),
	}
}

// len returns the number of active routines.
func (g *group) len() int {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.handling {
		return 0
	}
	return len(g.active)
}

// lifecycle

// start starts the group.
func (g *group) start() {
	g.mainMu.Lock()
	defer g.mainMu.Unlock()

	g.mu.Lock()
	defer g.mu.Unlock()

	g.handling = true
}

// stop stops all routines, waits for them to stop and clears the group.
func (g *group) stop() {
	g.mainMu.Lock()
	defer g.mainMu.Unlock()

	// Disable handling
	g.mu.Lock()
	g.handling = false
	g.mu.Unlock()

	// From now on, active and stopped cannot be
	// accessed concurrently by other threads.

	// Stop all
	for r := range g.active {
		r.Stop()
	}

	// Await all
	for r := range g.active {
		<-r.Wait()
	}

	// Clear all
	g.active.Clear()
	g.stopped.Clear()
}

// methods

// add adds a routine to the group, or panics if the group is stopped.
func (g *group) add(r routine.RoutineVoid) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Check handling
	if !g.handling {
		panic("worker group is stopped")
	}

	// Add callback
	ok := r.OnStop(g.onStop)
	if !ok {
		g.stopped.Push(r)
		return
	}

	// Add to active
	g.active.Add(r)
}

// poll returns the next stopped routine or false.
func (g *group) poll() (routine.RoutineVoid, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.handling {
		return nil, false
	}

	return g.stopped.Poll()
}

// wait waits for the next stopped routine.
func (g *group) wait() <-chan struct{} {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.handling {
		return chans.Closed()
	}

	return g.stopped.Wait()
}

// internal

// onStop moves a routine from active to stopped.
func (g *group) onStop(r routine.RoutineVoid) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.handling {
		return
	}

	g.active.Remove(r)
	g.stopped.Push(r)
}
