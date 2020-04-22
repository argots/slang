package eval

import (
	"github.com/argots/slang/pkg/cast"
)

var _ Value = strValue("")

// NewString creates a string value
func NewString(s string) Value {
	return strValue(s)
}

type strValue string

func (s strValue) Type() string {
	return "sys.string"
}

func (s strValue) Code() Code {
	return Code{cast.Quote(string(s)).Node}
}

func (s strValue) Get(v Valuable) Valuable {
	return strFields().Get(s, v)
}

func (s strValue) Value() Value {
	return s
}

func strFields() Fields {
	return Fields{
		"length": func(receiver Value) Valuable {
			l := len(string(receiver.(strValue)))
			return NewNumber(float64(l))
		},
	}
}
