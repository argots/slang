// eval implements a simple interpreter for slang
package eval

import (
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
		return strValue(n.Val)
	case ast.Number:
		return NewError(NewString("NYI"))
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
