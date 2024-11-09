// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package context

import "github.com/basecomplextech/baselibrary/status"

var no Context = &noContext{}

type noContext struct{}

func (*noContext) Cancel()                 {}
func (*noContext) Done() bool              { return false }
func (*noContext) Wait() <-chan struct{}   { return nil }
func (*noContext) Status() status.Status   { return status.OK }
func (*noContext) AddCallback(Callback)    {}
func (*noContext) RemoveCallback(Callback) {}
func (*noContext) Free()                   {}
