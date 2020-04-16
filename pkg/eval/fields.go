package eval

// Fields manages a static set of string fields.
//
// This acts like a method table or prototype.
type Fields map[string]func(receiver Value) Valuable

// Get implements accessing a specific field for a given receiver.
func (f Fields) Get(receiver Value, field Valuable) Valuable {
	if s, ok := field.Value().(strValue); ok {
		if fn, ok := f[string(s)]; ok {
			return fn(receiver)
		}
	}
	return NewError(NewString("no such field"))
}
