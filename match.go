package match

import (
	"errors"
	"reflect"
)

func When(args ...interface{}) *when {
	w := &when{args: args}
	m := &Match{order: []*when{w}}
	w.match = m
	return w
}

type Match struct {
	order []*when
}

func (m *Match) When(args ...interface{}) *when {
	w := &when{
		match: m,
		args:  args,
	}

	m.order = append(m.order, w)

	return w
}

func (m *Match) Int(args ...interface{}) int {
	res, err := m.result(args...)
	if err != nil {
		panic(err)
	}
	intval, err := res.Int()
	if err != nil {
		panic(err)
	}
	return intval
}

func (m *Match) result(args ...interface{}) (Result, error) {
	for _, when := range m.order {
		pass, err := when.pass(args...)
		if err != nil {
			return Result{nil}, err
		}
		if !pass {
			continue
		}
		return when.result(args...)
	}

	return Result{nil}, errors.New("not matched")
}

type Result struct {
	val interface{}
}

func result(val interface{}) Result {
	return Result{val: val}
}

func (r Result) Int() (int, error) {
	val, ok := r.val.(int)
	if !ok {
		return 0, errors.New("placeholder")
	}
	return val, nil
}

type when struct {
	match *Match
	args  []interface{}
	ret   interface{}
}

func (w *when) pass(args ...interface{}) (bool, error) {
	if len(w.args) != len(args) {
		return false, nil
	}

	for i := range w.args {
		switch {
		case w.args[i] == Any, w.args[i] == args[i]:
			continue
		}
		return false, nil
	}

	return true, nil
}

func (w *when) Return(ret interface{}) *Match {
	w.ret = ret
	return w.match
}

func (w *when) result(args ...interface{}) (Result, error) {
	ret := reflect.TypeOf(w.ret)

	if ret.Kind() != reflect.Func {
		return Result{w.ret}, nil
	}

	out := reflect.ValueOf(w.ret).Call([]reflect.Value{
		reflect.ValueOf(w.match),
		reflect.ValueOf(args[0]),
	})

	if len(out) > 2 {
		return Result{nil}, errors.New("more than two outputs")
	}

	if len(out) == 2 {
		errI := out[1].Interface()
		err, ok := errI.(error)
		if !ok {
			return Result{nil}, errors.New("second output must be an error")
		}
		return Result{out[0].Interface()}, err
	}

	if len(out) == 1 {
		return Result{out[0].Interface()}, nil
	}

	return Result{nil}, errors.New("no outputs")
}

type handle int

const (
	Invalid handle = iota
	Any
)
