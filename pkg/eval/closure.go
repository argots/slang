package eval

import "github.com/argots/slang/pkg/ast"

// NewClosure takes an arg definition and makes it callable.
//
// When called, the called args are matched against the definition to
// form a local scope. This is passed to the provided fn for actual
// execution.
func NewClosure(op string, args []string, fn func(args map[Value]Value) Valuable) Valuable {
	c := &closure{args: args, fn: fn}
	c.params = c.seqParams
	if op == "{}" {
		c.params = c.setParams
	}
	return c
}

type closure struct {
	args   []string
	params func(x, y ast.Node, s Scope) (map[Value]Value, Valuable)
	fn     func(args map[Value]Value) Valuable
}

func (c *closure) Type() string {
	return "sys.closure"
}

func (c *closure) Code() string {
	panic("NYI")
}

func (c *closure) Value() Value {
	return c
}

func (c *closure) Get(key Valuable) Valuable {
	return NewError(NewString("no such field " + key.Value().Code()))
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
		names[NewString(name).Code()] = true
	}
	args := Args{
		NoKey: func(val ast.Node) bool {
			err = NewError(NewString("invalid args"))
			return true
		},
		StringKey: func(key string, val ast.Node) bool {
			if code := NewString(key).Code(); !names[code] {
				err = NewError(NewString("invalid arg " + code))
				return true
			}
			params[NewString(key)] = Node(val, s).Value()
			return false
		},
		NodeKey: func(key, val ast.Node) bool {
			keyval := Node(key, s).Value()
			if code := keyval.Code(); !names[code] {
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
