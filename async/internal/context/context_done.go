// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package context

import (
	"github.com/basecomplextech/baselibrary/collect/chans"
	"github.com/basecomplextech/baselibrary/status"
)

var done MutContext = &doneContext{}

type doneContext struct{}

func (*doneContext) Cancel()                    {}
func (*doneContext) Done() bool                 { return true }
func (*doneContext) Wait() <-chan struct{}      { return chans.Closed() }
func (*doneContext) Status() status.Status      { return status.OK }
func (*doneContext) AddCallback(cb Callback)    { cb.OnCancelled(status.Cancelled) }
func (*doneContext) RemoveCallback(cb Callback) {}
func (*doneContext) Free()                      {}
