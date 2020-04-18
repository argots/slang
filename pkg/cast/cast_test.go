package cast_test

import (
	"bytes"
	"testing"

	"github.com/argots/slang/pkg/ast"
	"github.com/argots/slang/pkg/cast"
)

func TestSuccess(t *testing.T) {
	tests := map[string]interface{}{
		"1":         1,
		"1.3":       1.3,
		"x":         "x",
		"[1, 2, 3]": []interface{}{1, 2, 3},
		"{x: 5}":    map[interface{}]interface{}{"x": 5},
		`{[1, 2, 3]: "a"}`: map[interface{}]interface{}{
			cast.Seq(nil, 1, 2, 3): cast.Quote("a"),
		},
		"f[1, 2, 3]":   cast.Seq("f", 1, 2, 3),
		"f2[1, 2, 3]":  cast.ToNode("f2").Seq(1, 2, 3),
		"f.g.h{x: 5}":  cast.Set(cast.Dot("f", "g", "h"), cast.Pair("x", 5)),
		"f2.g.h{x: 5}": cast.ToNode("f2").Dot("g", "h").Set(cast.Pair("x", 5)),
		"f.(g.h)(1, 2)": cast.Call(
			cast.Dot("f", cast.Dot("g", "h")),
			1,
			2,
		),
		"f2.(g.h)(1, 2)":        cast.ToNode("f2").Dot(cast.ToNode("g").Dot("h")).Call(1, 2),
		"x + y * 3 / 4":         cast.ToNode("x").Add(cast.ToNode("y").Mul(3).Div(4)),
		"- x - 5":               cast.ToNode("x").Neg().Sub(5),
		"z":                     cast.ToNode(cast.ToNode("z").Node),
		`'hello" world'`:        cast.Quote("hello\" world"),
		"`hello\"' world`":      cast.Quote("hello\"' world"),
		"\"hello\\\"'` world\"": cast.Quote("hello\"'` world"),
	}

	for want, val := range tests {
		if got := toString(cast.ToNode(val).Node); got != want {
			t.Errorf("Expected %s, got %s", want, got)
		}
	}
}

func toString(n ast.Node) string {
	f := &ast.TextFormatter{}
	var buf bytes.Buffer
	if err := f.Format(&buf, n, &ast.FormatOptions{Formatter: f}); err != nil {
		panic(err)
	}
	return buf.String()
}
