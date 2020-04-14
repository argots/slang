package ast

// type assertions
var _ = []Node{&Expr{}, Number{}, Quote{}, Ident{}}

// Node is the main interface implemented by all nodes in the AST
type Node interface {
	NodeInfo() (value string, loc Loc)
}

// Expr represents an expression of form X Op Y.
// For unary expressions, X will be nil.
type Expr struct {
	Op string
	Loc
	X, Y Node
}

// NodeInfo returns the operator string and the location.
func (x *Expr) NodeInfo() (value string, loc Loc) {
	return x.Op, x.Loc
}

// Number represents a numeric literal
type Number struct {
	Val string
	Loc
}

// NodeInfo returns the raw numeric string and its location in the source code.
func (n Number) NodeInfo() (value string, loc Loc) {
	return n.Val, n.Loc
}

// Quote represents a quoted string
type Quote struct {
	Val string
	Loc
}

// NodeInfo returns the raw string (including the open/close quote and
// any embedded slashes) and its location in the source code.
func (q Quote) NodeInfo() (value string, loc Loc) {
	return q.Val, q.Loc
}

// Ident represents an identifier
type Ident struct {
	Val string
	Loc
}

// Ident returns the raw identifier string (including if has a quoted string)
func (i Ident) NodeInfo() (value string, loc Loc) {
	return i.Val, i.Loc
}
