package eval

import "github.com/argots/slang/pkg/cast"

var _ Value = &Set{}

type setItem struct {
	Key, Value Valuable
}

// Set implememnts a generic set type
type Set struct {
	// items has toString(Key) as the actual key
	items map[string]setItem
}

func (s *Set) Add(key, value Valuable) {
	s.items[toString(key)] = setItem{key, value}
}

// Type returns the type of the set
func (s *Set) Type() string {
	return "sys.operators.set{}"
}

// Code returns the code for a set
func (s *Set) Code() Code {
	args := []interface{}{}
	for _, item := range s.items {
		args = append(args, cast.Pair(item.Key.Value().Code(), item.Value.Value().Code()))
	}
	return Code{cast.Set(nil, args...).Node}
}

// Value returns the set itself
func (s *Set) Value() Value {
	return s
}

// Get returns the value for a key
func (s *Set) Get(key Valuable) Valuable {
	code := toString(key)
	if v, ok := s.items[code]; ok {
		return v.Value
	}
	return NewError(NewString("not found: " + code))
}
