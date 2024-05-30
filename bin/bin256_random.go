package bin

import (
	"crypto/rand"
	"encoding/binary"
	"time"
)

// Random256 returns a random bin256.
func Random256() Bin256 {
	u := Bin256{}
	if _, err := rand.Read(u[:]); err != nil {
		panic(err)
	}
	return u
}

// TimeRandom256 returns a time-random bin256 with a millisecond resolution.
func TimeRandom256() Bin256 {
	u := Bin256{}

	now := time.Now()
	ts := now.UnixNano() / int64(time.Millisecond)
	binary.BigEndian.PutUint64(u[:], uint64(ts))

	if _, err := rand.Read(u[8:]); err != nil {
		panic(err)
	}
	return u
}
