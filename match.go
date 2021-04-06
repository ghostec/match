package match

import (
	"errors"
	"reflect"
)

type Match []Case

type Case struct {
	When when
	Do   interface{}
}

func (m Match) Int(args ...interface{}) int {
	res, err := m.Result(args...)
	if err != nil {
		panic(err)
	}
	val, ok := res.(int)
	if !ok {
		panic("not int")
	}
	return val
}

func (m Match) String(args ...interface{}) string {
	res, err := m.Result(args...)
	if err != nil {
		panic(err)
	}
	val, ok := res.(string)
	if !ok {
		panic("not string")
	}
	return val
}

func (m Match) Slice(args ...interface{}) *SliceType {
	return NewSliceType([]interface{}{})
}

func (m Match) Result(args ...interface{}) (interface{}, error) {
	c := m.match(args)

	if c == nil {
		return nil, errors.New("not matched")
	}

	dotype := reflect.TypeOf(c.Do)
	if dotype.Kind() != reflect.Func {
		return c.Do, nil
	}

	// from now on, c.Do is a func

	dovalue := reflect.ValueOf(c.Do)

	input := make([]reflect.Value, dotype.NumIn())

	rshift := 0
	if dotype.NumIn() > 0 && dotype.In(0) == reflect.TypeOf(Match{}) {
		input[rshift] = reflect.ValueOf(m)
		rshift += 1
	}

	for i := range c.Args {
		input[rshift+i] = reflect.ValueOf(c.Args[i])
	}

	out := dovalue.Call(input)

	switch len(out) {
	case 1:
		return out[0].Interface(), nil
	case 2:
		return out[0].Interface(), out[1].Interface().(error)
	default:
		return nil, errors.New("more than 2 outputs")
	}
}

func (m Match) match(args []interface{}) *CaseWithArgs {
	for i := range m {
		ok, cargs := m[i].When.match(args)
		if ok {
			return &CaseWithArgs{
				Case: &m[i],
				Args: cargs,
			}
		}
	}

	return nil
}

type CaseWithArgs struct {
	*Case
	Args []interface{}
}
