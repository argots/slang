package eval

import "math/big"

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

func (n numValue) Code() string {
	if n.IsInt() {
		return n.RatString()
	}

	return "(" + n.RatString() + ")"
}

func (n numValue) Value() Value {
	return n
}

func (n numValue) Get(v Valuable) Valuable {
	return NewError(NewString("no such field " + v.Value().Code()))
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
