package match

type hAny struct {
	argn *int
}

func Any(args ...interface{}) (ret hAny) {
	if len(args) == 0 {
		return
	}

	if len(args) > 1 {
		panic("handles: any must have zero or one args")
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

func (h hAny) Match(arg interface{}) (ret []Arg, err error) {
	if h.argn == nil {
		return nil, nil
	}

	return append(ret, Arg{N: *h.argn, Value: arg}), nil
}
