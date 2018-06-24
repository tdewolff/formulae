package formulae

import (
	"fmt"
	"math/cmplx"

	"github.com/tdewolff/formulae/hash"
)

type Node interface {
	String() string
	Calc(Vars) (complex128, error)
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

func (n *Func) Calc(vars Vars) (complex128, error) {
	x, err := n.X.Calc(vars)
	if err != nil {
		return cmplx.NaN(), err
	}

	switch n.Name {
	case hash.Sqrt:
		return cmplx.Sqrt(x), nil
	case hash.Sin:
		return cmplx.Sin(x), nil
	case hash.Cos:
		return cmplx.Cos(x), nil
	case hash.Tan:
		return cmplx.Tan(x), nil
	case hash.Arcsin:
		return cmplx.Asin(x), nil
	case hash.Arccos:
		return cmplx.Acos(x), nil
	case hash.Arctan:
		return cmplx.Atan(x), nil
	case hash.Sinh:
		return cmplx.Sinh(x), nil
	case hash.Cosh:
		return cmplx.Cosh(x), nil
	case hash.Tanh:
		return cmplx.Tanh(x), nil
	case hash.Arcsinh:
		return cmplx.Asinh(x), nil
	case hash.Arccosh:
		return cmplx.Acosh(x), nil
	case hash.Arctanh:
		return cmplx.Atanh(x), nil
	case hash.Exp:
		return cmplx.Exp(x), nil
	case hash.Log, hash.Ln:
		return cmplx.Log(x), nil
	case hash.Log10:
		return cmplx.Log10(x), nil
	default:
		return cmplx.NaN(), ParseErrorf(n.Pos, "unknown function '%s'", n.Name)
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

func (n *Expr) Calc(vars Vars) (complex128, error) {
	x, err := n.X.Calc(vars)
	if err != nil {
		return cmplx.NaN(), err
	}

	y, err := n.Y.Calc(vars)
	if err != nil {
		return cmplx.NaN(), err
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
			return cmplx.Inf(), ParseErrorf(n.Pos, "division by zero") // TODO: set sign
		}
		return x / y, nil
	case PowerOp:
		return cmplx.Pow(x, y), nil
	default:
		return cmplx.NaN(), ParseErrorf(n.Pos, "unknown operation '%s'", n.Op)
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

func (n *UnaryExpr) Calc(vars Vars) (complex128, error) {
	x, err := n.X.Calc(vars)
	return -x, err
}

////////////////

type Variable struct {
	Pos  int
	Name string
}

func (n *Variable) String() string {
	return fmt.Sprintf("'%s'", n.Name)
}

func (n *Variable) Calc(vars Vars) (complex128, error) {
	if val, ok := vars[n.Name]; ok {
		return val, nil
	}
	return cmplx.NaN(), ParseErrorf(n.Pos, "undeclared variable '%s'", n.Name)
}

////////////////

type Number struct {
	Pos int
	Val complex128
}

func (n *Number) String() string {
	if imag(n.Val) == 0 {
		return fmt.Sprintf("'%v'", real(n.Val))
	}
	return fmt.Sprintf("'%v'", n.Val)
}

func (n *Number) Calc(vars Vars) (complex128, error) {
	return n.Val, nil
}
