// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import "github.com/basecomplextech/baselibrary/async/internal/stop"

// StopGroup stops all operations in the group and awaits their completion.
type StopGroup = stop.StopGroup

// NewStopGroup creates a new stop group.
func NewStopGroup() StopGroup {
	return *stop.NewStopGroup()
}
