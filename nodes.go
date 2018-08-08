package formulae

import (
	"fmt"
	"math/cmplx"

	"github.com/tdewolff/formulae/hash"
)

type Node interface {
	String() string
    LaTeX() string
    Equal(Node) bool
	Calc(complex128, Vars) (complex128, error)
}

var ZeroNode = &Number{val: 0+0i}
var OneNode = &Number{val: 1+0i}
var MinusOneNode = &Number{val: -1+0i}

func Optimize(in Node) Node {
    switch n := in.(type) {
    case *Expr:
        n.l = Optimize(n.l)
        n.r = Optimize(n.r)
        switch n.op {
        case AddOp:
            if n.l.Equal(ZeroNode) {
                return n.r
            } else if n.r.Equal(ZeroNode) {
                return n.l
            }
        case SubtractOp:
            if n.l.Equal(ZeroNode) {
                return Optimize(&UnaryExpr{op: MinusOp, a: n.r})
            } else if n.r.Equal(ZeroNode) {
                return n.l
            } else if rNumber, ok := n.r.(*Number); ok && real(rNumber.val) < 0.0 {
                return Optimize(&Expr{op: AddOp, l: n.l, r: &Number{val: -rNumber.val}})
            }
        case MultiplyOp:
            if n.l.Equal(ZeroNode) || n.r.Equal(ZeroNode) {
                return ZeroNode
            } else if n.l.Equal(OneNode) {
                return n.r
            } else if n.r.Equal(OneNode) {
                return n.l
            } else if n.l.Equal(MinusOneNode) {
                return Optimize(&UnaryExpr{op: MinusOp, a: n.r})
            } else if n.r.Equal(MinusOneNode) {
                return Optimize(&UnaryExpr{op: MinusOp, a: n.l})
            }
        case DivideOp:
            if n.r.Equal(OneNode) {
                return n.l
            } else if n.r.Equal(MinusOneNode) {
                return Optimize(&UnaryExpr{op: MinusOp, a: n.l})
            }
        case PowerOp:
            if n.l.Equal(ZeroNode) {
                return ZeroNode
            } else if n.l.Equal(OneNode) || n.r.Equal(ZeroNode) {
                return OneNode
            } else if n.r.Equal(OneNode) {
                return n.l
            } else if n.r.Equal(MinusOneNode) {
                return Optimize(&Expr{op: DivideOp, l: OneNode, r: n.l})
            }
        }
    case *UnaryExpr:
        n.a = Optimize(n.a)
        if n.op == MinusOp {
            if aNumber, ok := n.a.(*Number); ok {
                return &Number{val: -aNumber.val}
            } else if aUnaryExpr, ok := n.a.(*UnaryExpr); ok && aUnaryExpr.op == MinusOp {
                return aUnaryExpr.a
            }
        }
    }
    return in
}

func nodeIsGroup(n Node) bool {
    _, ok := n.(*Expr)
    return ok
}

////////////////

type Func struct {
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

func (n *Func) Equal(iother Node) bool {
    other, ok := iother.(*Func)
    return ok && n.name == other.name && n.a.Equal(other.a)
}

func (n *Func) Derivative() Node {
    switch n.name {
    case hash.Sin:
        return &Expr{
            op: MultiplyOp,
            l: &Func{name: hash.Cos, a: n.a},
            r: n.Derivative(),
        }
    default:
        panic("unknown function")
    }
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
        return cmplx.NaN(), fmt.Errorf("unknown function '%s'", n.name)
    }
    return f(y), nil
}

////////////////

type Expr struct {
	op  Operator
	l   Node
	r   Node
}

func (n *Expr) String() string {
	return fmt.Sprintf("(%v %v %v)", n.l, n.op, n.r)
}

func (n *Expr) LaTeX() string {
    l := n.l.LaTeX()
    if lExpr, ok := n.l.(*Expr); ok && OpPrec[n.op] > OpPrec[lExpr.op] {
        l = "("+l+")"
    } else if lFunc, ok := n.l.(*Func); ok && OpPrec[n.op] > OpPrec[FuncOp] {
        l = fmt.Sprintf("\\%v(%s)", lFunc.name, lFunc.a.LaTeX())
    }

    r := n.r.LaTeX()
    if n.op == PowerOp {
        r = "{"+r+"}"
    } else if rExpr, ok := n.r.(*Expr); ok && OpPrec[n.op] > OpPrec[rExpr.op] {
        r = "("+r+")"
    } else if _, ok := n.r.(*Func); ok && OpPrec[n.op] > OpPrec[FuncOp] {
        r = "("+r+")"
    }

    if n.op == DivideOp && (len(l) > 1 || len(r) > 1) {
	    return fmt.Sprintf("\\frac{%s}{%s}", l, r)
    }
	return fmt.Sprintf("%s%v%s", l, n.op, r)
}

func (n *Expr) Equal(iother Node) bool {
    other, ok := iother.(*Expr)
    return ok && n.op == other.op && n.l.Equal(other.l) && n.r.Equal(other.r)
}

func (n *Expr) Derivative() Node {
    panic("unimplemented")
}

func (n *Expr) Calc(x complex128, vars Vars) (complex128, error) {
	l, err := n.l.Calc(x, vars)
	if err != nil {
		return cmplx.NaN(), err
	}

	r, err := n.r.Calc(x, vars)
	if err != nil {
		return cmplx.NaN(), err
	}

    var y complex128
	switch n.op {
	case AddOp:
        y = l + r
	case SubtractOp:
        y = l - r
	case MultiplyOp:
        y = l * r
	case DivideOp:
        if r == 0 {
            return cmplx.NaN(), fmt.Errorf("division by zero")
        }
        y = l / r
	case PowerOp:
        y = cmplx.Pow(l, r)
	default:
		return cmplx.NaN(), fmt.Errorf("unknown operation '%s'", n.op)
	}
    return y, nil
}

////////////////

type UnaryExpr struct {
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

func (n *UnaryExpr) Equal(iother Node) bool {
    other, ok := iother.(*UnaryExpr)
    return ok && n.op == other.op && n.a.Equal(other.a)
}

func (n *UnaryExpr) Derivative() Node {
    panic("unimplemented")
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
	name string
}

func (n *Variable) String() string {
	return fmt.Sprintf("'%s'", n.name)
}

func (n *Variable) LaTeX() string {
	return fmt.Sprintf("%s", n.name)
}

func (n *Variable) Equal(iother Node) bool {
    other, ok := iother.(*Variable)
    return ok && n.name == other.name
}

func (n *Variable) Derivative() Node {
    return ZeroNode
}

func (n *Variable) Calc(x complex128, vars Vars) (complex128, error) {
    if _, ok := vars[n.name]; !ok {
        return cmplx.NaN(), fmt.Errorf("undefined variable '%s'", n.name)
    }
    return vars[n.name], nil
}

////////////////

type Number struct {
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

func (n *Number) Equal(iother Node) bool {
    other, ok := iother.(*Number)
    return ok && n.val == other.val
}

func (n *Number) Derivative() Node {
    return ZeroNode
}

func (n *Number) Calc(x complex128, vars Vars) (complex128, error) {
    return n.val, nil
}

////////////////

type Argument struct {
}

func (n *Argument) String() string {
	return "'x'"
}

func (n *Argument) LaTeX() string {
	return "x"
}

func (n *Argument) Equal(iother Node) bool {
    _, ok := iother.(*Argument)
    return ok
}

func (n *Argument) Derivative() Node {
    return OneNode
}

func (n *Argument) Calc(x complex128, vars Vars) (complex128, error) {
    return x, nil
}
