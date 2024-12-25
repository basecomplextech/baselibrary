// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package bin

import (
	"encoding/hex"
	"errors"
)

// Parse64 parses a bin64 value from a 16-byte array.
func Parse64(p []byte) (Bin64, error) {
	switch {
	case p == nil:
		return Bin64{}, nil
	case len(p) == 0:
		return Bin64{}, nil
	case len(p) != Len64:
		return Bin64{}, errors.New("bin64: invalid bin64 length")
	}

	b := Bin64{}
	copy(b[:], p)
	return b, nil
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

	b := Bin64{}
	_, err := hex.Decode(b[:], []byte(s))
	if err != nil {
		return Bin64{}, err
	}
	return b, nil
}

// MustParseString64 parses a bin64 from 16-char string or panics.
func MustParseString64(s string) Bin64 {
	v, err := ParseString64(s)
	if err != nil {
		panic(err)
	}
	return v
}
