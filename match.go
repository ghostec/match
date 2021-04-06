package match

import "reflect"

type Match []Case

type Case struct {
	When when
	Do   interface{}
}

func (m Match) Int(args ...interface{}) int {
	return 0
}

func (m Match) String(args ...interface{}) string {
	return ""
}

func (m Match) Slice(args ...interface{}) *SliceType {
	return NewSliceType([]interface{}{})
}

type handle struct {
	kind     handleKind
	argn     *int
	children []interface{}
}

type handleFunc func(args ...interface{}) handle

type handleKind int

const (
	hkInvalid handleKind = iota
	hkEmpty
	hkString
	hkSlice
	hkHead
	hkTail
	hkAny
)

func Any(args ...interface{}) (ret handle) {
	switch len(args) {
	case 0:
		ret.kind = hkAny
		return
	case 1:
		if argn, ok := args[0].(int); ok {
			ret.argn = &argn
			ret.kind = hkAny
		}
	}

	return
}

func Empty() handle {
	return handle{kind: hkEmpty}
}

func String(args ...interface{}) (ret handle) {
	switch len(args) {
	case 0:
		ret.kind = hkString
		return
	case 1:
		if argn, ok := args[0].(int); ok {
			ret.argn = &argn
			ret.kind = hkString
		}
	}

	return
}

func Slice(args ...interface{}) (ret handle) {
	switch len(args) {
	case 0:
		ret.kind = hkSlice
		return
	case 1:
		if argn, ok := args[0].(int); ok {
			ret.argn = &argn
			ret.kind = hkSlice
		}
	}

	return
}

func Head(argn, size int) (ret handle) {
	ret.kind = hkHead

	if argn > 0 {
		ret.argn = &argn
	}

	return
}

func Tail(argn, size int) (ret handle) {
	ret.kind = hkTail

	if argn > 0 {
		ret.argn = &argn
	}

	return
}

func When(args ...interface{}) when {
	all := make([]interface{}, len(args))
	for i := range args {
		all[i] = args[i]
	}
	return all
}

type when []interface{}

type SliceType struct {
	wrapped interface{}
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
