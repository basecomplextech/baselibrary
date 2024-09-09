// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package inject

import "reflect"

type provider interface {
	typ() reflect.Type
	init(x *context) reflect.Value
}

func newProvider(v any) provider {
	typ := reflect.TypeOf(v)
	kind := typ.Kind()

	if kind == reflect.Func {
		return newFuncProvider(v)
	}
	return newObjectProvider(v)
}
