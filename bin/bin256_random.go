package bin

import (
	"encoding/binary"
	"time"
)

// Random256 returns a random bin256.
func Random256() Bin256 {
	return random.read256()
}

// TimeRandom256 returns a time-random bin256 with a millisecond resolution.
func TimeRandom256() Bin256 {
	b := random.read256()

	now := time.Now()
	timestamp := now.UnixNano() / int64(time.Millisecond)
	binary.BigEndian.PutUint64(b[:], uint64(timestamp))
	return b
}
