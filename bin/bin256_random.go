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
	b[0] = Bin64(binary.BigEndian.Uint64(p[:8]))
	b[1] = Bin64(binary.BigEndian.Uint64(p[8:]))
	b[2] = Bin64(binary.BigEndian.Uint64(p[16:]))
	b[3] = Bin64(binary.BigEndian.Uint64(p[24:]))
	return b
}

// TimeRandom256 returns a time-random bin256 with a millisecond resolution.
func TimeRandom256() Bin256 {
	now := time.Now().UnixNano() / int64(time.Millisecond)

	b := Random256()
	b[0] = Bin64(uint64(now))
	return b
}
