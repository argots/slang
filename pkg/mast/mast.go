// Package mast implements utilities to match ASTs.
package mast

import (
	"strconv"

	"github.com/argots/slang/pkg/ast"
)

// Matcher is a general purpose matcher function to see if a node
// matches some condition.
//
// It exports most top-level functions as methods to support ease of
// use via chaining.  For example, to test if an expression is `x +
// y`, One can write `Op("+").X(Ident("x")).Y(Ident("y"))`
type Matcher func(n *ast.Node, tx *Tx) bool

// Or returns a matcher which matches the current or any of the
// provided matches.
func (m Matcher) Or(others ...Matcher) Matcher {
	return Or(m, others...)
}

// Contains matches m against X or Y recursively
func (m Matcher) Contains(other Matcher) Matcher {
	return m.And(Contains(other))
}

// MaybeParen matches m or (m)
func (m Matcher) MaybeParen() Matcher {
	return MaybeParen(m)
}

// And returns a matcher which ensures all matchers succeed.
func (m Matcher) And(others ...Matcher) Matcher {
	return And(m, others...)
}

// Call matches m(args)
func (m Matcher) Call(args ...Matcher) Matcher {
	return Paren().X(m).Y(ArgListSeq(args...))
}

// Seq matches m[args]
func (m Matcher) Seq(args ...Matcher) Matcher {
	return Seq().X(m).Y(ArgListSeq(args...))
}

// Set matches m{args}
func (m Matcher) Set(args ...Matcher) Matcher {
	return Set().X(m).Y(ArgListSeq(args...))
}

// Dot matches m.f.
func (m Matcher) Dot(f Matcher) Matcher {
	return Op(".").X(m).Y(f)
}

// X matches a node with the X part of it.
func (m Matcher) X(other Matcher) Matcher {
	return m.And(X(other))
}

// Y matches a node with the Y part of it.
func (m Matcher) Y(other Matcher) Matcher {
	return m.And(Y(other))
}

// HasItem matches if a set, seq or arglist has a specific item.
func (m Matcher) HasItem(item Matcher) Matcher {
	return m.And(HasItem(item))
}

// HasKeyValue matches if a set, seq or arglist has a specific key value.
func (m Matcher) HasKeyValue(key, value Matcher) Matcher {
	return m.And(HasKeyValue(key, value))
}

// Capture captures the node into the pointer to it if the condition
// is met.
func (m Matcher) Capture(n *ast.Node) Matcher {
	return m.And(Capture(n))
}

// CaptureVal captures the node value if the condition is met
func (m Matcher) CaptureVal(s *string) Matcher {
	return m.And(CaptureVal(s))
}

// Replace replaces the node with another.
func (m Matcher) Replace(fn func() ast.Node) Matcher {
	return m.And(Replace(fn))
}

// Any matches any node.
func Any() Matcher {
	return func(n *ast.Node, tx *Tx) bool {
		return true
	}
}

// Or returns a matcher which matches the current or any of the
// provided matches.
func Or(m Matcher, ms ...Matcher) Matcher {
	return func(np *ast.Node, tx *Tx) bool {
		nested := tx.Begin()
		defer nested.End()

		if m(np, nested) {
			return true
		}
		nested.Cancel()
		for _, other := range ms {
			if other(np, nested) {
				return true
			}
			nested.Cancel()
		}
		return false
	}
}

// Contains matches m against X or Y recursively.
func Contains(m Matcher) Matcher {
	recurse := Matcher(func(np *ast.Node, tx *Tx) bool {
		return Contains(m)(np, tx)
	})
	return Or(X(m.Or(recurse)), Y(m.Or(recurse)))
}

// And returns a matcher which ensures all matchers succeed.
func And(m Matcher, ms ...Matcher) Matcher {
	return func(np *ast.Node, tx *Tx) bool {
		n := *np
		if !m(np, tx) {
			*np = n
			return false
		}
		for _, other := range ms {
			if !other(np, tx) {
				*np = n
				return false
			}
		}
		return true
	}
}

// Op returns a matcher which checks if the node has one of the
// specified ops.  Node that () {} or [] are not included.
func Op(ops ...string) Matcher {
	return func(np *ast.Node, tx *Tx) bool {
		var found bool

		x, ok := (*np).(*ast.Expr)
		for kk := 0; ok && !found && kk < len(ops); kk++ {
			found = x.Op == ops[kk]
			switch ops[kk] {
			case "()":
				found = found || Paren()(np, tx)
			case "[]":
				found = found || Seq()(np, tx)
			case "{}":
				found = found || Set()(np, tx)
			}
		}
		return found
	}
}

// X matches a node with the X part of it.
func X(m Matcher) Matcher {
	return func(np *ast.Node, tx *Tx) bool {
		n := *np
		matched := false
		switch n := n.(type) {
		case *ast.Expr:
			copy := *n
			n = &copy
			matched = m(&n.X, tx)
		case *ast.Paren:
			copy := *n
			n = &copy
			matched = m(&n.X, tx)
		case *ast.Set:
			copy := *n
			n = &copy
			matched = m(&n.X, tx)
		case *ast.Seq:
			copy := *n
			n = &copy
			matched = m(&n.X, tx)
		}
		if matched {
			*np = n
		}

		return matched
	}
}

// Y matches a node with the Y part of it.
func Y(m Matcher) Matcher {
	return func(np *ast.Node, tx *Tx) bool {
		n := *np
		matched := false
		switch n := n.(type) {
		case *ast.Expr:
			copy := *n
			n = &copy
			matched = m(&n.Y, tx)
		case *ast.Paren:
			copy := *n
			n = &copy
			matched = m(&n.Y, tx)
		case *ast.Set:
			copy := *n
			n = &copy
			matched = m(&n.Y, tx)
		case *ast.Seq:
			copy := *n
			n = &copy
			matched = m(&n.Y, tx)
		}
		if matched {
			*np = n
		}

		return matched
	}
}

// ArgListSeq matches all the matchers against the provided args.
func ArgListSeq(ms ...Matcher) Matcher {
	l := len(ms)
	switch l {
	case 0:
		return Nil()
	case 1:
		return ms[0]
	}
	return Op(",").Y(ms[l-1]).X(ArgListSeq(ms[:l-1]...))
}

// MaybeParen matches either m or (m).
func MaybeParen(m Matcher) Matcher {
	return m.Or(Paren().X(Nil()).Y(m))
}

// KeyValue matches x:y.
func KeyValue(x, y Matcher) Matcher {
	return Op(":").X(x).Y(y)
}

// ArgItem matches if a set, seq or arglist has a specific item.
func ArgItem(item Matcher) Matcher {
	recurse := Matcher(func(np *ast.Node, tx *Tx) bool {
		return ArgItem(item)(np, tx)
	})
	return item.Or(Op(",").Y(item), Op(",").X(recurse))
}

// HasItem matches if the current node is a set, seq or arglist and
// its arglist has the item.
func HasItem(item Matcher) Matcher {
	return Y(ArgItem(item))
}

// HasKeyValue matches if a set, seq or arglist has a specific key value.
func HasKeyValue(key, value Matcher) Matcher {
	return HasItem(KeyValue(key, value))
}

// Capture captures the node.
func Capture(n *ast.Node) Matcher {
	return func(np *ast.Node, tx *Tx) bool {
		*n = *np
		return true
	}
}

// CaptureVal captures the node value.
func CaptureVal(s *string) Matcher {
	return func(np *ast.Node, tx *Tx) bool {
		n := *np
		switch n := n.(type) {
		case ast.Ident:
			*s = n.Val
		case ast.Number:
			*s = n.Val
		case ast.Quote:
			*s = n.Val
		}
		return true
	}
}

// Replace replaces the node with another.
func Replace(fn func() ast.Node) Matcher {
	return func(np *ast.Node, tx *Tx) bool {
		*np = fn()
		return true
	}
}

// Paren matches a Paren node.
func Paren() Matcher {
	return func(np *ast.Node, tx *Tx) bool {
		n := *np
		_, ok := n.(*ast.Paren)
		return ok
	}
}

// Set matches a Set node.
func Set() Matcher {
	return func(np *ast.Node, tx *Tx) bool {
		n := *np
		_, ok := n.(*ast.Set)
		return ok
	}
}

// Seq matches a Seq node.
func Seq() Matcher {
	return func(np *ast.Node, tx *Tx) bool {
		n := *np
		_, ok := n.(*ast.Seq)
		return ok
	}
}

// Nil matches a nil node.
func Nil() Matcher {
	return func(np *ast.Node, tx *Tx) bool {
		return *np == nil
	}
}

// NonNil matches a non-nil node.
func NonNil() Matcher {
	return func(np *ast.Node, tx *Tx) bool {
		return *np != nil
	}
}

// Identf matches against an identifier.
//
// The input func, if provided, is used to match for specific values
// of the identifier. Use Ident to match for fixed string names.
func Identf(fn func(val string, tx *Tx) bool) Matcher {
	fn = defaultStringMatcher(fn)
	return func(np *ast.Node, tx *Tx) bool {
		ident, ok := (*np).(ast.Ident)
		return ok && fn(ident.Val, tx)
	}
}

// Ident matches against specific identifier strings.
func Ident(ss ...string) Matcher {
	return Identf(func(s string, tx *Tx) bool {
		for _, str := range ss {
			if s == str {
				return true
			}
		}
		return false
	})
}

// Numberf matches against an number.
//
// The input func, if provided, is used to match for specific values
// of the number. Use Number to match for fixed numbers.
func Numberf(fn func(val string, tx *Tx) bool) Matcher {
	fn = defaultStringMatcher(fn)
	return func(np *ast.Node, tx *Tx) bool {
		num, ok := (*np).(ast.Number)
		return ok && fn(num.Val, tx)
	}
}

// Number matches against specific numeric strings.
func Number(ns ...float64) Matcher {
	return Numberf(func(s string, tx *Tx) bool {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return false
		}
		for _, nx := range ns {
			if f == nx {
				return true
			}
		}
		return false
	})
}

// Quotef matches against an quote.
//
// The input func, if provided, is used to match for specific values
// of the quoted string. Use Quote to match for fixed strings.
func Quotef(fn func(val string, tx *Tx) bool) Matcher {
	fn = defaultStringMatcher(fn)
	return func(np *ast.Node, tx *Tx) bool {
		q, ok := (*np).(ast.Quote)
		return ok && fn(q.Val, tx)
	}
}

// Quote matches against specific quoted strings.
//
// Note that the passed string must be unquoted.
func Quote(ss ...string) Matcher {
	unquote := func(s string) string {
		skip := false
		result := []rune{}
		for _, r := range s {
			skip = !skip && r == '\\'
			if !skip {
				result = append(result, r)
			}
		}
		return string(result)
	}
	return Quotef(func(val string, tx *Tx) bool {
		val = unquote(val[1 : len(val)-1])
		for _, s := range ss {
			if val == s {
				return true
			}
		}
		return false
	})
}

func defaultStringMatcher(fn func(val string, tx *Tx) bool) func(val string, tx *Tx) bool {
	if fn != nil {
		return fn
	}
	return func(val string, tx *Tx) bool {
		return true
	}
}
