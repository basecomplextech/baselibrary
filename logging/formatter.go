package logging

import "io"

type Formatter interface {
	Format(w io.Writer, rec Record) error
}
