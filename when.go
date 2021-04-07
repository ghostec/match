package match

import (
	"reflect"

	ha "github.com/ghostec/match/handles"
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

	var matchedArgs []ha.Arg

	for i := range w {
		if reflect.DeepEqual(w[i], args[i]) {
			continue
		}

		m, isMatcher := w[i].(interface {
			Match(interface{}) ([]ha.Arg, error)
		})

		if !isMatcher {
			return false, nil
		}

		wiargs, err := m.Match(args[i])
		if err != nil {
			return false, nil
		}

		matchedArgs = append(matchedArgs, wiargs...)
	}

	wargs = make([]interface{}, len(matchedArgs))

	for _, arg := range matchedArgs {
		wargs[arg.N] = arg.Value
	}

	return true, wargs
}
