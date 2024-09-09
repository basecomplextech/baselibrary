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
	case len(b) != ByteLen64:
		return Bin64{}, errors.New("bin64: invalid bin64 length")
	}

	u := Bin64{}
	copy(u[:], b)
	return u, nil
}

// ParseString64 parses a bin64 from 32-char string.
func ParseString64(s string) (Bin64, error) {
	switch {
	case s == "":
		return Bin64{}, nil
	case len(s) == 0:
		return Bin64{}, nil
	case len(s) != CharLen64:
		return Bin64{}, errors.New("bin64: invalid bin64 string length")
	}

	u := Bin64{}
	_, err := hex.Decode(u[:], []byte(s))
	if err != nil {
		return u, err
	}
	return u, nil
}

// MustParseString64 parses a bin64 from 16-char string or panics.
func MustParseString64(s string) Bin64 {
	u, err := ParseString64(s)
	if err != nil {
		panic(err)
	}
	return u
}
