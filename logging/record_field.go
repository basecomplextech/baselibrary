// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package logging

import "fmt"

type Field struct {
	Key   string
	Value any
}

// NewField returns a new field.
func NewField(key string, value any) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}

// NewFields builds fields from key/value pairs.
func NewFields(keyValues ...any) []Field {
	if len(keyValues) == 0 {
		return nil
	}

	fields := make([]Field, 0, len(keyValues)/2)

	for i := 0; i < len(keyValues); i += 2 {
		k := keyValues[i]
		key, ok := k.(string)
		if !ok {
			key = fmt.Sprintf("%v", k)
		}

		var value any
		if i < len(keyValues)-1 {
			value = keyValues[i+1]
		}

		field := NewField(key, value)
		fields = append(fields, field)
	}

	return fields
}
