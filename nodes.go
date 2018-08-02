package formulae

import (
	"fmt"
	"math/cmplx"

	"github.com/tdewolff/formulae/hash"
)

type Node interface {
	String() string
    LaTeX() string
	Calc(complex128, Vars) (complex128, error)
}

func nodeIsGroup(n Node) bool {
    _, ok := n.(*Expr)
    return ok
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

func (n *Func) LaTeX() string {
    if nodeIsGroup(n.a) {
	    return fmt.Sprintf("\\%v(%s)", n.name, n.a.LaTeX())
    }
    return fmt.Sprintf("\\%v %s", n.name, n.a.LaTeX())
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

func (n *Expr) LaTeX() string {
    a := n.a.LaTeX()
    if aExpr, ok := n.a.(*Expr); ok && OpPrec[n.op] > OpPrec[aExpr.op] {
        a = "("+a+")"
    } else if aFunc, ok := n.a.(*Func); ok && OpPrec[n.op] > OpPrec[FuncOp] {
        a = fmt.Sprintf("\\%v(%s)", aFunc.name, aFunc.a.LaTeX())
    }

    b := n.b.LaTeX()
    if n.op == PowerOp {
        b = "{"+b+"}"
    } else if bExpr, ok := n.b.(*Expr); ok && OpPrec[n.op] > OpPrec[bExpr.op] {
        b = "("+b+")"
    } else if _, ok := n.b.(*Func); ok && OpPrec[n.op] > OpPrec[FuncOp] {
        b = "("+b+")"
    }

    if n.op == DivideOp && (len(a) > 1 || len(b) > 1) {
	    return fmt.Sprintf("\\frac{%s}{%s}", a, b)
    }
	return fmt.Sprintf("%s%v%s", a, n.op, b)
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

func (n *UnaryExpr) LaTeX() string {
    if nodeIsGroup(n.a) {
	    return fmt.Sprintf("%v(%s)", n.op, n.a.LaTeX())
    }
	return fmt.Sprintf("%v%s", n.op, n.a.LaTeX())
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

func (n *Variable) LaTeX() string {
	return fmt.Sprintf("%s", n.name)
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
	return "'x'"
}

func (n *Argument) LaTeX() string {
	return "x"
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

func (n *Number) LaTeX() string {
	if imag(n.val) == 0 {
		return fmt.Sprintf("%v", real(n.val))
	}
	return fmt.Sprintf("%v", n.val)
}

func (n *Number) Calc(x complex128, vars Vars) (complex128, error) {
    return n.val, nil
}
