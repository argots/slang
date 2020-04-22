package eval

import (
	"github.com/argots/slang/pkg/ast"
	"github.com/argots/slang/pkg/cast"
)

// NewClosure takes an arg definition and makes it callable.
//
// When called, the called args are matched against the definition to
// form a local scope. This is passed to the provided fn for actual
// execution.
func NewClosure(op string, args []string, fn func(args map[Value]Value) Valuable) Valuable {
	c := &closure{op: op, args: args, fn: fn}
	c.params = c.seqParams
	if op == "{}" {
		c.params = c.setParams
	}
	return c
}

type closure struct {
	op     string
	args   []string
	params func(x, y ast.Node, s Scope) (map[Value]Value, Valuable)
	fn     func(args map[Value]Value) Valuable
	fnCode Code
}

func (c *closure) Type() string {
	return "sys.closure"
}

func (c *closure) Code() Code {
	code := cast.ToNode("closure")
	args := []interface{}{}
	for _, arg := range c.args {
		args = append(args, cast.ToNode(arg))
	}
	var key cast.Node
	switch c.op {
	case "()":
		key = code.Call(args...)
	case "{}":
		key = code.Set(args...)
	case "[]":
		key = code.Seq(args...)
	}
	return Code{cast.Set(nil, cast.Pair(key, c.fnCode)).Dot("closure").Node}
}

func (c *closure) Value() Value {
	return c
}

func (c *closure) Get(key Valuable) Valuable {
	return NewError(NewString("no such field " + toString(key)))
}

func (c *closure) seqParams(x, y ast.Node, s Scope) (map[Value]Value, Valuable) {
	params := map[Value]Value{}
	var err Valuable

	args := Args{
		NoKey: func(val ast.Node) bool {
			if len(c.args) == 0 {
				err = NewError(NewString("invalid args"))
				return true
			}
			params[NewString(c.args[0])] = Node(val, s).Value()
			c.args = c.args[1:]
			return false
		},
		StringKey: func(_ string, val ast.Node) bool {
			err = NewError(NewString("invalid args"))
			return true
		},
		NodeKey: func(key, val ast.Node) bool {
			err = NewError(NewString("invalid args"))
			return true
		},
		ParenKey: func(name string, args, val ast.Node) bool {
			err = NewError(NewString("invalid args"))
			return true
		},
		SetKey: func(name string, args, val ast.Node) bool {
			err = NewError(NewString("invalid args"))
			return true
		},
		SeqKey: func(name string, args, val ast.Node) bool {
			err = NewError(NewString("invalid args"))
			return true
		},
	}
	args.Visit(y)
	return params, err
}

func (c *closure) setParams(x, y ast.Node, s Scope) (map[Value]Value, Valuable) {
	params := map[Value]Value{}
	var err Valuable

	names := map[string]bool{}
	for _, name := range c.args {
		names[toString(NewString(name))] = true
	}
	args := Args{
		NoKey: func(val ast.Node) bool {
			err = NewError(NewString("invalid args"))
			return true
		},
		StringKey: func(key string, val ast.Node) bool {
			if code := toString(NewString(key)); !names[code] {
				err = NewError(NewString("invalid arg " + code))
				return true
			}
			params[NewString(key)] = Node(val, s).Value()
			return false
		},
		NodeKey: func(key, val ast.Node) bool {
			keyval := Node(key, s).Value()
			if code := toString(keyval); !names[code] {
				err = NewError(NewString("invalid arg " + code))
				return true
			}
			params[keyval] = Node(val, s).Value()
			return false
		},
		ParenKey: func(name string, args, val ast.Node) bool {
			err = NewError(NewString("fn args not yet impleemnted: " + name))
			return true
		},
		SetKey: func(name string, args, val ast.Node) bool {
			err = NewError(NewString("fn args not yet impleemnted: " + name))
			return true
		},
		SeqKey: func(name string, args, val ast.Node) bool {
			err = NewError(NewString("fn args not yet impleemnted: " + name))
			return true
		},
	}
	args.Visit(y)
	return params, err
}

func (c *closure) Call(x, y ast.Node, s Scope) Valuable {
	params, err := c.params(x, y, s)
	if err != nil {
		return err
	}
	return c.fn(params)
}
