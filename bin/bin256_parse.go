// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package bin

import (
	"encoding/hex"
	"errors"
)

// Parse256 parses a bin256 from a 32-byte array.
func Parse256(b []byte) (Bin256, error) {
	switch {
	case b == nil:
		return Bin256{}, nil
	case len(b) == 0:
		return Bin256{}, nil
	case len(b) != ByteLen256:
		return Bin256{}, errors.New("bin256: invalid bin256 length")
	}

	u := Bin256{}
	copy(u[:], b)
	return u, nil
}

// ParseString256 parses a bin256 from 67-char string.
func ParseString256(s string) (Bin256, error) {
	switch {
	case s == "":
		return Bin256{}, nil
	case len(s) == 0:
		return Bin256{}, nil
	case len(s) != CharLen256:
		return Bin256{}, errors.New("bin256: invalid bin256 length")
	}

	u := Bin256{}
	b := unsafeByteString(s)

	_, err := hex.Decode(u[:], b[:16])
	if err != nil {
		return u, err
	}
	_, err = hex.Decode(u[8:], b[17:33])
	if err != nil {
		return u, err
	}
	_, err = hex.Decode(u[16:], b[34:50])
	if err != nil {
		return u, err
	}
	_, err = hex.Decode(u[24:], b[51:])
	if err != nil {
		return u, err
	}
	return u, nil
}

// MustParseString256 parses a bin256 from 67-char string or panics.
func MustParseString256(s string) Bin256 {
	u, err := ParseString256(s)
	if err != nil {
		panic(err)
	}
	return u
}
