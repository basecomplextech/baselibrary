// Copyright 2022 Ivan Korobkov. All rights reserved.

package logging

import "io"

type Formatter interface {
	Format(w io.Writer, rec *Record) error
}
