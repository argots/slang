package eval

// Scope defines a current scope.
type Scope interface {
	Get(key Value) Valuable
	Add(key Value, value Valuable)
}

// NewScope creates a new scope, possibly from another scope.
func NewScope(parent Scope) Scope {
	return &scope{parent: parent}
}

type scopeItem struct {
	key   Value
	value Valuable
}

type scope struct {
	parent Scope
	items  []scopeItem
}

func (s *scope) Get(key Value) Valuable {
	for _, item := range s.items {
		if s.equals(item.key, key) {
			return item.value
		}
	}
	if s.parent != nil {
		return s.parent.Get(key)
	}
	return NewError(NewString("no such key"))
}

func (s *scope) Add(key Value, value Valuable) {
	s.items = append(s.items, scopeItem{key, value})
}

func (s *scope) equals(key1, key2 Value) bool {
	// TODO: use a more efficient equals implementation
	return key1.Type() == key2.Type() && key1.Code() == key2.Code()
}
