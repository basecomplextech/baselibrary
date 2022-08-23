package bin

import "regexp"

var Pattern128 = regexp.MustCompile(`^[0-9A-Za-z]{32}$`)

// Match128 returns true if a byte string matches a bin128 pattern.
func Match128(s []byte) bool {
	return Pattern128.Match(s)
}

// MatchString128 returns true if a string matches a bin128 pattern.
func MatchString128(s string) bool {
	return Pattern128.MatchString(s)
}
