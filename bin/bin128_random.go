// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package bin

import (
	"encoding/binary"
	"time"
)

// Random128 returns a random bin128.
func Random128() Bin128 {
	p := random.read128()

	b := Bin128{}
	copy(b[0][:], p[:8])
	copy(b[1][:], p[8:])
	return b
}

// TimeRandom128 returns a time-random bin128 with a millisecond resolution.
func TimeRandom128() Bin128 {
	now := time.Now().UnixNano() / int64(time.Millisecond)

	b := Random128()
	binary.BigEndian.PutUint64(b[0][:], uint64(now))
	return b
}
