package match

type Match []Case

type Case struct {
	When when
	Do   interface{}
}

func (m Match) Int(args ...interface{}) int {
	return 0
}

func (m Match) String(args ...interface{}) string {
	return ""
}

func (m Match) Slice(args ...interface{}) *SliceType {
	return NewSliceType([]interface{}{})
}
