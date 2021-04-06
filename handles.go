package match

import (
	"errors"
	"reflect"
)

type handle struct {
	kind     handleKind
	argn     *int
	size     int
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
	ret.kind = hkSlice

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
			argValue := reflect.ValueOf(args[i])
			isArgFunc := argValue.Kind() == reflect.Func

			isHeadFunc := isArgFunc && argValue.Pointer() == reflect.ValueOf(Head).Pointer()
			isTailFunc := isArgFunc && argValue.Pointer() == reflect.ValueOf(Tail).Pointer()
			isSliceFunc := isArgFunc && argValue.Pointer() == reflect.ValueOf(Slice).Pointer()

			h, isHandle := args[i].(handle)

			isHead := isHeadFunc || (isHandle && h.kind == hkHead)
			isTail := isTailFunc || (isHandle && h.kind == hkTail)
			isSlice := isSliceFunc || (isHandle && h.kind == hkSlice)

			switch {
			case isHead, isTail, isSlice:
				ret.children = append(ret.children, args[i])
			default:
				panic("not head, tail or slice")
			}
		}
	}

	return
}

func Head(argn, size int) (ret handle) {
	ret.kind = hkHead
	ret.size = size

	if argn >= 0 {
		ret.argn = &argn
	}

	return
}

func Tail(argn, size int) (ret handle) {
	ret.kind = hkTail
	ret.size = size

	if argn > 0 {
		ret.argn = &argn
	}

	return
}

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

func (h handle) Match(arg interface{}) (ret []Arg, err error) {
	switch h.kind {
	case hkEmpty:
		argValue := reflect.ValueOf(arg)
		switch argValue.Kind() {
		case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
			if argValue.Len() == 0 {
				return nil, nil
			}
		}

		return nil, errors.New("not empty")

	case hkAny:
		if h.argn != nil {
			ret = append(ret, Arg{N: *h.argn, Value: arg})
		}

	case hkString:
		if _, ok := arg.(string); !ok {
			return nil, errors.New("arg isn't string")
		}

		if h.argn != nil {
			ret = append(ret, Arg{N: *h.argn, Value: arg})
		}

	case hkHead, hkTail:
		argValue := reflect.ValueOf(arg)

		if argValue.Kind() != reflect.Slice {
			return nil, errors.New("arg isn't slice")
		}

		var value interface{}

		if h.size == 1 {
			switch h.kind {
			case hkHead:
				value = argValue.Index(0)
			case hkTail:
				value = argValue.Index(argValue.Len() - 1)
			}
		} else {
			switch h.kind {
			case hkHead:
				ret = append(ret, Arg{N: *h.argn, Value: argValue.Slice(0, h.size).Interface()})
			case hkTail:
				l := argValue.Len()
				ret = append(ret, Arg{N: *h.argn, Value: argValue.Slice(l-h.size, l).Interface()})
			}
		}

		ret = append(ret, Arg{N: *h.argn, Value: value})

	case hkSlice:
		argValue := reflect.ValueOf(arg)

		if argValue.Kind() != reflect.Slice {
			return nil, errors.New("arg isn't slice")
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
			hh, ok := h.children[i].(handle)
			if !ok {
				continue
			}

			switch hh.kind {
			case hkSlice:
				start, end := -1, -1

				if i == 0 {
					next := h.children[1]
					nextValue := reflect.ValueOf(next)
					isNextFunc := nextValue.Kind() == reflect.Func

					isTailFunc := isNextFunc && nextValue.Pointer() == reflect.ValueOf(Tail).Pointer()
					hh, isHandle := next.(handle)
					isTail := isTailFunc || (isHandle && hh.kind == hkTail)

					if !isTail {
						return nil, errors.New("slice neighbors must be head or tail")
					}

					size := 1
					if !isTailFunc {
						size = hh.size
					}

					start, end = 0, argValue.Len()-size
				} else {
					prev := h.children[0]
					prevValue := reflect.ValueOf(prev)
					isPrevFunc := prevValue.Kind() == reflect.Func

					isHeadFunc := isPrevFunc && prevValue.Pointer() == reflect.ValueOf(Head).Pointer()
					hh, isHandle := prev.(handle)
					isHead := isHeadFunc || (isHandle && hh.kind == hkHead)

					if !isHead {
						return nil, errors.New("slice neighbors must be head or tail")
					}

					size := 1
					if !isHeadFunc {
						size = hh.size
					}

					start, end = size, argValue.Len()
				}

				ret = append(ret, Arg{
					N:     *hh.argn,
					Value: argValue.Slice(start, end).Interface(),
				})

			case hkHead, hkTail:
				hargs, err := hh.Match(arg)
				if err != nil {
					return nil, err
				}

				ret = append(ret, hargs...)
			}
		}
	}

	return
}

type Arg struct {
	N     int
	Value interface{}
}
