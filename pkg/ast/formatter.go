package ast

import "io"

// Formatter is the interface to format a Node
type Formatter interface {
	Format(w io.Writer, n Node, options *FormatOptions) error
}

// FormatOptions defines the set of format options availqble
type FormatOptions struct {
	// Formatter is the default formatter to use to format a
	// node. This is used when a Node is recursively formaatted
	// allowing callers to wrap a formatter with another.
	Formatter
}

// TextFormatter implements a simple text formatting of a node
type TextFormatter struct{}

// Format formats a node.
func (f *TextFormatter) Format(w io.Writer, n Node, options *FormatOptions) error {
	if n == nil {
		return nil
	}

	x, ok := n.(*Expr)
	if !ok {
		v, _ := n.NodeInfo()
		_, err := w.Write([]byte(v))
		return err
	}

	ew := errWriter{nil, w, f}
	if options != nil && options.Formatter != nil {
		ew.f = options.Formatter
	}

	ew.format(x.X, options, f.needParen(x.Op, x.X, true))
	if x.X != nil {
		ew.write(" ")
	}
	ew.write(x.Op)
	if x.Y != nil {
		ew.write(" ")
	}
	ew.format(x.Y, options, f.needParen(x.Op, x.Y, false))
	return ew.err
}

func (f *TextFormatter) needParen(op string, n Node, isLeft bool) bool {
	x, ok := n.(*Expr)
	if !ok {
		return false
	}
	ownPri, xPri := priority(op), priority(x.Op)
	switch {
	case ownPri < xPri:
		return false
	case ownPri > xPri:
		return true
	case isLeft:
		return isRightAssoc(x.Op)
	default: // !isleft
		return !isRightAssoc(x.Op)
	}
}

type errWriter struct {
	err error
	w   io.Writer
	f   Formatter
}

func (ew *errWriter) write(s string) {
	if ew.err != nil {
		return
	}
	_, ew.err = ew.w.Write([]byte(s))
}

func (ew *errWriter) format(n Node, options *FormatOptions, useParen bool) {
	if ew.err != nil {
		return
	}
	if useParen {
		ew.write("(")
	}
	ew.err = ew.f.Format(ew.w, n, options)
	if useParen {
		ew.write(")")
	}
}
