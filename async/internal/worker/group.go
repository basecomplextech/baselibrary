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

type routineGroup interface {
	// len returns the number of active routines.
	len() int

	// start/stop

	// start starts the group.
	start()

	// stop stops all routines, waits for them to stop and clears the group.
	stop()

	// methods

	// add adds a routine to the group, or returns false if the group is stopped.
	add(r routine.RoutineVoid) bool

	// poll returns the next stopped routine or false.
	poll() (routine.RoutineVoid, bool)

	// wait waits for the next stopped routine.
	wait() <-chan struct{}
}

// internal

var _ routineGroup = (*group)(nil)

type group struct {
	// guards start/stop
	mainMu sync.Mutex

	// guards running and state
	runMu   sync.Mutex
	running bool

	// state can be accessed only when running is true
	active  sets.Set[routine.RoutineVoid]
	stopped queue.Queue[routine.RoutineVoid]
}

func newRoutineGroup() *group {
	return &group{
		active:  sets.New[routine.RoutineVoid](),
		stopped: queue.New[routine.RoutineVoid](),
	}
}

// len returns the number of active routines.
func (g *group) len() int {
	g.runMu.Lock()
	defer g.runMu.Unlock()

	if !g.running {
		return 0
	}

	return len(g.active)
}

// lifecycle

// start starts the group.
func (g *group) start() {
	g.mainMu.Lock()
	defer g.mainMu.Unlock()

	g.runMu.Lock()
	defer g.runMu.Unlock()

	g.running = true
}

// stop stops all routines, waits for them to stop and clears the group.
func (g *group) stop() {
	g.mainMu.Lock()
	defer g.mainMu.Unlock()

	// Disable running
	g.runMu.Lock()
	g.running = false
	g.runMu.Unlock()

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

// add adds a routine to the group, or returns false if the group is stopped.
func (g *group) add(r routine.RoutineVoid) bool {
	g.runMu.Lock()
	defer g.runMu.Unlock()

	// Check running
	if !g.running {
		panic("worker group is stopped")
	}

	// Add callback
	ok := r.OnStop(g.onStop)
	if !ok {
		g.stopped.Push(r)
		return false
	}

	// Add to active
	g.active.Add(r)
	return true
}

// poll returns the next stopped routine or false.
func (g *group) poll() (routine.RoutineVoid, bool) {
	g.runMu.Lock()
	defer g.runMu.Unlock()

	if !g.running {
		return nil, false
	}

	return g.stopped.Poll()
}

// wait waits for the next stopped routine.
func (g *group) wait() <-chan struct{} {
	g.runMu.Lock()
	defer g.runMu.Unlock()

	if !g.running {
		return chans.Closed()
	}

	return g.stopped.Wait()
}

// internal

// onStop moves a routine from active to stopped.
func (g *group) onStop(r routine.RoutineVoid) {
	g.runMu.Lock()
	defer g.runMu.Unlock()

	if !g.running {
		return
	}

	g.active.Remove(r)
	g.stopped.Push(r)
}
