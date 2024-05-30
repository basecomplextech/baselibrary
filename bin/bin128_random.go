package bin

import (
	"encoding/binary"
	"time"
)

// Random128 returns a random bin128.
func Random128() Bin128 {
	return random.read128()
}

// TimeRandom128 returns a time-random bin128 with a millisecond resolution.
func TimeRandom128() Bin128 {
	b := random.read128()

	now := time.Now()
	timestamp := now.UnixNano() / int64(time.Millisecond)
	binary.BigEndian.PutUint64(b[:], uint64(timestamp))
	return b
}
