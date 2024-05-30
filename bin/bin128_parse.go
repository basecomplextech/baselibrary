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
	case len(b) != ByteLen128:
		return Bin128{}, errors.New("bin128: invalid bin128 length")
	}

	u := Bin128{}
	copy(u[:], b)
	return u, nil
}

// ParseString128 parses a bin128 from 33-char string.
func ParseString128(s string) (Bin128, error) {
	switch {
	case s == "":
		return Bin128{}, nil
	case len(s) == 0:
		return Bin128{}, nil
	case len(s) != CharLen128:
		return Bin128{}, errors.New("bin128: invalid bin128 string length")
	}

	u := Bin128{}
	b := unsafeByteString(s)

	_, err := hex.Decode(u[:], b[:16])
	if err != nil {
		return u, err
	}
	_, err = hex.Decode(u[8:], b[17:])
	if err != nil {
		return u, err
	}
	return u, nil
}

// MustParseString128 parses a bin128 from 33-char string or panics.
func MustParseString128(s string) Bin128 {
	u, err := ParseString128(s)
	if err != nil {
		panic(err)
	}
	return u
}
