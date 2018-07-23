package formulae

import (
	"fmt"
	"math/cmplx"

	"github.com/tdewolff/formulae/hash"
)

type Node interface {
	String() string
	Calc(complex128, Vars) (complex128, error)
	CalcN([]complex128, Vars) ([]complex128, error)
}

var tmp []complex128

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
    return f(a), nil
}

func (n *Func) CalcN(xs []complex128, vars Vars) ([]complex128, error) {
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
		return nil, ParseErrorf(n.pos, "unknown function '%s'", n.name)
	}

    xs, err := n.a.CalcN(xs, vars)
    if err != nil {
        return nil, err
    }

    for i := range xs {
        xs[i] = f(xs[i])
    }
    return xs, nil

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
            return cmplx.NaN(), ParseErrorf(n.pos, "division by zero")
        }
        return a / b, nil
	case PowerOp:
        return cmplx.Pow(a, b), nil
	default:
		return cmplx.NaN(), ParseErrorf(n.pos, "unknown operation '%s'", n.op)
	}
}

func (n *Expr) CalcN(xs []complex128, vars Vars) ([]complex128, error) {
	as, err := n.a.CalcN(xs, vars)
	if err != nil {
		return nil, err
	}

    if len(tmp) < len(xs) {
        tmp = make([]complex128, len(xs))
    }
    tmp = tmp[:len(xs)]
    copy(tmp, xs)

	bs, err := n.b.CalcN(tmp, vars)
	if err != nil {
		return nil, err
	}

	switch n.op {
	case AddOp:
        for i := range xs {
            xs[i] = as[i] + bs[i]
        }
	case SubtractOp:
        for i := range xs {
            xs[i] = as[i] - bs[i]
        }
	case MultiplyOp:
        for i := range xs {
            xs[i] = as[i] * bs[i]
        }
	case DivideOp:
        for i := range xs {
            if bs[i] == 0 {
                xs[i] = cmplx.NaN()
                continue
            }
            xs[i] = as[i] / bs[i]
        }
	case PowerOp:
        for i := range xs {
            xs[i] = cmplx.Pow(as[i], bs[i])
        }
	default:
		return nil, ParseErrorf(n.pos, "unknown operation '%s'", n.op)
	}
    return xs, nil
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

func (n *UnaryExpr) CalcN(xs []complex128, vars Vars) ([]complex128, error) {
    var err error
	xs, err = n.a.CalcN(xs, vars)
    if err != nil {
        return nil, err
    }

    for i := range xs {
        xs[i] = -xs[i]
    }
    return xs, nil
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

func (n *Variable) CalcN(xs []complex128, vars Vars) ([]complex128, error) {
	if n.name == "x" {
        return xs, nil
	} else if val, ok := vars[n.name]; ok {
        for i := range xs {
            xs[i] = val
        }
        return xs, nil
	}
	return nil, ParseErrorf(n.pos, "undeclared variable '%s'", n.name)
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

func (n *Number) CalcN(xs []complex128, vars Vars) ([]complex128, error) {
    for i := range xs {
        xs[i] = n.val
    }
    return xs, nil
}
