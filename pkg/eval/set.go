package eval

var _ Value = &Set{}

// Set implememnts a generic set type
type Set struct {
	items map[string]Valuable
}

// Type returns the type of the set
func (s *Set) Type() string {
	return "sys.operators.set{}"
}

// Code returns the code for a set
func (s *Set) Code() string {
	result := "{"
	first := true
	for k, v := range s.items {
		if !first {
			result += ", "
		}
		first = false
		result += k + ": " + v.Value().Code()
	}
	return result + "}"
}

// Value returns the set itself
func (s *Set) Value() Value {
	return s
}

// Get returns the value for a key
func (s *Set) Get(key Valuable) Valuable {
	code := key.Value().Code()
	if v, ok := s.items[code]; ok {
		return v
	}
	return NewError(NewString("not found: " + code))
}
