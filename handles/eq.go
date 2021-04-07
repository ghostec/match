package match

import (
	"errors"
	"reflect"
)

type hEq struct {
	wrapped interface{}
}

func Eq(wrapped interface{}) (ret hEq) {
	ret.wrapped = wrapped
	return
}

func (h hEq) Match(arg interface{}) (ret []Arg, err error) {
	if !reflect.DeepEqual(h.wrapped, arg) {
		return nil, errors.New("handles: no match")
	}

	return nil, nil
}
