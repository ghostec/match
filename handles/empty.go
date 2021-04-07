package match

import (
	"errors"
	"reflect"

	ty "github.com/ghostec/match/types"
)

type hEmpty struct{}

func Empty() hEmpty {
	return hEmpty{}
}

func (h hEmpty) Match(arg interface{}) (ret []Arg, err error) {
	if sl, ok := arg.(*ty.Slice); ok {
		arg = sl.Get()
	}

	argValue := reflect.ValueOf(arg)

	switch argValue.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		if argValue.Len() == 0 {
			return nil, nil
		}
	}

	return nil, errors.New("not empty")
}
