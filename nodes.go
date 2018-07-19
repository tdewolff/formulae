package formulae

import (
	"fmt"
	"math/cmplx"

	"github.com/tdewolff/formulae/hash"
)

type Node interface {
	String() string
	Calc(complex128, Vars) (complex128, error)
	Lua() string
}

////////////////

type Func struct {
	pos  int
	name hash.Hash
	a    Node
}

func (n *Func) String() string {
	return fmt.Sprintf("(%v %v)", n.name, n.a)
}

func (n *Func) Calc(x complex128, vars Vars) (complex128, error) {
	a, err := n.a.Calc(x, vars)
	if err != nil {
		return cmplx.NaN(), err
	}

	switch n.name {
	case hash.Sqrt:
		return cmplx.Sqrt(a), nil
	case hash.Sin:
		return cmplx.Sin(a), nil
	case hash.Cos:
		return cmplx.Cos(a), nil
	case hash.Tan:
		return cmplx.Tan(a), nil
	case hash.Arcsin:
		return cmplx.Asin(a), nil
	case hash.Arccos:
		return cmplx.Acos(a), nil
	case hash.Arctan:
		return cmplx.Atan(a), nil
	case hash.Sinh:
		return cmplx.Sinh(a), nil
	case hash.Cosh:
		return cmplx.Cosh(a), nil
	case hash.Tanh:
		return cmplx.Tanh(a), nil
	case hash.Arcsinh:
		return cmplx.Asinh(a), nil
	case hash.Arccosh:
		return cmplx.Acosh(a), nil
	case hash.Arctanh:
		return cmplx.Atanh(a), nil
	case hash.Exp:
		return cmplx.Exp(a), nil
	case hash.Log, hash.Ln:
		return cmplx.Log(a), nil
	case hash.Log10:
		return cmplx.Log10(a), nil
	default:
		return cmplx.NaN(), ParseErrorf(n.pos, "unknown function '%s'", n.name)
	}
}

func (n *Func) Lua() string {
	return "math." + n.name.String() + "(" + n.a.Lua() + ")"
}

////////////////

type Expr struct {
	pos int
	op  Operator
	a   Node
	b   Node
}

func (n *Expr) String() string {
	return fmt.Sprintf("(%v %v %v)", n.a, n.op, n.b)
}

func (n *Expr) Calc(x complex128, vars Vars) (complex128, error) {
	a, err := n.a.Calc(x, vars)
	if err != nil {
		return cmplx.NaN(), err
	}

	b, err := n.b.Calc(x, vars)
	if err != nil {
		return cmplx.NaN(), err
	}

	switch n.op {
	case AddOp:
		return a + b, nil
	case SubtractOp:
		return a - b, nil
	case MultiplyOp:
		return a * b, nil
	case DivideOp:
		if b == 0 {
			return cmplx.Inf(), ParseErrorf(n.pos, "division by zero") // TODO: set sign
		}
		return a / b, nil
	case PowerOp:
		return cmplx.Pow(a, b), nil
	default:
		return cmplx.NaN(), ParseErrorf(n.pos, "unknown operation '%s'", n.op)
	}
}

func (n *Expr) Lua() string {
	return n.a.Lua() + n.op.String() + n.b.Lua()
}

////////////////

type UnaryExpr struct {
	pos int
	op  Operator
	a   Node
}

func (n *UnaryExpr) String() string {
	return fmt.Sprintf("(%v %v)", n.op, n.a)
}

func (n *UnaryExpr) Calc(x complex128, vars Vars) (complex128, error) {
	a, err := n.a.Calc(x, vars)
	return -a, err
}

func (n *UnaryExpr) Lua() string {
	return n.op.String() + n.a.Lua()
}

////////////////

type Variable struct {
	pos  int
	name string
}

func (n *Variable) String() string {
	return fmt.Sprintf("'%s'", n.name)
}

func (n *Variable) Calc(x complex128, vars Vars) (complex128, error) {
	if n.name == "x" {
		return x, nil
	} else if val, ok := vars[n.name]; ok {
		return val, nil
	}
	return cmplx.NaN(), ParseErrorf(n.pos, "undeclared variable '%s'", n.name)
}

func (n *Variable) Lua() string {
	return n.name
}

////////////////

type Number struct {
	pos int
	val complex128
}

func (n *Number) String() string {
	if imag(n.val) == 0 {
		return fmt.Sprintf("'%v'", real(n.val))
	}
	return fmt.Sprintf("'%v'", n.val)
}

func (n *Number) Calc(x complex128, vars Vars) (complex128, error) {
	return n.val, nil
}

func (n *Number) Lua() string {
	return fmt.Sprintf("%v", real(n.val))
}
