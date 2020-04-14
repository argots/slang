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
	n, _, err := p.parse("")
	return n, err
}

type parser struct {
	tokenizer
	lastWasTerm bool
	ops         []*token
	terms       []Node
}

func (p *parser) parse(end string) (Node, Loc, error) {
	for {
		tok, err := p.Next()
		switch {
		case err == io.EOF && end == "":
			return p.finish(Loc(0), false)
		case err == io.EOF:
			return nil, Loc(0), io.ErrUnexpectedEOF
		case err != nil:
		case tok.Kind == operatorToken && tok.Value == end:
			return p.finish(tok.Loc, true)
		case tok.Kind == operatorToken && tok.Value == "(":
			err = p.handleParen()
		case tok.Kind == operatorToken && (tok.Value == "[" || tok.Value == "]"):
			err = p.handleSetOrSeq(tok)
		case tok.Kind == operatorToken:
			if tok.Value == ")" || tok.Value == "]" || tok.Value == "}" {
				err = errors.New("unexpected close")
			} else {
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
	if !p.lastWasTerm {
		return errors.New("missing term")
	}

	if err := p.unwindOps(tok.Value); err != nil {
		return err
	}

	p.ops = append(p.ops, tok)
	p.lastWasTerm = false
	return nil
}

func (p *parser) handleParen() error {
	p2 := parser{tokenizer: p.tokenizer}
	n, _, err := p2.parse(")")
	if err != nil {
		return err
	}
	p.tokenizer = p2.tokenizer
	return p.handleTerm(n)
}

func (p *parser) handleSetOrSeq(tok *token) error {
	close := "]"
	if tok.Value == "{" {
		close = "}"
	}
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
	y, l, err := p2.parse(close)
	if err != nil {
		return err
	}
	p.tokenizer = p2.tokenizer

	if close == "]" {
		return p.handleTerm(&Seq{tok.Value, close, tok.Loc, l, x, y})
	}
	return p.handleTerm(&Set{tok.Value, close, tok.Loc, l, x, y})
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
		p.terms = append(p.terms, &Expr{Op: tok.Value, Loc: tok.Loc, X: x, Y: y})
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
