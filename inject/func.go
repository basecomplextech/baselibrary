// Copyright 2024 Ivan Korobkov. All rights reserved.

package inject

import (
	"fmt"
	"reflect"
)

var _ provider = (*funcProvider)(nil)

type funcProvider struct {
	new    reflect.Value
	deps   []reflect.Type
	result reflect.Type
}

func newFuncProvider(fn any) *funcProvider {
	f := reflect.TypeOf(fn)
	if f.Kind() != reflect.Func {
		panic(fmt.Sprintf("constructor must be function, %v", f))
	}

	deps := make([]reflect.Type, f.NumIn())
	for i := range deps {
		dep := f.In(i)
		deps[i] = dep
	}

	if f.NumOut() != 1 {
		panic(fmt.Sprintf("constructor must return one value, %v", f))
	}

	result := f.Out(0)
	return &funcProvider{
		new:    reflect.ValueOf(fn),
		deps:   deps,
		result: result,
	}
}

func (p *funcProvider) typ() reflect.Type {
	return p.result
}

func (p *funcProvider) init(x *context) reflect.Value {
	args := make([]reflect.Value, len(p.deps))
	for i, dep := range p.deps {
		args[i] = x.get(dep)
	}
	return p.new.Call(args)[0]
}
