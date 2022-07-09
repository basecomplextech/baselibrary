package errors2

import (
	"fmt"
	"strings"
)

// MultiError wraps multiple nonnil errors into one error.
type MultiError struct {
	Errors []error
	Format string
}

// Combines combines multiple nonnil errors into a *MultiError and returns it or nil.
func Combine(err ...error) error {
	return Combinef("%v", err...)
}

// Combines combines multiple nonnil errors into a *MultiError and returns it or nil.
func Combinef(format string, err ...error) error {
	if len(err) == 0 {
		return nil
	}

	// filter nonnil errors in place
	var nonnil = err[:0]
	for _, e := range err {
		if e == nil {
			continue
		}
		nonnil = append(nonnil, e)
	}

	// return nil or the only error
	switch len(nonnil) {
	case 0:
		return nil
	case 1:
		return nonnil[0]
	}

	// return multi error
	return &MultiError{
		Errors: nonnil,
		Format: format,
	}
}

func (e *MultiError) Error() string {
	ss := make([]string, 0, len(e.Errors))
	for _, err := range e.Errors {
		ss = append(ss, err.Error())
	}

	format := e.Format
	if format == "" {
		format = "%v"
	}

	s := strings.Join(ss, ", ")
	return fmt.Sprintf(format, s)
}
