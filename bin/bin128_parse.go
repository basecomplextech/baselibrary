// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package bin

import (
	"encoding/binary"
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
	v[0] = Bin64(binary.BigEndian.Uint64(b))
	v[1] = Bin64(binary.BigEndian.Uint64(b[8:]))
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
	p := [Len128]byte{}

	_, err := hex.Decode(p[:8], b[:16])
	if err != nil {
		return Bin128{}, err
	}
	_, err = hex.Decode(p[8:], b[17:])
	if err != nil {
		return Bin128{}, err
	}
	return Parse128(p[:])
}

// MustParseString128 parses a bin128 from 33-char string or panics.
func MustParseString128(s string) Bin128 {
	v, err := ParseString128(s)
	if err != nil {
		panic(err)
	}
	return v
}
