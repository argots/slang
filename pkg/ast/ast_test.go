package ast_test

import (
	"bytes"
	"testing"

	"github.com/argots/slang/pkg/ast"
)

func TestParseAndFormat(t *testing.T) { //nolint: funlen
	tests := [][]string{
		{"1"},
		{"1.5"},
		{"   1.5  \n", "1.5"},
		{`""`},
		{`"x"`},
		{`'x\'y'`},
		{"`\nx`"},
		{"x"},
		{"x123"},
		{"x'a b'"},
		{"-x", "- x"},
		{"x + y"},
		{"(x < y) | (x > y)", "x < y | x > y"},
		{"(x <= y) & (x >= y)", "x <= y & x >= y"},
		{"(x = y) | (x != y)", "x = y | x != y"},
		{"x + y + z"},
		{"x + (y + z)", "x + (y + z)"},
		{"(x + y) - z", "x + y - z"},
		{"x + -5", "x + - 5"},
		{"x ---5", "x - - - 5"},
		{"x - y - z"},
		{"x - ( y - z )", "x - (y - z)"},
		{"x + y * 5"},
		{"x: y: z"},
		{"x <= y"},
		{"x & y | a & b"},
		{"(x&y) | (a&b)", "x & y | a & b"},
		{"(x)", "x"},
		{"(x + y)*5", "(x + y) * 5"},
		{"(x , a) : 5", "(x, a): 5"},
		{"(x : y) : z", "(x: y): z"},
		{"x : (y : z)", "x: y: z"},
		{"[  ]", "[]"},
		{"x []", "x[]"},
		{"x. y. z []", "x.y.z[]"},
		{"(x + y)[5, 2]", "(x + y)[5, 2]"},
		{"x + (f[23])", "x + f[23]"},
		{"x + [23, 24]", "x + [23, 24]"},
		{"{}"},
		{"map{[1, 2]: 42}"},
	}

	run := func(test []string) func(t *testing.T) {
		return func(t *testing.T) {
			s, lm := &ast.Sources{}, ast.NewLocMap()
			s.AddStringSource("test", test[0])
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
			if x := buf.String(); x != test[len(test)-1] {
				t.Error("parse/format diverged", x)
			}
		}
	}

	for _, test := range tests {
		t.Run(test[0], run(test))
	}
}
