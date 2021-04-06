package match

func When(args ...interface{}) when {
	all := make([]interface{}, len(args))
	for i := range args {
		all[i] = args[i]
	}
	return all
}

type when []interface{}

func (w when) match(args []interface) bool {
	if len(args) != len(w) {
		return false
	}

	for i := range args {

	}
}
