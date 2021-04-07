package match

import (
	"errors"
	"reflect"

	ty "github.com/ghostec/match/types"
)

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

type hSlice struct {
	argn     *int
	children []interface{}
}

func Slice(args ...interface{}) (ret hSlice) {
	if len(args) == 0 {
		return
	}

	switch len(args) {
	case 1:
		argn, ok := args[0].(int)
		if !ok {
			panic("only argument must be argn")
		}

		ret.argn = &argn
	default:
		argn, ok := args[0].(int)
		if ok {
			ret.argn = &argn
		}

		for i := range args {
			passed := false

			if _, ok := args[i].(hHead); ok {
				passed = true
			} else if _, ok := args[i].(hTail); ok {
				passed = true
			} else if _, ok := args[i].(hSlice); ok {
				passed = true
			}

			if !passed {
				panic("not head, tail or slice")
			}

			ret.children = append(ret.children, args[i])
		}
	}

	return
}

// TODO: write this as a ma.Match expression
// maybe Arg needs to go on ma instead of ha
func (h hSlice) Match(arg interface{}) (ret []Arg, err error) {
	if sl, ok := arg.(*ty.Slice); ok {
		arg = sl.Get()
	}

	argValue := reflect.ValueOf(arg)

	switch argValue.Kind() {
	case reflect.Array, reflect.Slice:
	default:
		return nil, errors.New("handles: arg isn't slice")
	}

	if h.argn == nil && len(h.children) == 0 {
		return
	}

	if len(h.children) == 1 {
		return nil, errors.New("slice can't be single child of slice")
	}

	if h.argn != nil {
		ret = append(ret, Arg{
			N:     *h.argn,
			Value: arg,
		})
	}

	for i := range h.children {
		switch child := h.children[i].(type) {
		case hSlice:
			start, end := -1, -1
			l := argValue.Len()

			if l > 0 && i == 0 {
				next := h.children[1]
				tl, ok := next.(hTail)
				if !ok {
					return nil, errors.New("slice neighbors must be head or tail")
				}

				start, end = 0, argValue.Len()-tl.size
			} else if l > 0 {
				prev := h.children[0]
				hd, ok := prev.(hHead)
				if !ok {
					return nil, errors.New("slice neighbors must be head or tail")
				}

				start, end = hd.size, argValue.Len()
			}

			var value interface{}
			if start != -1 && end != -1 {
				value = argValue.Slice(start, end).Interface()
			}

			if child.argn != nil {
				ret = append(ret, Arg{
					N:     *child.argn,
					Value: value,
				})
			}

		case hHead, hTail:
			matcher, _ := child.(interface {
				Match(interface{}) ([]Arg, error)
			})

			hargs, err := matcher.Match(arg)
			if err != nil {
				return nil, err
			}

			ret = append(ret, hargs...)

		default:
			return nil, errors.New("handles: invalid child for slice")
		}
	}

	return
}

type hHead struct {
	argn *int
	size int
}

type hTail struct {
	argn *int
	size int
}

func Head(argn, size int) (ret hHead) {
	ret.size = size

	if argn >= 0 {
		ret.argn = &argn
	}

	return
}

func Tail(argn, size int) (ret hTail) {
	ret.size = size

	if argn > 0 {
		ret.argn = &argn
	}

	return
}

func (h hHead) Match(arg interface{}) (ret []Arg, err error) {
	argValue := reflect.ValueOf(arg)

	switch argValue.Kind() {
	case reflect.Array, reflect.Slice:
	default:
		return nil, errors.New("handles: arg isn't slice")
	}

	var value interface{}

	if argValue.Len() == 0 {
		return nil, errors.New("handles: empty list can't have a head")
	}

	if h.size == 1 {
		value = argValue.Index(0).Interface()
	} else {
		value = argValue.Slice(0, h.size).Interface()
	}

	ret = append(ret, Arg{N: *h.argn, Value: value})

	return
}

func (h hTail) Match(arg interface{}) (ret []Arg, err error) {
	argValue := reflect.ValueOf(arg)

	switch argValue.Kind() {
	case reflect.Array, reflect.Slice:
	default:
		return nil, errors.New("handles: arg isn't slice")
	}

	if argValue.Len() == 0 {
		return nil, errors.New("handles: empty list can't have a tail")
	}

	var value interface{}

	if h.size == 1 {
		value = argValue.Index(argValue.Len() - 1).Interface()
	} else {
		l := argValue.Len()
		value = argValue.Slice(l-h.size, l).Interface()
	}

	ret = append(ret, Arg{N: *h.argn, Value: value})

	return
}

type Arg struct {
	N     int
	Value interface{}
}
