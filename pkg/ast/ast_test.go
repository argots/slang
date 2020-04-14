package ast_test

import (
	"bytes"
	"testing"

	"github.com/argots/slang/pkg/ast"
)

func TestSuccessfulCanonicalParsing(t *testing.T) {
	tests := []string{
		"1",
		"1.5",
		`""`,
		`"x"`,
		`'x\'y'`,
		"`x`",
		"x",
		"x123",
		"x'a b'",
		"x + y",
		"x + y + z",
		"x + y - z",
		"x - y - z",
		"x + y * 5",
		"x : y : z",
	}

	run := func(test string) func(t *testing.T) {
		return func(t *testing.T) {
			s, lm := &ast.Sources{}, ast.NewLocMap()
			s.AddStringSource("test", test)
			n, err := ast.Parse(s.ReadSource("test"), "test", lm)
			if err != nil {
				t.Fatal("parse", err)
			}
			var buf bytes.Buffer
			f := &ast.TextFormatter{}
			err = f.Format(&buf, n, &ast.FormatOptions{Formatter: f})
			if err != nil {
				t.Fatal("format", err)
			}
			if x := buf.String(); x != test {
				t.Error("parse/format diverged", x)
			}
		}
	}

	for _, test := range tests {
		t.Run(test, run(test))
	}
}
