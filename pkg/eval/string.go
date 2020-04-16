package eval

import (
	"log"
	"strconv"
	"strings"
)

// NewString creates a string value
func NewString(s string) Value {
	ends := `"`
	switch {
	case !strings.Contains(s, `"`):
		ends = `"`
	case !strings.Contains(s, `'`):
		ends = `'`
	case !strings.Contains(s, "`"):
		ends = "`"
	default:
		s = strings.ReplaceAll(s, `"`, `\"`)
	}

	return strValue(ends + s + ends)
}

type strValue string

func (s strValue) Type() string {
	return "sys.string"
}

func (s strValue) Code() string {
	return string(s)
}

func (s strValue) Get(v Valuable) Valuable {
	return strFields().Get(s, v)
}

func (s strValue) Value() Value {
	return s
}

func strFields() Fields {
	return Fields{
		`"length"`: func(receiver Value) Valuable {
			log.Println("calculating length", receiver.Code())
			// TODO: get the raw string
			l := len(string(receiver.(strValue))) - 2
			return strValue(`"` + strconv.Itoa(l) + `"`)
		},
	}
}
