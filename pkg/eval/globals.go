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
	s.Add(NewString("()"), operator{"sys.operators.call", call})
	s.Add(NewString("[]"), operator{"sys.operators.seq", seq})
	s.Add(NewString("sys"), sys())
	return s
}

func sys() Value {
	result := &Set{items: map[string]setItem{}}
	result.Add(NewString("operators"), operators())
	return result
}

func operators() Value {
	ops := &Set{items: map[string]setItem{}}
	ops.Add(NewString("dot"), operator{"sys.operators.dot", dot})
	ops.Add(NewString("add"), operator{"sys.operators.add", arithmetic("+")})
	ops.Add(NewString("sub"), operator{"sys.operators.sub", arithmetic("-")})
	ops.Add(NewString("mul"), operator{"sys.operators.mul", arithmetic("*")})
	ops.Add(NewString("div"), operator{"sys.operators.div", arithmetic("/")})

	ops.Add(NewString("set"), operator{"sys.operators.set", set})
	ops.Add(NewString("call"), operator{"sys.operators.call", call})
	ops.Add(NewString("seq"), operator{"sys.operators.seq", seq})
	return ops
}
