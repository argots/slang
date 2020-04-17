package eval

import (
	"strings"
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

func (s strValue) Code() string {
	sx := string(s)
	ends := `"`
	switch {
	case !strings.Contains(sx, `"`):
		ends = `"`
	case !strings.Contains(sx, `'`):
		ends = `'`
	case !strings.Contains(sx, "`"):
		ends = "`"
	default:
		sx = strings.ReplaceAll(sx, `"`, `\"`)
	}

	return ends + sx + ends
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
