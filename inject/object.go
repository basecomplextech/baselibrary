package inject

import "reflect"

var _ provider = (*objectProvider)(nil)

type objectProvider struct {
	object reflect.Value
	result reflect.Type
}

func newObjectProvider(obj any) *objectProvider {
	return &objectProvider{
		object: reflect.ValueOf(obj),
		result: reflect.TypeOf(obj),
	}
}

func (p *objectProvider) typ() reflect.Type {
	return p.result
}

func (p *objectProvider) init(x *context) reflect.Value {
	return p.object
}
