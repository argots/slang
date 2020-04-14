package ast

// priority returns the priority of an operator.
// returns -1 if the operator isn't found
//nolint: gomnd
func priority(op string) int {
	switch op {
	case ",":
		return 200
	case ":":
		return 300
	case "|":
		return 400
	case "&":
		return 500
	case "=":
		return 600
	case "!=":
		return 700
	case "<", ">", "<=", ">=":
		return 800
	case "+", "-":
		return 900
	case "*", "/":
		return 1000
	case "(", ")", "{", "}", "[", "]":
		return 1200
	case ".":
		return 1300
	}

	return -1
}

func isRightAssoc(op string) bool {
	return op == ":"
}

func isUnary(op string) bool {
	return op == "-"
}
