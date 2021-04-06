package match

import (
	"errors"
	"reflect"
)

func Match(args ...interface{}) *match {
	return &match{args: args}
}

type match struct {
	args  []interface{}
	order []*when
}

func (m *match) When(args ...interface{}) *when {
	w := &when{
		match: m,
		args:  args,
	}

	m.order = append(m.order, w)

	return w
}

func (m *match) Int() int {
	res, err := m.result()
	if err != nil {
		panic(err)
	}
	intval, err := res.Int()
	if err != nil {
		panic(err)
	}
	return intval
}

func (m *match) result() (Result, error) {
	for _, when := range m.order {
		pass, err := when.pass()
		if err != nil {
			return Result{nil}, err
		}
		if !pass {
			continue
		}
		return when.result()
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
	match *match
	args  []interface{}
	ret   interface{}
}

func (w *when) pass() (bool, error) {
	if len(w.args) != len(w.match.args) {
		return false, nil
	}

	for i := range w.args {
		switch {
		case w.args[i] == Any, w.args[i] == w.match.args[i]:
			continue
		}
		return false, nil
	}

	return true, nil
}

func (w *when) Do(f interface{}) *match {
	w.ret = f
	return w.match
}

func (w *when) Return(ret interface{}) *match {
	w.ret = ret
	return w.match
}

func (w *when) result() (Result, error) {
	ret := reflect.TypeOf(w.ret)

	if ret.Kind() != reflect.Func {
		return Result{w.ret}, nil
	}

	out := reflect.ValueOf(w.ret).Call([]reflect.Value{reflect.ValueOf(w.match.args[0])})

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
