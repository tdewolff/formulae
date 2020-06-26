package formulae

import (
	"fmt"
	"math"
	"math/cmplx"

	"github.com/tdewolff/formulae/hash"
)

type Node interface {
	String() string
	LaTeX() string
	Equal(Node) bool
	Derivative() Node
	Calc(complex128, Vars) (complex128, error)
}

var ZeroNode = &Number{val: 0 + 0i}
var OneNode = &Number{val: 1 + 0i}
var TwoNode = &Number{val: 2 + 0i}
var MinusOneNode = &Number{val: -1 + 0i}

func isNegative(n Node) bool {
	if nNumber, ok := n.(*Number); ok && real(nNumber.val) < 0.0 {
		return true
	} else if nUnaryExpr, ok := n.(*UnaryExpr); ok && nUnaryExpr.op == MinusOp {
		return true
	}
	return false
}

func negateNode(n Node) Node {
	if nNumber, ok := n.(*Number); ok {
		return &Number{val: -nNumber.val}
	} else if nUnaryExpr, ok := n.(*UnaryExpr); ok && nUnaryExpr.op == MinusOp {
		return nUnaryExpr.a
	}
	return &UnaryExpr{op: MinusOp, a: n}
}

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
			} else {
				lNumber, _ := n.l.(*Number)
				rNumber, _ := n.r.(*Number)
				if lNumber != nil && rNumber != nil {
					return &Number{val: lNumber.val + rNumber.val}
				}

				if isNegative(n.r) {
					return &Expr{op: SubtractOp, l: n.l, r: negateNode(n.r)}
				}
			}
		case SubtractOp:
			if n.l.Equal(ZeroNode) {
				return negateNode(n.r)
			} else if n.r.Equal(ZeroNode) {
				return n.l
			} else {
				lNumber, _ := n.l.(*Number)
				rNumber, _ := n.r.(*Number)
				if lNumber != nil && rNumber != nil {
					return &Number{val: lNumber.val - rNumber.val}
				}

				if isNegative(n.r) {
					return &Expr{op: AddOp, l: n.l, r: negateNode(n.r)}
				}
			}
		case MultiplyOp:
			if n.l.Equal(ZeroNode) || n.r.Equal(ZeroNode) {
				return ZeroNode
			} else if n.l.Equal(OneNode) {
				return n.r
			} else if n.r.Equal(OneNode) {
				return n.l
			} else if n.l.Equal(MinusOneNode) {
				return negateNode(n.r)
			} else if n.r.Equal(MinusOneNode) {
				return negateNode(n.l)
			} else {
				lNumber, _ := n.l.(*Number)
				rNumber, _ := n.r.(*Number)
				if lNumber != nil && rNumber != nil {
					return &Number{val: lNumber.val * rNumber.val}
				} else if rNumber != nil {
					n.l, n.r = n.r, n.l
				}

				// move unary minus outwards
				if isNegative(n.l) && isNegative(n.r) {
					return &Expr{
						op: MultiplyOp,
						l:  negateNode(n.l),
						r:  negateNode(n.r),
					}
				} else if isNegative(n.l) {
					return negateNode(&Expr{
						op: MultiplyOp,
						l:  negateNode(n.l),
						r:  n.r,
					})
				} else if isNegative(n.r) {
					return negateNode(&Expr{
						op: MultiplyOp,
						l:  n.l,
						r:  negateNode(n.r),
					})
				}
			}
		case DivideOp:
			if n.r.Equal(OneNode) {
				return n.l
			} else if n.r.Equal(MinusOneNode) {
				return negateNode(n.l)
			} else {
				lNumber, _ := n.l.(*Number)
				rNumber, _ := n.r.(*Number)
				if lNumber != nil && rNumber != nil {
					return &Number{val: lNumber.val / rNumber.val}
				}

				// move unary minus outwards
				if isNegative(n.l) && isNegative(n.r) {
					return &Expr{
						op: DivideOp,
						l:  negateNode(n.l),
						r:  negateNode(n.r),
					}
				} else if isNegative(n.l) {
					return negateNode(&Expr{
						op: DivideOp,
						l:  negateNode(n.l),
						r:  n.r,
					})
				} else if isNegative(n.r) {
					return negateNode(&Expr{
						op: DivideOp,
						l:  n.l,
						r:  negateNode(n.r),
					})
				}
			}
		case PowerOp:
			if n.l.Equal(ZeroNode) {
				return ZeroNode
			} else if n.l.Equal(OneNode) || n.r.Equal(ZeroNode) {
				return OneNode
			} else if n.r.Equal(OneNode) {
				return n.l
			} else if isNegative(n.r) {
				return Optimize(&Expr{
					op: DivideOp,
					l:  OneNode,
					r: &Expr{
						op: PowerOp,
						l:  n.l,
						r:  negateNode(n.r),
					},
				})
			} else {
				lNumber, _ := n.l.(*Number)
				rNumber, _ := n.r.(*Number)
				if lNumber != nil && rNumber != nil {
					return &Number{val: cmplx.Pow(lNumber.val, rNumber.val)}
				} else if lNumber != nil && lNumber.val == 10+0i {
					if rFunc, ok := n.r.(*Func); ok && rFunc.name == hash.Log10 {
						return rFunc.a
					}
				}

				if isNegative(n.l) && rNumber != nil && imag(rNumber.val) == 0.0 && real(rNumber.val) == math.Trunc(real(rNumber.val)) {
					if pow := int64(real(rNumber.val)); pow%2 == 0 {
						return &Expr{
							op: PowerOp,
							l:  negateNode(n.l),
							r:  n.r,
						}
					} else {
						return negateNode(&Expr{
							op: PowerOp,
							l:  negateNode(n.l),
							r:  n.r,
						})
					}
				}

				if lVariable, ok := n.l.(*Variable); ok && lVariable.name == "e" {
					if rFunc, ok := n.r.(*Func); ok && rFunc.name == hash.Log {
						return rFunc.a
					}
				}
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
	case *Func:
		n.a = Optimize(n.a)
		switch n.name {
		case hash.Log:
			if aVariable, ok := n.a.(*Variable); ok && aVariable.name == "e" {
				return OneNode
			}
		case hash.Log10:
			if aNumber, ok := n.a.(*Number); ok && aNumber.val == 10+0i {
				return OneNode
			}
		case hash.Sin, hash.Tan:
			if aNumber, ok := n.a.(*Number); ok && aNumber.val == 0+0i {
				return ZeroNode
			} else if isNegative(n.a) {
				return negateNode(&Func{name: n.name, a: negateNode(n.a)})
			}
		case hash.Cos:
			if aNumber, ok := n.a.(*Number); ok && aNumber.val == 0+0i {
				return OneNode
			} else if isNegative(n.a) {
				return &Func{name: hash.Cos, a: negateNode(n.a)}
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
	return fmt.Sprintf("%v(%v)", n.name, n.a)
}

func (n *Func) LaTeX() string {
	name := n.name.String()
	if n.name == hash.Log {
		name = "log_e"
	} else if n.name == hash.Log {
		name = "log_{10}"
	} else if n.name == hash.Sqrt {
		return fmt.Sprintf("\\sqrt{%s}", n.a.LaTeX())
	}
	return fmt.Sprintf("\\%s\\left(%s\\right)", name, n.a.LaTeX())
}

func (n *Func) Equal(iother Node) bool {
	other, ok := iother.(*Func)
	return ok && n.name == other.name && n.a.Equal(other.a)
}

func (n *Func) Derivative() Node {
	switch n.name {
	case hash.Sin:
		return &Expr{ // cos(a) * da/dx
			op: MultiplyOp,
			l:  &Func{name: hash.Cos, a: n.a},
			r:  n.a.Derivative(),
		}
	case hash.Cos:
		return &Expr{ // -sin(a) * da/dx
			op: MultiplyOp,
			l:  &UnaryExpr{op: MinusOp, a: &Func{name: hash.Sin, a: n.a}},
			r:  n.a.Derivative(),
		}
	case hash.Tan:
		return &Expr{ // 1/cos(a)^2 * da/dx
			op: MultiplyOp,
			l: &Expr{
				op: DivideOp,
				l:  OneNode,
				r: &Expr{
					op: PowerOp,
					l:  &Func{name: hash.Cos, a: n.a},
					r:  TwoNode,
				},
			},
			r: n.a.Derivative(),
		}
	case hash.Arcsin:
		return &Expr{ // 1/sqrt(1-a^2) * da/dx
			op: MultiplyOp,
			l: &Expr{
				op: DivideOp,
				l:  OneNode,
				r: &Func{name: hash.Sqrt, a: &Expr{
					op: SubtractOp,
					l:  OneNode,
					r:  &Expr{op: PowerOp, l: n.a, r: TwoNode},
				}},
			},
			r: n.a.Derivative(),
		}
	case hash.Arccos:
		return &Expr{ // -1/sqrt(1-a^2) * da/dx
			op: MultiplyOp,
			l: &Expr{
				op: DivideOp,
				l:  MinusOneNode,
				r: &Func{name: hash.Sqrt, a: &Expr{
					op: SubtractOp,
					l:  OneNode,
					r:  &Expr{op: PowerOp, l: n.a, r: TwoNode},
				}},
			},
			r: n.a.Derivative(),
		}
	case hash.Arctan:
		return &Expr{ // 1/(1+a^2) * da/dx
			op: MultiplyOp,
			l: &Expr{
				op: DivideOp,
				l:  OneNode,
				r: &Expr{
					op: AddOp,
					l:  OneNode,
					r:  &Expr{op: PowerOp, l: n.a, r: TwoNode},
				},
			},
			r: n.a.Derivative(),
		}
	case hash.Sinh:
		return &Expr{ // cosh(a) * da/dx
			op: MultiplyOp,
			l:  &Func{name: hash.Cosh, a: n.a},
			r:  n.a.Derivative(),
		}
	case hash.Cosh:
		return &Expr{ // sinh(a) * da/dx
			op: MultiplyOp,
			l:  &Func{name: hash.Sinh, a: n.a},
			r:  n.a.Derivative(),
		}
	case hash.Tanh:
		return &Expr{ // 1/cosh(a)^2 * da/dx
			op: MultiplyOp,
			l: &Expr{
				op: DivideOp,
				l:  OneNode,
				r: &Expr{
					op: PowerOp,
					l:  &Func{name: hash.Cosh, a: n.a},
					r:  TwoNode,
				},
			},
			r: n.a.Derivative(),
		}
	case hash.Arcsinh:
		return &Expr{ // 1/sqrt(a^2+1) * da/dx
			op: MultiplyOp,
			l: &Expr{
				op: DivideOp,
				l:  OneNode,
				r: &Func{name: hash.Sqrt, a: &Expr{
					op: AddOp,
					l:  &Expr{op: PowerOp, l: n.a, r: TwoNode},
					r:  OneNode,
				}},
			},
			r: n.a.Derivative(),
		}
	case hash.Arccosh:
		return &Expr{ // 1/sqrt(a^2-1) * da/dx
			op: MultiplyOp,
			l: &Expr{
				op: DivideOp,
				l:  OneNode,
				r: &Func{name: hash.Sqrt, a: &Expr{
					op: SubtractOp,
					l:  &Expr{op: PowerOp, l: n.a, r: TwoNode},
					r:  OneNode,
				}},
			},
			r: n.a.Derivative(),
		}
	case hash.Arctanh:
		return &Expr{ // 1/(1-a^2) * da/dx
			op: MultiplyOp,
			l: &Expr{
				op: DivideOp,
				l:  OneNode,
				r: &Expr{
					op: SubtractOp,
					l:  OneNode,
					r:  &Expr{op: PowerOp, l: n.a, r: TwoNode},
				},
			},
			r: n.a.Derivative(),
		}
	case hash.Sqrt:
		return &Expr{ // 1/(2*sqrt(a)) * da/dx
			op: MultiplyOp,
			l: &Expr{
				op: DivideOp,
				l:  OneNode,
				r: &Expr{
					op: MultiplyOp,
					l:  TwoNode,
					r:  &Func{name: hash.Sqrt, a: n.a},
				},
			},
			r: n.a.Derivative(),
		}
	case hash.Log:
		return &Expr{
			op: MultiplyOp,
			l:  &Expr{op: DivideOp, l: OneNode, r: n.a},
			r:  n.a.Derivative(),
		}
	case hash.Log10:
		return &Expr{ // 1/(a*ln(10)) * da/dx
			op: MultiplyOp,
			l: &Expr{
				op: DivideOp,
				l:  OneNode,
				r: &Expr{
					op: MultiplyOp,
					l:  n.a,
					r:  &Func{name: hash.Log, a: &Number{val: 10 + 0i}},
				},
			},
			r: n.a.Derivative(),
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
	case hash.Sqrt:
		f = cmplx.Sqrt
	case hash.Log:
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
	op Operator
	l  Node
	r  Node
}

func (n *Expr) String() string {
	l := n.l.String()
	if lExpr, ok := n.l.(*Expr); ok && (OpPrec[n.op] > OpPrec[lExpr.op] || OpRightAssoc[n.op] && OpPrec[n.op] == OpPrec[lExpr.op]) {
		l = "(" + l + ")"
	}

	r := n.r.String()
	if rExpr, ok := n.r.(*Expr); ok && (OpPrec[n.op] > OpPrec[rExpr.op] || !OpRightAssoc[n.op] && OpPrec[n.op] == OpPrec[rExpr.op] && n.op != rExpr.op) {
		r = "(" + r + ")"
	}
	return fmt.Sprintf("%s%v%s", l, n.op, r)
}

func (n *Expr) LaTeX() string {
	l := n.l.LaTeX()
	if lExpr, ok := n.l.(*Expr); ok && (OpPrec[n.op] > OpPrec[lExpr.op] || OpRightAssoc[n.op] && OpPrec[n.op] == OpPrec[lExpr.op]) {
		l = "\\left(" + l + "\\right)"
	}

	r := n.r.LaTeX()
	if n.op == PowerOp {
		r = "{" + r + "}"
	} else if rExpr, ok := n.r.(*Expr); ok && (OpPrec[n.op] > OpPrec[rExpr.op] || !OpRightAssoc[n.op] && OpPrec[n.op] == OpPrec[rExpr.op] && n.op != rExpr.op) {
		r = "\\left(" + r + "\\right)"
	}

	if n.op == DivideOp {
		return fmt.Sprintf("\\frac{%s}{%s}", l, r)
	} else if n.op == MultiplyOp {
		// TODO: do more checking for when we have to add the multiplier explicitly
		if _, isNumber := n.r.(*Number); !isNumber {
			return fmt.Sprintf("%s %s", l, r)
		}
	}
	return fmt.Sprintf("%s%v%s", l, n.op, r)
}

func (n *Expr) Equal(iother Node) bool {
	other, ok := iother.(*Expr)
	return ok && n.op == other.op && n.l.Equal(other.l) && n.r.Equal(other.r)
}

func (n *Expr) Derivative() Node {
	switch n.op {
	case AddOp:
		return &Expr{
			op: AddOp,
			l:  n.l.Derivative(),
			r:  n.r.Derivative(),
		}
	case SubtractOp:
		return &Expr{
			op: SubtractOp,
			l:  n.l.Derivative(),
			r:  n.r.Derivative(),
		}
	case MultiplyOp:
		return &Expr{ // r * dl/dx + l * dr/dx
			op: AddOp,
			l: &Expr{
				op: MultiplyOp,
				l:  n.r,
				r:  n.l.Derivative(),
			},
			r: &Expr{
				op: MultiplyOp,
				l:  n.l,
				r:  n.r.Derivative(),
			},
		}
	case DivideOp:
		return &Expr{
			op: DivideOp,
			l: &Expr{ // r * dl/dx - l * dr/dx
				op: SubtractOp,
				l: &Expr{
					op: MultiplyOp,
					l:  n.r,
					r:  n.l.Derivative(),
				},
				r: &Expr{
					op: MultiplyOp,
					l:  n.l,
					r:  n.r.Derivative(),
				},
			},
			r: &Expr{ // r^2
				op: PowerOp,
				l:  n.r,
				r:  TwoNode,
			},
		}
	case PowerOp:
		return &Expr{
			op: AddOp,
			l: &Expr{ // r * l^(r-1) * dl/dx
				op: MultiplyOp,
				l: &Expr{
					op: MultiplyOp,
					l:  n.r,
					r: &Expr{
						op: PowerOp,
						l:  n.l,
						r: &Expr{
							op: SubtractOp,
							l:  n.r,
							r:  OneNode,
						},
					},
				},
				r: n.l.Derivative(),
			},
			r: &Expr{ // l^r * ln(l) * dr/dx
				op: MultiplyOp,
				l: &Expr{
					op: MultiplyOp,
					l:  n,
					r: &Func{
						name: hash.Log,
						a:    n.l,
					},
				},
				r: n.r.Derivative(),
			},
		}
	default:
		panic("unknown operation")
	}
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
	op Operator
	a  Node
}

func (n *UnaryExpr) String() string {
	if nodeIsGroup(n.a) {
		return fmt.Sprintf("%v(%v)", n.op, n.a)
	}
	return fmt.Sprintf("%v%v", n.op, n.a)
}

func (n *UnaryExpr) LaTeX() string {
	if nodeIsGroup(n.a) {
		return fmt.Sprintf("%v\\left(%s\\right)", n.op, n.a.LaTeX())
	}
	return fmt.Sprintf("%v%s", n.op, n.a.LaTeX())
}

func (n *UnaryExpr) Equal(iother Node) bool {
	other, ok := iother.(*UnaryExpr)
	return ok && n.op == other.op && n.a.Equal(other.a)
}

func (n *UnaryExpr) Derivative() Node {
	return &UnaryExpr{op: n.op, a: n.a.Derivative()}
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
	return fmt.Sprintf("%s", n.name)
}

func (n *Variable) LaTeX() string {
	return fmt.Sprintf("%s", n.name)
}

func (n *Variable) Equal(iother Node) bool {
	other, ok := iother.(*Variable)
	return ok && n.name == other.name
}

func (n *Variable) Derivative() Node {
	if n.name == "x" {
		return OneNode
	}
	return ZeroNode
}

func (n *Variable) Calc(x complex128, vars Vars) (complex128, error) {
	if n.name == "x" {
		return x, nil
	} else if _, ok := vars[n.name]; !ok {
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
		return fmt.Sprintf("%v", real(n.val))
	}
	return fmt.Sprintf("%v", n.val)
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
