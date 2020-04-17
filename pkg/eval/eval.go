// eval implements a simple interpreter for slang
package eval

import (
	"strconv"

	"github.com/argots/slang/pkg/ast"
)

// Value represents a value.  Most cases, Valuable is a better
// interface to use.
type Value interface {
	Code() string
	Type() string
	Get(v Valuable) Valuable
	Valuable
}

// Valuable represents something that can provide a value.
type Valuable interface {
	Value() Value
}

// Node evaluates a node
func Node(n ast.Node, s Scope) Valuable {
	switch n := n.(type) {
	case ast.Quote:
		return strValue(decodeString(n.Val))
	case ast.Number:
		f, err := strconv.ParseFloat(n.Val, 64)
		if err != nil {
			return NewError(NewString(err.Error()))
		}
		return NewNumber(f)
	case ast.Ident:
		return s.Get(NewString(n.Val)).Value()
	case *ast.Expr:
		return Call(s.Get(NewString(n.Op)), n.X, n.Y, s)
	case *ast.Paren:
		if n.X == nil {
			return Node(n.Y, s)
		}
		return Call(s.Get(NewString("()")), n.X, n.Y, s)
	case *ast.Seq:
		return Call(s.Get(NewString("[]")), n.X, n.Y, s)
	case *ast.Set:
		return Call(s.Get(NewString("{}")), n.X, n.Y, s)
	}

	return NewError(NewString("nil"))
}

func decodeString(s string) string {
	rs := []rune{}
	skip := false
	for _, r := range s[1 : len(s)-1] {
		if skip || r != '\\' {
			rs = append(rs, r)
		}
		skip = !skip && r == '\\'
	}
	return string(rs)
}
