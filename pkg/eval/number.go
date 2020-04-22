package eval

import (
	"math/big"

	"github.com/argots/slang/pkg/ast"
	"github.com/argots/slang/pkg/cast"
)

var _ Value = numValue{}

// NewNumber creates a numeric value from a float64
func NewNumber(f float64) Value {
	v := numValue{&big.Rat{}}
	v.SetFloat64(f)
	return v
}

type numValue struct {
	*big.Rat
}

func (n numValue) Type() string {
	return "sys.number"
}

func (n numValue) Code() Code {
	if n.IsInt() {
		return Code{ast.Number{Val: n.RatString()}}
	}

	num, denom := n.Num().String(), n.Denom().String()
	return Code{cast.Expr("/", ast.Number{Val: num}, ast.Number{Val: denom}).Node}
}

func (n numValue) Value() Value {
	return n
}

func (n numValue) Get(v Valuable) Valuable {
	return NewError(NewString("no such field " + toString(v)))
}

func (n numValue) Arithmetic(op string, other numValue) Valuable {
	var r big.Rat

	switch op {
	case "+":
		return numValue{r.Add(n.Rat, other.Rat)}
	case "-":
		return numValue{r.Sub(n.Rat, other.Rat)}
	case "*":
		return numValue{r.Mul(n.Rat, other.Rat)}
	case "/":
		var yinv big.Rat
		yinv.Inv(other.Rat)
		return numValue{r.Mul(n.Rat, &yinv)}
	}
	return NewError(NewString("Unknown op: " + op))
}
