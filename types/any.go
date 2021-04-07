package types

type Any struct {
	wrapped interface{}
}

func NewAny(wrapped interface{}) *Any {
	return &Any{wrapped}
}

func (any *Any) Get() interface{} {
	return any.wrapped
}
