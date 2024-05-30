package bin

import (
	"encoding/binary"
	"time"
)

// Random64 returns a random bin64.
func Random64() Bin64 {
	return random.read64()
}

// TimeRandom64 returns a time-random bin64 with a second resolution.
func TimeRandom64() Bin64 {
	b := random.read64()

	now := time.Now()
	timestamp := now.UnixNano() / int64(time.Second)
	binary.BigEndian.PutUint32(b[:], uint32(timestamp))
	return b
}
