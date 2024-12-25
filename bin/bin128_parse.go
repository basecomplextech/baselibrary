// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package bin

import (
	"encoding/hex"
	"errors"
)

// Parse128 parses a bin128 value from a 16-byte array.
func Parse128(b []byte) (Bin128, error) {
	switch {
	case b == nil:
		return Bin128{}, nil
	case len(b) == 0:
		return Bin128{}, nil
	case len(b) != Len128:
		return Bin128{}, errors.New("bin128: invalid bin128 length")
	}

	v := Bin128{}
	copy(v[0][:], b[:8])
	copy(v[1][:], b[8:])
	return v, nil
}

// ParseString128 parses a bin128 from 33-char string.
func ParseString128(s string) (Bin128, error) {
	switch {
	case s == "":
		return Bin128{}, nil
	case len(s) == 0:
		return Bin128{}, nil
	case len(s) != Len128Char:
		return Bin128{}, errors.New("bin128: invalid bin128 string length")
	}

	b := unsafeByteString(s)
	p := Bin128{}

	_, err := hex.Decode(p[0][:], b[:16])
	if err != nil {
		return Bin128{}, err
	}
	_, err = hex.Decode(p[1][:], b[17:])
	if err != nil {
		return Bin128{}, err
	}
	return p, nil
}

// MustParseString128 parses a bin128 from 33-char string or panics.
func MustParseString128(s string) Bin128 {
	v, err := ParseString128(s)
	if err != nil {
		panic(err)
	}
	return v
}
