// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package bin

import (
	"encoding/binary"
	"time"
)

// Random192 returns a random bin192.
func Random192() Bin192 {
	p := random.read192()

	b := Bin192{}
	copy(b[0][:], p[:8])
	copy(b[1][:], p[8:16])
	copy(b[2][:], p[16:24])
	return b
}

// TimeRandom192 returns a time-random bin192 with a millisecond resolution.
func TimeRandom192() Bin192 {
	now := time.Now().UnixNano() / int64(time.Millisecond)

	b := Random192()
	binary.BigEndian.PutUint64(b[0][:], uint64(now))
	return b
}
