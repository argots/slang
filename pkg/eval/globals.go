package eval

// Globals returns a global scope
func Globals() Scope {
	s := NewScope(nil)
	s.Add(NewString("."), operator{".", dot})
	return s
}
