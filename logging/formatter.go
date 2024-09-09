// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package logging

import "io"

type Formatter interface {
	Format(w io.Writer, rec *Record) error
}
