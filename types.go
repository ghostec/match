package match

import "reflect"

type SliceType struct {
	wrapped interface{}
}

func (sl *SliceType) Get() interface{} {
	return sl.wrapped
}

func NewSliceType(wrapped interface{}) *SliceType {
	if reflect.TypeOf(wrapped).Kind() != reflect.Slice {
		panic("not slice")
	}

	return &SliceType{wrapped}
}

func (sl *SliceType) Append(el interface{}) *SliceType {
	ret := reflect.Append(reflect.ValueOf(sl.wrapped), reflect.ValueOf(el))
	sl.wrapped = ret.Interface()
	return sl
}
