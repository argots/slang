package ast_test

import (
	"testing"

	"github.com/argots/slang/pkg/ast"
)

func TestLocMap(t *testing.T) {
	lm := ast.NewLocMap()
	h := lm.Add("hello", 2, 3)
	h2 := lm.Add("hello", 2, 3)
	if h != h2 {
		t.Fatal("Unexpected duplicate add", h, h2)
	}

	if loc, start, end := h.Offset(lm); loc != "hello" || start != 2 || end != 3 {
		t.Fatal("Unexpected loc, start, end", loc, start, end)
	}

	s := &ast.Sources{}
	s.AddStringSource("hello", "world")

	if tok, err := h.Token(lm, s); tok != "r" || err != nil {
		t.Fatal("Unexpected token", tok, err)
	}
}
