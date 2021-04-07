package match

import (
	"errors"
)

type hString struct {
	argn *int
}

func String(args ...interface{}) (ret hString) {
	if len(args) == 0 {
		return
	}

	switch len(args) {
	case 1:
		if argn, ok := args[0].(int); ok {
			ret.argn = &argn
			return
		}

		panic("handles: any's argument must be an index")
	}

	return
}

func (h hString) Match(arg interface{}) (ret []Arg, err error) {
	if _, ok := arg.(string); !ok {
		return nil, errors.New("handles: arg isn't string")
	}

	if h.argn == nil {
		return
	}

	ret = append(ret, Arg{N: *h.argn, Value: arg})
	return
}
