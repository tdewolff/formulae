package formulae

import (
	"fmt"
	"math/cmplx"

	"github.com/tdewolff/formulae/hash"
)

type Node interface {
	String() string
	Calc(complex128, Vars) (complex128, error)
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
    y, err := n.a.Calc(x, vars)
    if err != nil {
        return cmplx.NaN(), err
    }

    var f func(complex128) complex128
    switch n.name {
    case hash.Sqrt:
        f = cmplx.Sqrt
    case hash.Sin:
        f = cmplx.Sin
    case hash.Cos:
        f = cmplx.Cos
    case hash.Tan:
        f = cmplx.Tan
    case hash.Arcsin:
        f = cmplx.Asin
    case hash.Arccos:
        f = cmplx.Acos
    case hash.Arctan:
        f = cmplx.Atan
    case hash.Sinh:
        f = cmplx.Sinh
    case hash.Cosh:
        f = cmplx.Cosh
    case hash.Tanh:
        f = cmplx.Tanh
    case hash.Arcsinh:
        f = cmplx.Asinh
    case hash.Arccosh:
        f = cmplx.Acosh
    case hash.Arctanh:
        f = cmplx.Atanh
    case hash.Exp:
        f = cmplx.Exp
    case hash.Log, hash.Ln:
        f = cmplx.Log
    case hash.Log10:
        f = cmplx.Log10
    default:
        return cmplx.NaN(), ParseErrorf(n.pos, "unknown function '%s'", n.name)
    }
    return f(y), nil
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

    var y complex128
	switch n.op {
	case AddOp:
        y = a + b
	case SubtractOp:
        y = a - b
	case MultiplyOp:
        y = a * b
	case DivideOp:
        if b == 0 {
            return cmplx.NaN(), ParseErrorf(n.pos, "division by zero")
        }
        y = a / b
	case PowerOp:
        y = cmplx.Pow(a, b)
	default:
		return cmplx.NaN(), ParseErrorf(n.pos, "unknown operation '%s'", n.op)
	}
    return y, nil
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
	y, err := n.a.Calc(x, vars)
    if err != nil {
        return cmplx.NaN(), err
    }
    return -y, nil
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
    if _, ok := vars[n.name]; !ok {
        return cmplx.NaN(), ParseErrorf(n.pos, "undefined variable '%s'", n.name)
    }
    return vars[n.name], nil
}

////////////////

type Argument struct {
	pos  int
}

func (n *Argument) String() string {
	return fmt.Sprintf("'x'")
}

func (n *Argument) Calc(x complex128, vars Vars) (complex128, error) {
    return x, nil
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
