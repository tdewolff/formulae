package formulae

import (
	"fmt"
	"math"

	"github.com/tdewolff/formulae/hash"
)

type Node interface {
	String() string
	Calc(Vars) (float64, error)
}

////////////////

type Func struct {
	Pos  int
	Name hash.Hash
	X    Node
}

func (n *Func) String() string {
	return fmt.Sprintf("(%v %v)", n.Name, n.X)
}

func (n *Func) Calc(vars Vars) (float64, error) {
	x, err := n.X.Calc(vars)
	if err != nil {
		return math.NaN(), err
	}

	switch n.Name {
	case hash.Sqrt:
		return math.Sqrt(x), nil
	case hash.Cbrt:
		return math.Cbrt(x), nil
	case hash.Sin:
		return math.Sin(x), nil
	case hash.Cos:
		return math.Cos(x), nil
	case hash.Tan:
		return math.Tan(x), nil
	case hash.Arcsin:
		return math.Asin(x), nil
	case hash.Arccos:
		return math.Acos(x), nil
	case hash.Arctan:
		return math.Atan(x), nil
	case hash.Sinh:
		return math.Sinh(x), nil
	case hash.Cosh:
		return math.Cosh(x), nil
	case hash.Tanh:
		return math.Tanh(x), nil
	case hash.Arcsinh:
		return math.Asinh(x), nil
	case hash.Arccosh:
		return math.Acosh(x), nil
	case hash.Arctanh:
		return math.Atanh(x), nil
	case hash.Exp:
		return math.Exp(x), nil
	case hash.Log, hash.Ln:
		if x < 0 {
			return math.NaN(), ParseErrorf(n.Pos, "logarithm of negative number")
		}
		return math.Log(x), nil
	case hash.Log2:
		if x < 0 {
			return math.NaN(), ParseErrorf(n.Pos, "logarithm of negative number")
		}
		return math.Log2(x), nil
	case hash.Log10:
		if x < 0 {
			return math.NaN(), ParseErrorf(n.Pos, "logarithm of negative number")
		}
		return math.Log10(x), nil
	case hash.Erf:
		return math.Erf(x), nil
	case hash.Gamma:
		return math.Gamma(x), nil
	default:
		return math.NaN(), ParseErrorf(n.Pos, "unknown function '%s'", n.Name)
	}
}

////////////////

type Expr struct {
	Pos int
	Op  Operator
	X   Node
	Y   Node
}

func (n *Expr) String() string {
	return fmt.Sprintf("(%v %v %v)", n.X, n.Op, n.Y)
}

func (n *Expr) Calc(vars Vars) (float64, error) {
	x, err := n.X.Calc(vars)
	if err != nil {
		return math.NaN(), err
	}

	y, err := n.Y.Calc(vars)
	if err != nil {
		return math.NaN(), err
	}

	switch n.Op {
	case AddOp:
		return x + y, nil
	case SubtractOp:
		return x - y, nil
	case MultiplyOp:
		return x * y, nil
	case DivideOp:
		if y == 0 {
			return math.NaN(), ParseErrorf(n.Pos, "division by zero")
		}
		return x / y, nil
	case PowerOp:
		return math.Pow(x, y), nil
	default:
		return math.NaN(), ParseErrorf(n.Pos, "unknown operation '%s'", n.Op)
	}
}

////////////////

type UnaryExpr struct {
	Pos int
	Op  Operator
	X   Node
}

func (n *UnaryExpr) String() string {
	return fmt.Sprintf("(%v %v)", n.Op, n.X)
}

func (n *UnaryExpr) Calc(vars Vars) (float64, error) {
	x, err := n.X.Calc(vars)
	if err != nil {
		return math.NaN(), err
	}
	return -x, nil
}

////////////////

type Variable struct {
	Pos  int
	Name string
}

func (n *Variable) String() string {
	return fmt.Sprintf("'%s'", n.Name)
}

func (n *Variable) Calc(vars Vars) (float64, error) {
	if val, ok := vars[n.Name]; ok {
		return val, nil
	}
	return math.NaN(), ParseErrorf(n.Pos, "undeclared variable '%s'", n.Name)
}

////////////////

type Number struct {
	Pos int
	Val float64
}

func (n *Number) String() string {
	return fmt.Sprintf("'%v'", n.Val)
}

func (n *Number) Calc(vars Vars) (float64, error) {
	return n.Val, nil
}
