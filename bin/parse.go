package bin

import (
	"encoding/hex"
	"errors"
)

// Parse64 parses a bin64 value from a 16-byte array.
func Parse64(b []byte) (Bin64, error) {
	switch {
	case b == nil:
		return Bin64{}, nil
	case len(b) == 0:
		return Bin64{}, nil
	case len(b) != ByteLen64:
		return Bin64{}, errors.New("bin64: invalid bin64 length")
	}

	u := Bin64{}
	copy(u[:], b)
	return u, nil
}

// ParseString64 parses a bin64 from 32-char string.
func ParseString64(s string) (Bin64, error) {
	switch {
	case s == "":
		return Bin64{}, nil
	case len(s) == 0:
		return Bin64{}, nil
	case len(s) != CharLen64:
		return Bin64{}, errors.New("bin64: invalid bin64 string length")
	}

	u := Bin64{}
	_, err := hex.Decode(u[:], []byte(s))
	if err != nil {
		return u, err
	}
	return u, nil
}

// MustParseString64 parses a bin64 from 16-char string or panics.
func MustParseString64(s string) Bin64 {
	u, err := ParseString64(s)
	if err != nil {
		panic(err)
	}
	return u
}

// Bin128

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

// ParseString128 parses a bin128 from 32-char string.
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
	_, err := hex.Decode(u[:], []byte(s))
	if err != nil {
		return u, err
	}
	return u, nil
}

// MustParseString128 parses a bin128 from 32-char string or panics.
func MustParseString128(s string) Bin128 {
	u, err := ParseString128(s)
	if err != nil {
		panic(err)
	}
	return u
}

// Bin256

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
	switch {
	case s == "":
		return Bin256{}, nil
	case len(s) == 0:
		return Bin256{}, nil
	case len(s) != CharLen256:
		return Bin256{}, errors.New("bin256: invalid bin256 length")
	}

	u := Bin256{}
	_, err := hex.Decode(u[:], []byte(s))
	if err != nil {
		return u, err
	}
	return u, nil
}

// MustParseString256 parses a bin256 from 64-char string or panics.
func MustParseString256(s string) Bin256 {
	u, err := ParseString256(s)
	if err != nil {
		panic(err)
	}
	return u
}
