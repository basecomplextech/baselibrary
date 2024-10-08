// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package logging

import (
	"os"
	"sync"

	"github.com/basecomplextech/baselibrary/alloc"
	"github.com/basecomplextech/baselibrary/terminal"
)

var _ Writer = (*consoleWriter)(nil)

type consoleWriter struct {
	level  Level
	format *textFormatter

	mu  sync.Mutex
	out *os.File
}

func newConsoleWriter(level Level, color bool, out *os.File) *consoleWriter {
	if color {
		color = terminal.CheckColor(out)
	}

	return &consoleWriter{
		level:  level,
		format: newTextFormatter(color),

		out: out,
	}
}

func initConsoleWriter(config *ConsoleConfig) (*consoleWriter, error) {
	if !config.Enabled {
		return nil, nil
	}

	level := config.Level
	color := config.Color
	return newConsoleWriter(level, color, os.Stdout), nil
}

// Enabled returns true if the writer is enabled for the given level.
func (w *consoleWriter) Enabled(level Level) bool {
	return level >= w.level
}

// Write writes a record.
func (w *consoleWriter) Write(rec *Record) error {
	// Check level
	ok := w.Enabled(rec.Level)
	if !ok {
		return nil
	}

	// Format record
	buf := alloc.NewBuffer()
	defer buf.Free()

	if err := w.format.Format(buf, rec); err != nil {
		return err
	}
	msg := buf.Bytes()

	// Write record
	w.mu.Lock()
	defer w.mu.Unlock()

	_, err := w.out.Write(msg)
	return err
}
