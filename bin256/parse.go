package bin256

import (
	"encoding/hex"
	"errors"
)

// Parse parses a B256 from a 32-byte array.
func Parse(b []byte) (B256, error) {
	switch {
	case b == nil:
		return B256{}, nil
	case len(b) == 0:
		return B256{}, nil
	case len(b) != ByteLen:
		return B256{}, errors.New("bin256: invalid B256 length")
	}

	u := B256{}
	copy(u[:], b)
	return u, nil
}

// ParseString parses a B256 from 64-char string.
func ParseString(s string) (B256, error) {
	return ParseByteString([]byte(s))
}

// ParseByteString parses a B256 from 64-char byte string.
func ParseByteString(s []byte) (B256, error) {
	switch {
	case s == nil:
		return B256{}, nil
	case len(s) == 0:
		return B256{}, nil
	case len(s) != CharLen:
		return B256{}, errors.New("bin256: invalid B256 length")
	}

	u := B256{}
	_, err := hex.Decode(u[:], s)
	if err != nil {
		return u, err
	}
	return u, nil
}

// MustParseString parses a B256 from 32-char string or panics.
func MustParseString(s string) B256 {
	u, err := ParseString(s)
	if err != nil {
		panic(err)
	}
	return u
}
