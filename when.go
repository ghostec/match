package match

import (
	ha "github.com/ghostec/match/handles"
)

type Matcher interface {
	Match(arg interface{}) ([]ha.Arg, error)
}

func When(args ...Matcher) when {
	return args
}

type when []Matcher

func (w when) match(args []interface{}) (ok bool, wargs []interface{}) {
	if len(args) != len(w) {
		return false, nil
	}

	var matchedArgs []ha.Arg

	for i, m := range w {
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
