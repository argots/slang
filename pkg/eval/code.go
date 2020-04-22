package eval

import (
	"bytes"

	"github.com/argots/slang/pkg/ast"
)

// Code represents a code fragment.
type Code struct {
	ast.Node
}

// String formats the code in a canonical text format.
func (c Code) String() string {
	f := &ast.TextFormatter{}
	var buf bytes.Buffer
	if err := f.Format(&buf, c.Node, &ast.FormatOptions{Formatter: f}); err != nil {
		panic(err)
	}
	return buf.String()
}
