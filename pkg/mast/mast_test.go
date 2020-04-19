package mast_test

import (
	"testing"

	"github.com/argots/slang/pkg/ast"
	"github.com/argots/slang/pkg/mast"
)

func TestExamples(t *testing.T) {
	x := mast.Ident("x")
	y := mast.Ident("y")
	f := mast.Ident("f")
	g := mast.Ident("g")
	one := mast.Number(1)
	two5 := mast.Number(2.5)
	boo := mast.Quote("boo").MaybeParen()
	var xy ast.Node
	rxy := func() ast.Node {
		return xy
	}
	var xs string
	tests := map[string]mast.Matcher{
		"x + y":               mast.Op("+").X(x).Y(y),
		"{x: y}":              mast.Set().X(mast.Nil()).Y(mast.KeyValue(x, y)),
		"{x: 2.5}":            mast.Nil().Set(mast.KeyValue(x, two5)),
		"{x: 1, y: 2.5}":      mast.Set().HasKeyValue(x, one).HasKeyValue(y, two5),
		"f.g(1, 2.5)":         f.Dot(g).Call(one, two5),
		"f{1, 2.5, 3}":        mast.Any().HasItem(two5),
		"(f.g)[1, 2.5]":       f.Dot(g).Seq(one, two5),
		`f.("boo")[1, 2.5]`:   f.Dot(boo).Seq(one, two5),
		`f.(("boo"))[1, 2.5]`: f.Dot(boo).Seq(one, two5),

		// capture the value and confirm it matches
		"x": x.CaptureVal(&xs).And(mast.Identf(func(s string, tx *mast.Tx) bool {
			return s == xs
		})),

		// capture the parent of 'x' into xy, replace everything
		// with that and then match it again to ensure it is x + y
		"f(x + y + z)": mast.NonNil().Contains(mast.X(x).Capture(&xy)).
			Replace(rxy).And(mast.Op("+").X(x).Y(y)),
	}

	for code, matcher := range tests {
		n, err := ast.ParseString(code)
		if err != nil {
			t.Fatal("Failed to parse", err)
		}
		if !matcher(&n, &mast.Tx{}) {
			t.Error("Failed to match", code)
		}
	}
}
