package bin

import (
	"crypto/rand"
	"encoding/binary"
	"time"
)

// Random128 returns a random bin128.
func Random128() Bin128 {
	u := Bin128{}
	if _, err := rand.Read(u[:]); err != nil {
		panic(err)
	}
	return u
}

// TimeRandom128 returns a time-random bin128 with a millisecond resolution.
func TimeRandom128() Bin128 {
	u := Bin128{}

	now := time.Now()
	ts := now.UnixNano() / int64(time.Millisecond)
	binary.BigEndian.PutUint64(u[:], uint64(ts))

	if _, err := rand.Read(u[8:]); err != nil {
		panic(err)
	}
	return u
}
