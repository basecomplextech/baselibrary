package u256

import (
	"encoding/hex"
	"errors"
)

// Parse parses a U256 from a 32-byte array.
func Parse(b []byte) (U256, error) {
	switch {
	case b == nil:
		return U256{}, nil
	case len(b) == 0:
		return U256{}, nil
	case len(b) != ByteLen:
		return U256{}, errors.New("u256: invalid U256 length")
	}

	u := U256{}
	copy(u[:], b)
	return u, nil
}

// ParseString parses a U256 from 64-char string.
func ParseString(s string) (U256, error) {
	return ParseByteString([]byte(s))
}

// ParseByteString parses a U256 from 64-char byte string.
func ParseByteString(s []byte) (U256, error) {
	switch {
	case s == nil:
		return U256{}, nil
	case len(s) == 0:
		return U256{}, nil
	case len(s) != CharLen:
		return U256{}, errors.New("u256: invalid U256 length")
	}

	u := U256{}
	_, err := hex.Decode(u[:], s)
	if err != nil {
		return u, err
	}
	return u, nil
}

// MustParseString parses a U256 from 32-char string or panics.
func MustParseString(s string) U256 {
	u, err := ParseString(s)
	if err != nil {
		panic(err)
	}
	return u
}
