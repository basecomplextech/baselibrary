// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package bin

import (
	"encoding/binary"
	"time"
)

// Random256 returns a random bin256.
func Random256() Bin256 {
	p := random.read256()

	b := Bin256{}
	copy(b[0][:], p[:8])
	copy(b[1][:], p[8:16])
	copy(b[2][:], p[16:24])
	copy(b[3][:], p[24:])
	return b
}

// TimeRandom256 returns a time-random bin256 with a millisecond resolution.
func TimeRandom256() Bin256 {
	now := time.Now().UnixNano() / int64(time.Millisecond)

	b := Random256()
	binary.BigEndian.PutUint64(b[0][:], uint64(now))
	return b
}
