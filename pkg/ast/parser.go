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
	return p.parse()
}

type parser struct {
	tokenizer
	lastWasTerm bool
	ops         []*token
	terms       []Node
}

func (p *parser) parse() (Node, error) {
	for {
		tok, err := p.Next()
		switch {
		case err == io.EOF:
			return p.finish()
		case err != nil:
		case tok.Kind == operatorToken:
			err = p.handleOp(tok)
		case tok.Kind == numberToken:
			err = p.handleTerm(Number{tok.Value, tok.Loc})
		case tok.Kind == quoteToken:
			err = p.handleTerm(Quote{tok.Value, tok.Loc})
		case tok.Kind == identToken:
			err = p.handleTerm(Ident{tok.Value, tok.Loc})
		}

		if err != nil {
			return nil, err
		}
	}
}

func (p *parser) finish() (Node, error) {
	if err := p.unwindOps(""); err != nil {
		return nil, err
	}
	if len(p.terms) != 1 {
		return nil, fmt.Errorf("unexpected terms count %v %v", p.terms, p.ops)
	}
	return p.terms[0], nil
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
