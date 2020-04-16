package eval

import (
	"github.com/argots/slang/pkg/ast"
)

// type assertion
var _ Value = operator{}

type callable interface {
	Call(x, y ast.Node, s Scope) Valuable
}

type operator struct {
	code string
	fn   func(x, y ast.Node, s Scope) Valuable
}

func (o operator) Type() string {
	return "operator"
}

func (o operator) Code() string {
	return o.code
}

func (o operator) Value() Value {
	return o
}

func (o operator) Get(v Valuable) Valuable {
	return NewError(NewString("unknown field " + v.Value().Code()))
}

func (o operator) Call(x, y ast.Node, s Scope) Valuable {
	return o.fn(x, y, s)
}

func dot(x, y ast.Node, s Scope) Valuable {
	xval := Node(x, s).Value()
	if ident, ok := y.(ast.Ident); ok {
		return xval.Get(NewString(ident.Val))
	}
	return xval.Get(Node(y, s))
}

// Call implements a operator (including (), [], and {})
func Call(v Valuable, x, y ast.Node, s Scope) Valuable {
	val := v.Value()
	if c, ok := val.(callable); ok {
		return c.Call(x, y, s)
	}
	return NewError(NewString("unknown operator " + val.Code()))
}
