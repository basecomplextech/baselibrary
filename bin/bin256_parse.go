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
	case len(b) != Len256:
		return Bin256{}, errors.New("bin256: invalid bin256 length")
	}

	v := Bin256{}
	copy(v[0][:], b)
	copy(v[1][:], b[8:])
	copy(v[2][:], b[16:])
	copy(v[3][:], b[24:])
	return v, nil
}

// ParseString256 parses a bin256 from 67-char string.
func ParseString256(s string) (Bin256, error) {
	switch {
	case s == "":
		return Bin256{}, nil
	case len(s) == 0:
		return Bin256{}, nil
	case len(s) != Len256Char:
		return Bin256{}, errors.New("bin256: invalid bin256 length")
	}

	v := Bin256{}
	b := unsafeByteString(s)

	_, err := hex.Decode(v[0][:], b[:16])
	if err != nil {
		return v, err
	}
	_, err = hex.Decode(v[1][:], b[17:33])
	if err != nil {
		return v, err
	}
	_, err = hex.Decode(v[2][:], b[34:50])
	if err != nil {
		return v, err
	}
	_, err = hex.Decode(v[3][:], b[51:])
	if err != nil {
		return v, err
	}
	return v, nil
}

// MustParseString256 parses a bin256 from 67-char string or panics.
func MustParseString256(s string) Bin256 {
	v, err := ParseString256(s)
	if err != nil {
		panic(err)
	}
	return v
}
