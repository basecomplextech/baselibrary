package u128

import (
	"encoding/hex"
	"errors"
)

// Parse parses a U128 ID from a 16-byte array.
func Parse(b []byte) (U128, error) {
	switch {
	case b == nil:
		return U128{}, nil
	case len(b) == 0:
		return U128{}, nil
	case len(b) != byteLen:
		return U128{}, errors.New("u128: invalid U128 length")
	}

	u := U128{}
	copy(u[:], b)
	return u, nil
}

// ParseString parses a U128 from 33-char string.
func ParseString(s string) (U128, error) {
	return ParseByteString([]byte(s))
}

// ParseString parses a U128 from 32-char string.
func ParseByteString(s []byte) (U128, error) {
	switch {
	case s == nil:
		return U128{}, nil
	case len(s) == 0:
		return U128{}, nil
	case len(s) != charLen:
		return U128{}, errors.New("u128: invalid U128 string length")
	}

	u := U128{}
	_, err := hex.Decode(u[:8], s[:16])
	if err != nil {
		return u, err
	}
	_, err = hex.Decode(u[8:], s[17:])
	return u, err
}
