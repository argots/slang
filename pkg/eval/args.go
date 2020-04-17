package eval

import "github.com/argots/slang/pkg/ast"

func Args(n ast.Node, s Scope, fn func(key Value, v ast.Node)) {
	loop := n != nil
	for loop {
		if comma, ok := n.(*ast.Expr); ok && comma.Op == "," {
			arg(comma.X, s, fn)
			n = comma.Y
		} else {
			arg(n, s, fn)
			loop = false
		}
	}
}

func arg(n ast.Node, s Scope, fn func(key Value, v ast.Node)) {
	if expr, ok := n.(*ast.Expr); ok && expr.Op == ":" {
		if ident, ok := expr.X.(ast.Ident); ok {
			fn(NewString(ident.Val), expr.Y)
		} else {
			fn(Node(expr.X, s).Value(), expr.Y)
		}
	} else if ident, ok := n.(ast.Ident); ok {
		fn(NewString(ident.Val), n)
	} else {
		fn(NewString(""), n)
	}
}
