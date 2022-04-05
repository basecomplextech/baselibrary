// package compactint provides a compact variable length integer encoding.
//
// Everything is big-endian.
//
// Encoding
//
//  Size                   Prefix  Example               Description
//	<= 0xfc                -       0x12                  uint8_t
//	<= 0xffff              0xfd    0xfd1234              0xfd followed by uint16
//	<= 0xffffffff          0xfe    0xfe12345678          0xfe followed by uint32
//	<= 0xffffffffffffffff  0xff    0xff1234567890abcdef  0xff followed by uint64
//
package compactint

import (
	"encoding/binary"
)

const MaxLen = 9

// Uint32 returns a uint32 value and the number of bytes read.
func Uint32(b []byte) (uint32, int) {
	if len(b) == 0 {
		return 0, 0
	}

	f := b[0]
	switch f {
	default:
		return uint32(f), 1

	case 0xfd:
		if len(b) < 3 {
			return 0, 0
		}

		v := binary.BigEndian.Uint16(b[1:3])
		return uint32(v), 3

	case 0xfe:
		if len(b) < 5 {
			return 0, 0
		}

		v := binary.BigEndian.Uint32(b[1:5])
		return uint32(v), 5

	case 0xff:
		return 0, -1
	}
}

// Uint64 returns a uint64 value and the number of bytes read.
func Uint64(b []byte) (uint64, int) {
	if len(b) == 0 {
		return 0, 0
	}

	f := b[0]
	switch f {
	default:
		return uint64(f), 1

	case 0xfd:
		if len(b) < 3 {
			return 0, 0
		}

		v := binary.BigEndian.Uint16(b[1:3])
		return uint64(v), 3

	case 0xfe:
		if len(b) < 5 {
			return 0, 0
		}

		v := binary.BigEndian.Uint32(b[1:5])
		return uint64(v), 5

	case 0xff:
		if len(b) < 9 {
			return 0, 0
		}

		v := binary.BigEndian.Uint64(b[1:9])
		return v, 9
	}
}

// PutUint32 appends a uint32 to b and returns the number of bytes written.
func PutUint32(b []byte, v uint32) int {
	switch {
	case v <= 0xfc:
		b[0] = byte(v)
		return 1

	case v <= 0xffff:
		b[0] = 0xfd
		binary.BigEndian.PutUint16(b[1:3], uint16(v))
		return 3

	default:
		b[0] = 0xfe
		binary.BigEndian.PutUint32(b[1:5], uint32(v))
		return 5
	}
}

// PutUint64 appends a uint64 to b and returns the number of bytes written.
func PutUint64(b []byte, v uint64) int {
	switch {
	case v <= 0xfc:
		b[0] = byte(v)
		return 1

	case v <= 0xffff:
		b[0] = 0xfd
		binary.BigEndian.PutUint16(b[1:3], uint16(v))
		return 3

	case v <= 0xffffffff:
		b[0] = 0xfe
		binary.BigEndian.PutUint32(b[1:5], uint32(v))
		return 5

	default:
		b[0] = 0xff
		binary.BigEndian.PutUint64(b[1:9], v)
		return 9
	}
}
