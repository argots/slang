package eval

// Globals returns a global scope
func Globals() Scope {
	s := NewScope(nil)
	s.Add(NewString("."), operator{"sys.operators.dot", dot})
	s.Add(NewString("+"), operator{"sys.operators.add", arithmetic("+")})
	s.Add(NewString("-"), operator{"sys.operators.sub", arithmetic("-")})
	s.Add(NewString("*"), operator{"sys.operators.mul", arithmetic("*")})
	s.Add(NewString("/"), operator{"sys.operators.div", arithmetic("/")})

	s.Add(NewString("{}"), operator{"sys.operators.set", set})
	s.Add(NewString("sys"), sys())
	return s
}

func sys() Value {
	ops := NewString("operators").Code()
	items := map[string]Valuable{
		ops: operators(),
		// number
		// string
		// error
	}
	return &Set{items: items}
}

func operators() Value {
	items := map[string]Valuable{
		NewString("dot").Code(): operator{"sys.operators.dot", dot},
		NewString("add").Code(): operator{"sys.operators.add", arithmetic("+")},
		NewString("sub").Code(): operator{"sys.operators.sub", arithmetic("-")},
		NewString("mul").Code(): operator{"sys.operators.mul", arithmetic("*")},
		NewString("div").Code(): operator{"sys.operators.div", arithmetic("/")},

		NewString("set").Code(): operator{"sys.operators.set", set},
	}
	return &Set{items: items}
}
