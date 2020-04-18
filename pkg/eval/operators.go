package eval

import "github.com/argots/slang/pkg/ast"

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

func arithmetic(op string) func(x, y ast.Node, s Scope) Valuable {
	return func(x, y ast.Node, s Scope) Valuable {
		xval, yval := Node(x, s).Value(), Node(y, s).Value()
		if x == nil {
			xval = NewNumber(0)
		}
		xnum, xok := xval.(numValue)
		ynum, yok := yval.(numValue)
		if !xok || !yok {
			return NewError(NewString("not a number: " + yval.Code()))
		}
		return xnum.Arithmetic(op, ynum)
	}
}

func call(x, y ast.Node, s Scope) Valuable {
	return Call(Node(x, s).Value().Get(NewString("()")), x, y, s)
}

func seq(x, y ast.Node, s Scope) Valuable {
	if x != nil {
		return Call(Node(x, s).Value().Get(NewString("[]")), x, y, s)
	}

	return NewError(NewString("sequences not yet implemented"))
}

func set(x, y ast.Node, s Scope) Valuable {
	if x != nil {
		return Call(Node(x, s).Value().Get(NewString("{}")), x, y, s)
	}
	items := map[string]Valuable{}
	calls := map[string]*Set{}

	var err Valuable
	args := Args{
		NoKey: func(val ast.Node) bool {
			items[NewString("").Code()] = Node(val, s)
			return false
		},
		StringKey: func(key string, val ast.Node) bool {
			items[NewString(key).Code()] = Node(val, s)
			return false
		},
		NodeKey: func(key, val ast.Node) bool {
			items[Node(key, s).Value().Code()] = Node(val, s)
			return false
		},
		ParenKey: func(name string, args, val ast.Node) bool {
			err = defineClosure(calls, name, "()", args, val, s)
			return err != nil
		},
		SetKey: func(name string, args, val ast.Node) bool {
			err = defineClosure(calls, name, "{}", args, val, s)
			return err != nil
		},
		SeqKey: func(name string, args, val ast.Node) bool {
			err = defineClosure(calls, name, "[]", args, val, s)
			return err != nil
		},
	}
	args.Visit(y)
	if err != nil {
		return err
	}
	for k, v := range calls {
		items[NewString(k).Code()] = v
	}
	return &Set{items}
}

func defineClosure(calls map[string]*Set, name, op string, args, val ast.Node, s Scope) Valuable {
	if _, ok := calls[name]; !ok {
		calls[name] = &Set{items: map[string]Valuable{}}
	}
	names := []string{}
	for args != nil {
		arg := args
		args = nil
		if comma, ok := arg.(*ast.Expr); ok && comma.Op == "," {
			arg, args = comma.X, comma.Y
		}
		if ident, ok := arg.(ast.Ident); ok {
			names = append(names, ident.Val)
		} else {
			return NewError(NewString("unexpected function args"))
		}
	}

	opcode := NewString(op).Code()
	calls[name].items[opcode] = NewClosure(op, names, func(args map[Value]Value) Valuable {
		inner := NewScope(s)
		for key, val := range args {
			inner.Add(key, val)
		}
		return Node(val, inner)
	})
	return nil
}

// Call implements a operator (including (), [], and {})
func Call(v Valuable, x, y ast.Node, s Scope) Valuable {
	val := v.Value()
	if c, ok := val.(callable); ok {
		return c.Call(x, y, s)
	}
	return NewError(NewString("unknown operator " + val.Code()))
}
