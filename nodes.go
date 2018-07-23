package formulae

import (
	"fmt"
	"math/cmplx"

	"github.com/tdewolff/formulae/hash"
)

type Node interface {
	String() string
	Calc([]complex128, []complex128) ([]complex128, error)
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

func (n *Func) Calc(xs, ys []complex128) ([]complex128, error) {
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

    ys, err := n.a.Calc(xs, ys)
    if err != nil {
        return nil, err
    }

    for i := range ys {
        ys[i] = f(ys[i])
    }
    return ys, nil
}

////////////////

type Expr struct {
	pos int
	op  Operator
	a   Node
	b   Node

    buf *[]complex128
}

func (n *Expr) String() string {
	return fmt.Sprintf("(%v %v %v)", n.a, n.op, n.b)
}

func (n *Expr) Calc(xs, ys []complex128) ([]complex128, error) {
	as, err := n.a.Calc(xs, ys)
	if err != nil {
		return nil, err
	}

    copy(*n.buf, xs)
    ys, *n.buf = *n.buf, ys
	bs, err := n.b.Calc(xs, ys)
	if err != nil {
		return nil, err
	}

	switch n.op {
	case AddOp:
        for i := range ys {
            ys[i] = as[i] + bs[i]
        }
	case SubtractOp:
        for i := range ys {
            ys[i] = as[i] - bs[i]
        }
	case MultiplyOp:
        for i := range ys {
            ys[i] = as[i] * bs[i]
        }
	case DivideOp:
        for i := range ys {
            if bs[i] == 0 {
                ys[i] = cmplx.NaN()
                continue
            }
            ys[i] = as[i] / bs[i]
        }
	case PowerOp:
        for i := range ys {
            ys[i] = cmplx.Pow(as[i], bs[i])
        }
	default:
		return nil, ParseErrorf(n.pos, "unknown operation '%s'", n.op)
	}
    return ys, nil
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

func (n *UnaryExpr) Calc(xs, ys []complex128) ([]complex128, error) {
    var err error
	ys, err = n.a.Calc(xs, ys)
    if err != nil {
        return nil, err
    }

    for i := range ys {
        ys[i] = -ys[i]
    }
    return ys, nil
}

////////////////

type Variable struct {
	pos  int
	name string

    vars *Vars
}

func (n *Variable) String() string {
	return fmt.Sprintf("'%s'", n.name)
}

func (n *Variable) Calc(xs, ys []complex128) ([]complex128, error) {
    if _, ok := (*n.vars)[n.name]; !ok {
        return nil, ParseErrorf(n.pos, "undeclared variable '%s'", n.name)
    }
    val := (*n.vars)[n.name]
    for i := range ys {
        ys[i] = val
    }
    return ys, nil
}

////////////////

type Argument struct {
	pos  int
}

func (n *Argument) String() string {
	return fmt.Sprintf("'x'")
}

func (n *Argument) Calc(xs, ys []complex128) ([]complex128, error) {
    return xs, nil
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

func (n *Number) Calc(xs, ys []complex128) ([]complex128, error) {
    for i := range ys {
        ys[i] = n.val
    }
    return ys, nil
}
