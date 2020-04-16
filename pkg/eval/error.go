package eval

// type assertion
var _ Value = &errorValue{}

// NewError creates an error value.
func NewError(v Valuable) Value {
	return &errorValue{v}
}

type errorValue struct {
	v Valuable
}

func (e *errorValue) Type() string {
	return "sys.error{" + e.v.Value().Type() + "}"
}

func (e *errorValue) Code() string {
	return `sys.error{` + e.v.Value().Code() + "}"
}

func (e *errorValue) Value() Value {
	return e
}

func (e *errorValue) Get(v Valuable) Valuable {
	return e
}
