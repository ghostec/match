package types

import "reflect"

type Slice struct {
	wrapped interface{}
}

func NewSlice(wrapped interface{}) *Slice {
	sl := &Slice{}

	switch reflect.ValueOf(wrapped).Kind() {
	case reflect.Invalid:
	case reflect.Array, reflect.Slice:
	default:
		if sl.init(wrapped) {
			return sl
		}
	}

	sl.wrapped = wrapped
	return sl
}

func (sl *Slice) Get() interface{} {
	return sl.wrapped
}

func (sl *Slice) Append(el interface{}) *Slice {
	if sl.init(el) {
		return sl
	}

	ret := reflect.Append(reflect.ValueOf(sl.wrapped), reflect.ValueOf(el))
	sl.wrapped = ret.Interface()

	return sl
}

func (sl *Slice) init(val interface{}) bool {
	if sl.wrapped != nil {
		return false
	}

	wrapped := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(val)), 0, 1)
	wrapped = reflect.Append(wrapped, reflect.ValueOf(val))
	sl.wrapped = wrapped.Interface()

	return true
}
