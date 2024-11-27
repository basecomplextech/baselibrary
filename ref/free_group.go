// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"github.com/basecomplextech/baselibrary/collect/slices2"
	"github.com/basecomplextech/baselibrary/pools"
)

// FreeGroup frees objects when the group is freed.
type FreeGroup interface {
	// Add an object to the free group.
	Add(obj Freer)

	// Free free the group and all objects in it.
	Free()
}

// NewFreeGroup returns a new unlimited free group.
func NewFreeGroup() FreeGroup {
	return newFreeGroup(0)
}

// NewFreeGroupCap returns a new free group with the given maximum capacity.
func NewFreeGroupCap(capacity int) FreeGroup {
	return newFreeGroup(capacity)
}

// internal

var _ (FreeGroup) = (*freeGroup)(nil)

type freeGroup struct {
	*freeGroupState
}

func newFreeGroup(maxCap int) *freeGroup {
	g := &freeGroup{acquireFreeGroupState()}
	g.maxCap = maxCap
	return g
}

type freeGroupState struct {
	maxCap  int
	objects []Freer
}

// Add an object to the free group.
func (g *freeGroup) Add(obj Freer) {
	if g.maxCap > 0 {
		if len(g.objects) >= g.maxCap {
			return
		}
	}

	g.objects = append(g.objects, obj)
}

// Free free the group and all objects in it.
func (g *freeGroup) Free() {
	s := g.freeGroupState
	g.freeGroupState = nil

	for _, obj := range s.objects {
		obj.Free()
	}
	s.objects = slices2.Truncate(s.objects)

	releaseFreeGroupState(s)
}

// pool

var freeGroupStatePool = pools.NewPoolFunc(
	func() *freeGroupState {
		return &freeGroupState{}
	},
)

func acquireFreeGroupState() *freeGroupState {
	return freeGroupStatePool.New()
}

func releaseFreeGroupState(s *freeGroupState) {
	freeGroupStatePool.Put(s)
}
