// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package bin

import (
	"encoding/binary"
	"time"
)

// Random64 returns a random bin64.
func Random64() Bin64 {
	p := random.read64()
	v := binary.BigEndian.Uint64(p[:])
	return Bin64(v)
}

// TimeRandom64 returns a time-random bin64 with a second resolution.
func TimeRandom64() Bin64 {
	p := random.read64()

	now := time.Now().UnixNano() / int64(time.Second)
	binary.BigEndian.PutUint32(p[:], uint32(now))

	v := binary.BigEndian.Uint64(p[:])
	return Bin64(v)
}
