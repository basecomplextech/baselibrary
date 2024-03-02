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
