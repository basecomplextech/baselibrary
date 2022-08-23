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

// ParseString256 parses a bin256 from 64-char string.
func ParseString256(s string) (Bin256, error) {
	return ParseByteString256([]byte(s))
}

// ParseByteString256 parses a bin256 from 64-char byte string.
func ParseByteString256(s []byte) (Bin256, error) {
	switch {
	case s == nil:
		return Bin256{}, nil
	case len(s) == 0:
		return Bin256{}, nil
	case len(s) != CharLen256:
		return Bin256{}, errors.New("bin256: invalid bin256 length")
	}

	u := Bin256{}
	_, err := hex.Decode(u[:], s)
	if err != nil {
		return u, err
	}
	return u, nil
}

// MustParseString256 parses a bin256 from 32-char string or panics.
func MustParseString256(s string) Bin256 {
	u, err := ParseString256(s)
	if err != nil {
		panic(err)
	}
	return u
}
