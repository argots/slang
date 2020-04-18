package main

import (
	"fmt"

	"github.com/argots/slang/pkg/ast"
	"github.com/argots/slang/pkg/eval"
)

func main() {
	tests := map[string]string{
		"x":                             `sys.error{'undefined variable "x"'}`,
		"{x: 5}":                        `{"x": 5}`,
		"{5: 22}":                       `{5: 22}`,
		`5 + 5`:                         `10`,
		`10 / 2`:                        `5`,
		"6/4":                           `(3/2)`,
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
			fmt.Printf("%s: wanted %s but got %s\n", test, want, got)
		} else {
			fmt.Printf("%s: ok\n", test)
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
