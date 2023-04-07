package bin

import (
	"crypto/rand"
	"encoding/binary"
	"time"
)

// Random64 returns a random bin64.
func Random64() Bin64 {
	u := Bin64{}
	if _, err := rand.Read(u[:]); err != nil {
		panic(err)
	}
	return u
}

// Random128 returns a random bin128.
func Random128() Bin128 {
	u := Bin128{}
	if _, err := rand.Read(u[:]); err != nil {
		panic(err)
	}
	return u
}

// Random256 returns a random bin256.
func Random256() Bin256 {
	u := Bin256{}
	if _, err := rand.Read(u[:]); err != nil {
		panic(err)
	}
	return u
}

// Time-random

// TimeRandom64 returns a time-random bin64 with a second resolution.
func TimeRandom64() Bin64 {
	u := Bin64{}

	now := time.Now()
	ts := now.UnixNano() / int64(time.Second)
	binary.BigEndian.PutUint32(u[:], uint32(ts))

	if _, err := rand.Read(u[4:]); err != nil {
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
