// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package pclock

import "time"

func TestHLTimestamp() HLTimestamp {
	now := time.Now().UnixNano()

	return HLTimestamp{
		Wall: int64(now),
	}
}
