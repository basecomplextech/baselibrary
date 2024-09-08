// Copyright 2023 Ivan Korobkov. All rights reserved.

package logging

import (
	"os"
	"path/filepath"

	"github.com/basecomplextech/baselibrary/alloc"
	"gopkg.in/natefinch/lumberjack.v2"
)

var _ Writer = (*fileWriter)(nil)

type fileWriter struct {
	level  Level
	format *textFormatter

	out *lumberjack.Logger
}

func initFileWriter(config *FileConfig) (*fileWriter, error) {
	if !config.Enabled {
		return nil, nil
	}

	path := config.Path
	dir, _ := filepath.Split(path)
	if err := createCreateDir(dir); err != nil {
		return nil, err
	}
	if err := checkCreateFile(path); err != nil {
		return nil, err
	}

	out := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    config.MaxSize,
		MaxAge:     config.MaxAge,
		MaxBackups: config.MaxBackups,
		LocalTime:  true,
	}

	level := config.Level
	formatter := newTextFormatter(false /* color */)

	w := &fileWriter{
		level:  level,
		format: formatter,

		out: out,
	}
	return w, nil
}

// Enabled returns true if the writer is enabled for the given level.
func (w *fileWriter) Enabled(level Level) bool {
	return level >= w.level
}

// Write writes a record.
func (w *fileWriter) Write(rec *Record) error {
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
	_, err := w.out.Write(msg)
	return err
}

// private

func createCreateDir(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return os.MkdirAll(path, 0700)
	}
	return err
}

func checkCreateFile(path string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	if err == nil {
		return file.Close()
	}

	if os.IsNotExist(err) {
		file, err = os.Create(path)
		if err != nil {
			return err
		}
		return file.Close()
	}
	return err
}
