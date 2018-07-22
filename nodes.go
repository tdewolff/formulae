package formulae

import (
	"fmt"
	"math/cmplx"

	"github.com/tdewolff/formulae/hash"
)

type Calc func(complex128) complex128

type Node interface {
	String() string
	Compile(Vars) (Calc, error)
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

func (n *Func) Compile(vars Vars) (Calc, error) {
	a, err := n.a.Compile(vars)
	if err != nil {
		return nil, err
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
		return nil, ParseErrorf(n.pos, "unknown function '%s'", n.name)
	}
	return func(x complex128) complex128 { return f(a(x)) }, nil
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

func (n *Expr) Compile(vars Vars) (Calc, error) {
	a, err := n.a.Compile(vars)
	if err != nil {
		return nil, err
	}

	b, err := n.b.Compile(vars)
	if err != nil {
		return nil, err
	}

	switch n.op {
	case AddOp:
		return func(x complex128) complex128 { return a(x) + b(x) }, nil
	case SubtractOp:
		return func(x complex128) complex128 { return a(x) - b(x) }, nil
	case MultiplyOp:
		return func(x complex128) complex128 { return a(x) * b(x) }, nil
	case DivideOp:
		return func(x complex128) complex128 {
			bVal := b(x)
			if bVal == 0 {
				return cmplx.Inf()
			}
			return a(x) / b(x)
		}, nil
	case PowerOp:
		return func(x complex128) complex128 { return cmplx.Pow(a(x), b(x)) }, nil
	default:
		return nil, ParseErrorf(n.pos, "unknown operation '%s'", n.op)
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

func (n *UnaryExpr) Compile(vars Vars) (Calc, error) {
	a, err := n.a.Compile(vars)
	if err != nil {
		return nil, err
	}
	return func(x complex128) complex128 { return -a(x) }, nil
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

func (n *Variable) Compile(vars Vars) (Calc, error) {
	if n.name == "x" {
		return func(x complex128) complex128 { return x }, nil
	} else if val, ok := vars[n.name]; ok {
		return func(x complex128) complex128 { return val }, nil
	}
	return nil, ParseErrorf(n.pos, "undeclared variable '%s'", n.name)
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

func (n *Number) Compile(vars Vars) (Calc, error) {
	val := n.val
	return func(x complex128) complex128 { return val }, nil
}

func (n *Number) Lua() string {
	return fmt.Sprintf("%v", real(n.val))
}
