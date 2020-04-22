package eval_test

import (
	"testing"

	"github.com/argots/slang/pkg/ast"
	"github.com/argots/slang/pkg/eval"
)

//nolint: lll
func TestEval(t *testing.T) {
	tests := map[string]string{
		"x":                             `sys.error{'undefined variable "x"'}`,
		"{x: 5}":                        `{"x": 5}`,
		"{5: 22}":                       `{5: 22}`,
		`5 + 5`:                         `10`,
		`10 / 2`:                        `5`,
		"6/4":                           `3 / 2`,
		"(-5)":                          `-5`,
		`"hello".length`:                `5`,
		`"hello".("length")`:            `5`,
		`{f(x,y): x + y}.f(1, 2)`:       `3`,
		`{f[x,y]: x + y}.f[1, 2]`:       `3`,
		`{f{x,y}: x + y}.f{y: 2, x: 1}`: `3`,
	}

	for test, want := range tests {
		got := evalString(test)
		if got != want {
			t.Errorf("%s: wanted %s but got %s", test, want, got)
		}
	}
}

func evalString(s string) string {
	n, err := ast.ParseString(s)
	if err != nil {
		return err.Error()
	}
	return eval.Node(n, eval.Globals()).Value().Code().String()
}
