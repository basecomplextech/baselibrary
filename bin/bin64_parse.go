// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package bin

import (
	"encoding/hex"
	"errors"
)

// Parse64 parses a bin64 value from a 16-byte array.
func Parse64(b []byte) (Bin64, error) {
	switch {
	case b == nil:
		return Bin64{}, nil
	case len(b) == 0:
		return Bin64{}, nil
	case len(b) != Len64:
		return Bin64{}, errors.New("bin64: invalid bin64 length")
	}

	v := Bin64{}
	copy(v[:], b)
	return v, nil
}

// ParseString64 parses a bin64 from 32-char string.
func ParseString64(s string) (Bin64, error) {
	switch {
	case s == "":
		return Bin64{}, nil
	case len(s) == 0:
		return Bin64{}, nil
	case len(s) != Len64Char:
		return Bin64{}, errors.New("bin64: invalid bin64 string length")
	}

	v := Bin64{}
	_, err := hex.Decode(v[:], []byte(s))
	if err != nil {
		return v, err
	}
	return v, nil
}

// MustParseString64 parses a bin64 from 16-char string or panics.
func MustParseString64(s string) Bin64 {
	v, err := ParseString64(s)
	if err != nil {
		panic(err)
	}
	return v
}
