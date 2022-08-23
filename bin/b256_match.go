package bin

import "regexp"

var Pattern256 = regexp.MustCompile(`^[0-9A-Za-z]{64}$`)

// Match256 returns true if a byte string matches a bin256 pattern.
func Match256(s []byte) bool {
	return Pattern256.Match(s)
}

// MatchString256 returns true if a string matches a bin256 pattern.
func MatchString256(s string) bool {
	return Pattern256.MatchString(s)
}
