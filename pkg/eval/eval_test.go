package eval_test

import (
	"testing"

	"github.com/argots/slang/pkg/ast"
	"github.com/argots/slang/pkg/eval"
)

func TestEval(t *testing.T) {
	tests := map[string]string{
		`"hello".length`:     `"5"`,
		`"hello".("length")`: `"5"`,
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
	return eval.Node(n, eval.Globals()).Value().Code()
}
