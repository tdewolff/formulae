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
		return 0, err
	}

	switch n.Name {
	case hash.Sqrt:
		return math.Sqrt(x), nil
	case hash.Sin:
		return math.Sin(x), nil
	default:
		return 0, ParseErrorf(n.Pos, "unknown function '%s'", n.Name)
	}
}

////////////////

type Group struct {
	X Node
}

func (n *Group) String() string {
	return fmt.Sprintf("%v", n.X)
}

func (n *Group) Calc(vars Vars) (float64, error) {
	return n.X.Calc(vars)
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
		return 0, err
	}

	y, err := n.Y.Calc(vars)
	if err != nil {
		return 0, err
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
			return 0, ParseErrorf(n.Pos, "division by zero")
		}
		return x / y, nil
	case PowerOp:
		return math.Pow(x, y), nil
	default:
		return 0, ParseErrorf(n.Pos, "unknown operation '%s'", n.Op)
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
		return 0, err
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
	return 0, ParseErrorf(n.Pos, "undeclared variable '%s'", n.Name)
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
