package eval

// Globals returns a global scope
func Globals() Scope {
	s := NewScope(nil)
	s.Add(NewString("."), operator{"sys.operators.dot", dot})
	s.Add(NewString("sys"), sys())
	return s
}

func sys() Value {
	ops := NewString("operators").Code()
	items := map[string]Valuable{ops: operators()}
	return &Set{items: items}
}

func operators() Value {
	items := map[string]Valuable{
		NewString("dot").Code(): operator{"sys.operators.dot", dot},
	}
	return &Set{items: items}
}
