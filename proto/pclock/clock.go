// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package pclock

import "fmt"

// IsZero returns true if this timestamp is zero.
func (t HLTimestamp) IsZero() bool {
	return t == HLTimestamp{}
}

// Compare this timestamp to another.
// The result is 0 if this == another, -1 if this < another, and 1 if this > another.
func (t HLTimestamp) Compare(t1 HLTimestamp) int {
	switch {
	case t.Wall == t1.Wall:
		switch {
		case t.Seq == t1.Seq:
			return 0
		case t.Seq < t1.Seq:
			return -1
		default:
			return 1
		}
	case t.Wall < t1.Wall:
		return -1
	default:
		return 1
	}
}

// Less returns true if this timestamp is less than another.
func (t HLTimestamp) Less(t1 HLTimestamp) bool {
	return t.Compare(t1) < 0
}

// LessOrEqual returns true if this timestamp is less than or equal to another.
func (t HLTimestamp) LessOrEqual(t1 HLTimestamp) bool {
	return t.Compare(t1) <= 0
}

// Min/Max

// Min returns the lesser of this timestamp and another.
func (t HLTimestamp) Min(t1 HLTimestamp) HLTimestamp {
	if t.Less(t1) {
		return t
	}
	return t1
}

// Max returns the greater of this timestamp and another.
func (t HLTimestamp) Max(t1 HLTimestamp) HLTimestamp {
	if t.Less(t1) {
		return t1
	}
	return t
}

// String

// String returns a "wall.seq" string.
func (t HLTimestamp) String() string {
	return fmt.Sprintf("%d.%d", t.Wall, t.Seq)
}
