package u128

import "regexp"

var Pattern = regexp.MustCompile(`^[0-9A-Za-z]{16}-[0-9A-Za-z]{16}$`)

// Match returns true if a byte string matches a U128 pattern.
func Match(s []byte) bool {
	return Pattern.Match(s)
}

// MatchString returns true if a string matches a U128 pattern.
func MatchString(s string) bool {
	return Pattern.MatchString(s)
}
