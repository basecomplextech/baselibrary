// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package bin

import (
	"encoding/hex"
	"errors"
)

// Parse192 parses a bin192 from a 24-byte array.
func Parse192(b []byte) (Bin192, error) {
	switch {
	case b == nil:
		return Bin192{}, nil
	case len(b) == 0:
		return Bin192{}, nil
	case len(b) != Len192:
		return Bin192{}, errors.New("bin192: invalid bin192 length")
	}

	v := Bin192{}
	copy(v[0][:], b[:8])
	copy(v[1][:], b[8:16])
	copy(v[2][:], b[16:24])
	return v, nil
}

// ParseString192 parses a bin192 from 67-char string.
func ParseString192(s string) (Bin192, error) {
	switch {
	case s == "":
		return Bin192{}, nil
	case len(s) == 0:
		return Bin192{}, nil
	case len(s) != Len192Char:
		return Bin192{}, errors.New("bin192: invalid bin192 length")
	}

	b := unsafeByteString(s)
	p := [Len192]byte{}

	_, err := hex.Decode(p[:8], b[:16])
	if err != nil {
		return Bin192{}, err
	}
	_, err = hex.Decode(p[8:16], b[17:33])
	if err != nil {
		return Bin192{}, err
	}
	_, err = hex.Decode(p[16:24], b[34:50])
	if err != nil {
		return Bin192{}, err
	}
	return Parse192(p[:])
}

// MustParseString192 parses a bin192 from 67-char string or panics.
func MustParseString192(s string) Bin192 {
	v, err := ParseString192(s)
	if err != nil {
		panic(err)
	}
	return v
}
