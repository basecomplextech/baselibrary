package logging

import (
	"github.com/complex1tech/baselibrary/slices"
)

// Writer writes log records.
type Writer interface {
	// Enabled returns true if the writer is enabled for the given level.
	Enabled(level Level) bool

	// Write writes a record.
	Write(rec Record) error
}

// internal

// var _ Writer = (*writer)(nil)

// const lflags = log.Ldate | log.Ltime | log.Lmicroseconds

// type writer struct {
// 	level     Level
// 	logger    *log.Logger
// 	formatter Formatter
// }

// func newWriter(level Level, logger *log.Logger, formatter Formatter) *writer {
// 	return &writer{
// 		level:     level,
// 		logger:    logger,
// 		formatter: formatter,
// 	}
// }

// // Write writes a record.
// func (w *writer) Write(rec Record) error {
// 	if rec.Level < w.level {
// 		return nil
// 	}

// 	msg := w.formatter.Format(rec)
// 	return w.logger.Output(2, msg)
// }

// null

var _ Writer = (*nullWriter)(nil)

type nullWriter struct{}

func newNullWriter() *nullWriter {
	return &nullWriter{}
}

// Enabled returns true if the writer is enabled for the given level.
func (w *nullWriter) Enabled(level Level) bool {
	return false
}

// Write writes a record.
func (w *nullWriter) Write(rec Record) error {
	return nil
}

// multi

var _ Writer = (*multiWriter)(nil)

type multiWriter struct {
	writers []Writer
}

func newMultiWriter(writers ...Writer) *multiWriter {
	return &multiWriter{
		writers: slices.Clone(writers),
	}
}

// Enabled returns true if the writer is enabled for the given level.
func (w *multiWriter) Enabled(level Level) bool {
	for _, w := range w.writers {
		if w.Enabled(level) {
			return true
		}
	}
	return false
}

// Write writes a record.
func (w *multiWriter) Write(rec Record) error {
	for _, w := range w.writers {
		if err := w.Write(rec); err != nil {
			return err
		}
	}
	return nil
}
