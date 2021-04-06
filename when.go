package match

import (
	"reflect"
)

func When(args ...interface{}) when {
	all := make([]interface{}, len(args))
	for i := range args {
		all[i] = args[i]
	}
	return all
}

type when []interface{}

func (w when) match(args []interface{}) (ok bool, wargs []interface{}) {
	if len(args) != len(w) {
		return false, nil
	}

	var matchedArgs []Arg

	for i := range w {
		h, isHandle := w[i].(handle)

		wiValue := reflect.ValueOf(w[i])

		isAnyFunc := wiValue.Kind() == reflect.Func && reflect.ValueOf(Any).Pointer() == wiValue.Pointer()
		isEmptyFunc := wiValue.Kind() == reflect.Func && reflect.ValueOf(Empty).Pointer() == wiValue.Pointer()
		isSliceFunc := wiValue.Kind() == reflect.Func && reflect.ValueOf(Slice).Pointer() == wiValue.Pointer()
		isStringFunc := wiValue.Kind() == reflect.Func && reflect.ValueOf(String).Pointer() == wiValue.Pointer()

		switch {
		case isEmptyFunc:
			h = handle{kind: hkEmpty}
		case isSliceFunc:
			h = handle{kind: hkSlice}
		case isStringFunc:
			h = handle{kind: hkString}
		case isAnyFunc, reflect.DeepEqual(w[i], args[i]):
			continue
		case !isHandle:
			return false, nil
		}

		hargs, err := h.Match(args[i])
		if err != nil {
			return false, nil
		}

		matchedArgs = append(matchedArgs, hargs...)
	}

	wargs = make([]interface{}, len(matchedArgs))

	for _, arg := range matchedArgs {
		wargs[arg.N] = arg.Value
	}

	return true, wargs
}
