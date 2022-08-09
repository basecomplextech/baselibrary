package bin128

import (
	"encoding/hex"
	"errors"
)

// Parse parses a B128 ID from a 16-byte array.
func Parse(b []byte) (B128, error) {
	switch {
	case b == nil:
		return B128{}, nil
	case len(b) == 0:
		return B128{}, nil
	case len(b) != ByteLen:
		return B128{}, errors.New("bin128: invalid B128 length")
	}

	u := B128{}
	copy(u[:], b)
	return u, nil
}

// ParseString parses a B128 from 33-char string.
func ParseString(s string) (B128, error) {
	return ParseByteString([]byte(s))
}

// ParseString parses a B128 from 32-char string.
func ParseByteString(s []byte) (B128, error) {
	switch {
	case s == nil:
		return B128{}, nil
	case len(s) == 0:
		return B128{}, nil
	case len(s) != CharLen:
		return B128{}, errors.New("bin128: invalid B128 string length")
	}

	u := B128{}
	_, err := hex.Decode(u[:], s)
	if err != nil {
		return u, err
	}
	return u, nil
}
