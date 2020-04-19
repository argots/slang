// Package cast provides methods to create and build AST nodes.
package cast

import (
	"strconv"
	"strings"

	"github.com/argots/slang/pkg/ast"
)

// Node provides some helpers on top of ast.Node to allow chaining
// methods.
//
// Example:
//
//      n := Node{..}
//      n.Dot("y").Add(5)
type Node struct {
	ast.Node
}

// Dot implements field access.
func (n Node) Dot(args ...interface{}) Node {
	return Dot(n, args...)
}

// Call implements function calls.
func (n Node) Call(args ...interface{}) Node {
	return Call(n, args...)
}

// Seq implements [].
func (n Node) Seq(args ...interface{}) Node {
	return Seq(n, args...)
}

// Set implements {}.
func (n Node) Set(args ...interface{}) Node {
	return Set(n, args...)
}

// Add implements +.
func (n Node) Add(other interface{}) Node {
	return Expr("+", n, other)
}

// Sub implements -
func (n Node) Sub(other interface{}) Node {
	return Expr("-", n, other)
}

// Neg implements negation (unary -).
func (n Node) Neg() Node {
	return Expr("-", nil, n)
}

// Mul implements *.
func (n Node) Mul(other interface{}) Node {
	return Expr("*", n, other)
}

// Div implements /.
func (n Node) Div(other interface{}) Node {
	return Expr("/", n, other)
}

// ToNode accepts strings, numbers, arrays and maps.
//
// Note that strings map to identifiers by default. If quotes are
// needed, use Quote explicitly.
//
// Arrays and maps must be typed as []interface{} and
// map[interface{}]interface{} respectively.
//
// For function/seq/set calls, use Call/Set/Seq instead of ToNode.
func ToNode(x interface{}) Node {
	switch x := x.(type) {
	case nil:
		return Node{nil}
	case string:
		return Node{ast.Ident{Val: FormatIdent(x)}}
	case float64:
		return Node{ast.Number{Val: strconv.FormatFloat(x, 'f', -1, 64)}}
	case int:
		return Node{ast.Number{Val: strconv.Itoa(x)}}
	case []interface{}:
		return Seq(nil, x...)
	case map[interface{}]interface{}:
		args := []interface{}{}
		for k, v := range x {
			args = append(args, Pair(k, v))
		}
		return Set(nil, args...)
	case Node:
		return x
	case ast.Node:
		return Node{x}
	}
	panic("unexpected type")
}

// FormatIdent formats any string into a form that is valid for use
// with slang.
//
// The first rune must be a letter but otherwise all strings are
// allowed.  Note that vaild strings like `xyz"a"` will still be
// double encoded (into `xzy"\"a\""`).
func FormatIdent(s string) string {
	// NYI
	return s
}

// Quote creates a quoted string.
func Quote(s string) Node {
	ends := `"`
	switch {
	case !strings.Contains(s, `"`):
		ends = `"`
	case !strings.Contains(s, `'`):
		ends = `'`
	case !strings.Contains(s, "`"):
		ends = "`"
	default:
		s = strings.ReplaceAll(s, `"`, `\"`)
	}

	return Node{ast.Quote{Val: ends + s + ends}}
}

// Expr creates a general expr node.
func Expr(op string, x, y interface{}) Node {
	return Node{&ast.Expr{Op: op, X: ToNode(x).Node, Y: ToNode(y).Node}}
}

// Pair creates a key value pair.
func Pair(x, y interface{}) Node {
	return Expr(":", x, y)
}

// Dot creates a dot expression
func Dot(x interface{}, y ...interface{}) Node {
	result := ToNode(x)
	for _, item := range y {
		result = Expr(".", result, item)
	}
	return result
}

// Call creates a Paren expr X(args...)
//
// x and args can be anything that can be passed to ToNode
func Call(x interface{}, args ...interface{}) Node {
	return Node{&ast.Paren{
		StartOp: "(",
		EndOp:   ")",
		X:       ToNode(x).Node,
		Y:       ArgsList(args...).Node,
	}}
}

// Seq creates a Seq expr X[args...]
//
// x and args can be anything that can be passed to ToNode
func Seq(x interface{}, args ...interface{}) Node {
	return Node{&ast.Set{
		StartOp: "[",
		EndOp:   "]",
		X:       ToNode(x).Node,
		Y:       ArgsList(args...).Node,
	}}
}

// Set creates a Seq expr X[args...]
//
// x and args can be anything that can be passed to ToNode
func Set(x interface{}, args ...interface{}) Node {
	return Node{&ast.Set{
		StartOp: "{",
		EndOp:   "}",
		X:       ToNode(x).Node,
		Y:       ArgsList(args...).Node,
	}}
}

// ArgsList converts the args into a comma separted list
func ArgsList(args ...interface{}) Node {
	var result Node
	for _, arg := range args {
		if result.Node == nil {
			result = ToNode(arg)
		} else {
			result = Expr(",", result, arg)
		}
	}
	return result
}
