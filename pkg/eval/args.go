package eval

import "github.com/argots/slang/pkg/ast"

// Args implements a simple visitor for the comma-separated list of
// args in function calls, sets and sequences.
type Args struct {
	NoKey     func(val ast.Node) bool
	StringKey func(key string, val ast.Node) bool
	NodeKey   func(key, val ast.Node) bool
	ParenKey  func(name string, args, val ast.Node) bool
	SetKey    func(name string, args, val ast.Node) bool
	SeqKey    func(name string, args, val ast.Node) bool
}

func (a Args) Visit(n ast.Node) {
	loop := true
	for n != nil && loop {
		if comma, ok := n.(*ast.Expr); ok && comma.Op == "," {
			loop = !a.visitArg(comma.X)
			n = comma.Y
			continue
		}
		a.visitArg(n)
		break
	}
}

func (a Args) visitArg(n ast.Node) bool {
	expr, ok := n.(*ast.Expr)
	if !ok || expr.Op != ":" {
		if ident, ok := n.(ast.Ident); ok {
			return a.StringKey(ident.Val, n)
		}
		return a.NoKey(n)
	}
	return a.visitArgWithKey(expr.X, expr.Y)
}

func (a Args) visitArgWithKey(key, val ast.Node) bool {
	switch key := key.(type) {
	case ast.Ident:
		return a.StringKey(key.Val, val)
	case *ast.Paren:
		if ident, ok := key.X.(ast.Ident); ok {
			return a.ParenKey(ident.Val, key.Y, val)
		}
	case *ast.Set:
		if ident, ok := key.X.(ast.Ident); ok {
			return a.SetKey(ident.Val, key.Y, val)
		}
	case *ast.Seq:
		if ident, ok := key.X.(ast.Ident); ok {
			return a.SeqKey(ident.Val, key.Y, val)
		}
	}
	return a.NodeKey(key, val)
}
