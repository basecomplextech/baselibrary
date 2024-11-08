// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package logging

import (
	"fmt"
	"time"

	"github.com/basecomplextech/baselibrary/collect/slices2"
	"github.com/basecomplextech/baselibrary/panics"
	"github.com/basecomplextech/baselibrary/pools"
	"github.com/basecomplextech/baselibrary/status"
)

type Record struct {
	Time    time.Time `json:"time"`
	Level   Level     `json:"level"`
	Logger  string    `json:"logger"`
	Message string    `json:"message"`
	Fields  []Field   `json:"fields"`
	Stack   []byte    `json:"stack"`
}

// NewRecord returns a new record with the current time.
func NewRecord(logger string, level Level) *Record {
	return &Record{
		Time:   time.Now(),
		Level:  level,
		Logger: logger,
	}
}

// newRecord acquires and returns a new record from the pool.
func newRecord(logger string, level Level) *Record {
	r := recordPool.New()
	r.Time = time.Now()
	r.Level = level
	r.Logger = logger
	return r
}

func (r *Record) WithTime(t time.Time) *Record {
	r.Time = t
	return r
}

func (r *Record) WithLevel(lv Level) *Record {
	r.Level = lv
	return r
}

func (r *Record) WithLogger(logger string) *Record {
	r.Logger = logger
	return r
}

func (r *Record) WithMessage(msg string) *Record {
	r.Message = msg
	return r
}

func (r *Record) WithMessagef(msg string, args ...any) *Record {
	r.Message = fmt.Sprintf(msg, args...)
	return r
}

func (r *Record) WithField(key string, value any) *Record {
	r.Fields = append(r.Fields, Field{
		Key:   key,
		Value: value,
	})
	return r
}

func (r *Record) WithFieldf(key string, format string, a ...any) *Record {
	r.Fields = append(r.Fields, Field{
		Key:   key,
		Value: fmt.Sprintf(format, a...),
	})
	return r
}

func (r *Record) WithFields(keysValues ...any) *Record {
	if len(keysValues) == 0 {
		return r
	}

	fields := NewFields(keysValues...)
	r.Fields = append(r.Fields, fields...)
	return r
}

func (r *Record) WithStack(stack []byte) *Record {
	r.Stack = stack
	return r
}

func (r *Record) WithStatus(st status.Status) *Record {
	r.Fields = append(r.Fields, Field{
		Key:   "status",
		Value: st,
	})

	err := st.Error
	if err == nil {
		return r
	}

	perr, ok := err.(*panics.Error)
	switch {
	case !ok:
		return r
	case perr.Stack == nil:
		return r
	}

	return r.WithStack(perr.Stack)
}

// pool

var recordPool = pools.NewPoolFunc(
	func() *Record {
		return &Record{}
	},
)

func releaseRecord(r *Record) {
	r.reset()
	recordPool.Put(r)
}

func (r *Record) reset() {
	fields := slices2.Truncate(r.Fields)

	*r = Record{}
	r.Fields = fields
}
