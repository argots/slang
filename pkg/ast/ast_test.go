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
		{"(x): y"},
		{"((x)): y", "(x): y"},
	}

	run := func(test []string) func(t *testing.T) {
		return func(t *testing.T) {
			n, err := ast.ParseString(test[0])
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

func TestParseErrors(t *testing.T) { //nolint: funlen
	tests := map[string]string{
		"":      "unexpected terms count 0 at string:0",
		"()":    "unexpected terms count 0 at string:0",
		"  -":   "insufficient terms at string:2",
		"1 ++2": "missing term at string:3",
		"x ! 2": "unexpected character ! at string:2",
		"x $ 2": "unexpected character $ at string:2",
		"x (":   "unexpected EOF",
		"x y":   "missing op at string:2",
		"x (}":  "unexpected close at string:3",
		"x }":   "unexpected close at string:2",
	}

	run := func(test string) func(t *testing.T) {
		return func(t *testing.T) {
			n, err := ast.ParseString(test)
			if err == nil {
				t.Error("parse", n)
			}
			if err.Error() != tests[test] {
				t.Error("Unexpected", err)
			}
		}
	}

	for test := range tests {
		t.Run(test, run(test))
	}
}
