package ast_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/argots/slang/pkg/ast"
)

func TestJSON(t *testing.T) {
	tests := []string{
		"1",
		"1.5",
		`""`,
		`"x"`,
		`'x\'y'`,
		"`\nx`",
		"x",
		"x123",
		"x'a b'",
		"-x",
		"x + y",
		"(x < y) | (x > y)",
		"x + y + z",
		"x + (y + z)",
		"x: y: z",
		"f.g(x[1], y{3:4})",
	}

	run := func(test string) func(t *testing.T) {
		return func(t *testing.T) {
			n, err := ast.ParseString(test)
			if err != nil {
				t.Fatal("parse", err)
			}
			data, err := json.Marshal(&ast.JSON{Node: n})
			if err != nil {
				t.Fatal("Marshal", err)
			}
			var result ast.JSON
			if err := json.Unmarshal(data, &result); err != nil {
				t.Fatal("Unmarshal", err)
			}

			var buf1, buf2 bytes.Buffer
			f := &ast.TextFormatter{}
			err = f.Format(&buf1, n, &ast.FormatOptions{Formatter: f})
			if err != nil {
				t.Fatal("format", err)
			}
			err = f.Format(&buf2, result.Node, &ast.FormatOptions{Formatter: f})
			if err != nil {
				t.Fatal("format", err)
			}
			if l, r := buf1.String(), buf2.String(); l != r {
				t.Error("parse/format diverged", l, "!=", r)
			}
		}
	}

	for _, test := range tests {
		t.Run(test, run(test))
	}
}
