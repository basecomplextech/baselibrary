package bin128

import "regexp"

var Pattern = regexp.MustCompile(`^[0-9A-Za-z]{32}$`)

// Match returns true if a byte string matches a B128 pattern.
func Match(s []byte) bool {
	return Pattern.Match(s)
}

// MatchString returns true if a string matches a B128 pattern.
func MatchString(s string) bool {
	return Pattern.MatchString(s)
}
