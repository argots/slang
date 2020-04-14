package ast

import (
	"errors"
	"fmt"
	"io"
)

// Parse parses the source and returns the parsed AST
func Parse(r io.Reader, location string, lm LocMap) (Node, error) {
	t := tokenizer{Reader: r, Location: location, LocMap: lm}
	p := parser{tokenizer: t}
	n, _, err := p.parse("", false)
	return p.stripParen(n), err
}

type parser struct {
	tokenizer
	lastWasTerm bool
	ops         []*token
	terms       []Node
}

func (p *parser) parse(end string, allowEmpty bool) (Node, Loc, error) {
	for {
		tok, err := p.Next()
		switch {
		case err == io.EOF && end == "":
			return p.finish(Loc(0), false)
		case err == io.EOF:
			return nil, Loc(0), io.ErrUnexpectedEOF
		case err != nil:
		case tok.Kind == operatorToken && tok.Value == end:
			return p.finish(tok.Loc, allowEmpty)
		case tok.Kind == operatorToken:
			switch tok.Value {
			case "(", "[", "{":
				err = p.handleSetSeqOrParen(tok)
			case ")", "]", "}":
				err = errors.New("unexpected close")
			default:
				err = p.handleOp(tok)
			}
		case tok.Kind == numberToken:
			err = p.handleTerm(Number{tok.Value, tok.Loc})
		case tok.Kind == quoteToken:
			err = p.handleTerm(Quote{tok.Value, tok.Loc})
		case tok.Kind == identToken:
			err = p.handleTerm(Ident{tok.Value, tok.Loc})
		}

		if err != nil {
			return nil, Loc(0), err
		}
	}
}

func (p *parser) finish(l Loc, allowEmpty bool) (Node, Loc, error) {
	if err := p.unwindOps(""); err != nil {
		return nil, l, err
	}
	if allowEmpty && len(p.terms) == 0 {
		return nil, l, nil
	}
	if len(p.terms) != 1 {
		return nil, l, fmt.Errorf("unexpected terms count %v %v", p.terms, p.ops)
	}
	return p.terms[0], l, nil
}

func (p *parser) handleOp(tok *token) error {
	if !p.lastWasTerm && !isUnary(tok.Value) {
		return errors.New("missing term")
	}

	if !p.lastWasTerm {
		p.terms = append(p.terms, nil)
	}

	if err := p.unwindOps(tok.Value); err != nil {
		return err
	}

	p.ops = append(p.ops, tok)
	p.lastWasTerm = false
	return nil
}

func (p *parser) handleSetSeqOrParen(tok *token) error {
	close := p.closeToken(tok.Value)
	allowEmpty := p.lastWasTerm || tok.Value != "("

	var x Node
	var l Loc
	if p.lastWasTerm {
		p.lastWasTerm = false
		if err := p.unwindOps(tok.Value); err != nil {
			return err
		}
		x = p.terms[len(p.terms)-1]
		p.terms = p.terms[:len(p.terms)-1]
	}

	p2 := parser{tokenizer: p.tokenizer}
	y, l, err := p2.parse(close, allowEmpty)
	if err != nil {
		return err
	}
	p.tokenizer = p2.tokenizer

	x, y = p.stripParen(x), p.stripParen(y)
	switch close {
	case ")":
		return p.handleTerm(&Paren{tok.Value, close, tok.Loc, l, x, y})
	case "]":
		return p.handleTerm(&Seq{tok.Value, close, tok.Loc, l, x, y})
	default:
		return p.handleTerm(&Set{tok.Value, close, tok.Loc, l, x, y})
	}
}

func (p *parser) unwindOps(op string) error {
	pri := priority(op)
	isRightAssociative := isRightAssoc(op)
	for l := len(p.ops) - 1; l >= 0 && priority(p.ops[l].Value) >= pri; l-- {
		tok := p.ops[l]
		if isRightAssociative && tok.Value == op {
			break
		}

		p.ops = p.ops[:l]
		if len(p.terms) <= 1 {
			return errors.New("insufficient terms")
		}
		x, y := p.terms[len(p.terms)-2], p.terms[len(p.terms)-1]
		p.terms = p.terms[:len(p.terms)-2]
		p.terms = append(p.terms, p.newExpr(tok.Value, tok.Loc, x, y))
	}
	return nil
}

func (p *parser) handleTerm(n Node) error {
	if p.lastWasTerm {
		return errors.New("missing op")
	}
	p.terms = append(p.terms, n)
	p.lastWasTerm = true
	return nil
}

func (p *parser) newExpr(op string, loc Loc, x, y Node) Node {
	if op != ":" {
		x = p.stripParen(x)
	}
	return &Expr{Op: op, Loc: loc, X: x, Y: p.stripParen(y)}
}

func (p *parser) stripParen(n Node) Node {
	if v, ok := n.(*Paren); ok && v.X == nil {
		return v.Y
	}
	return n
}

func (p *parser) closeToken(s string) string {
	switch s {
	case "(":
		return ")"
	case "[":
		return "]"
	default:
		return "}"
	}
}
