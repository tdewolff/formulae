package formulae

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/tdewolff/formulae/hash"
)

type SYToken struct {
	pos      int
	tt       TokenType
	data     []byte
	op       Operator
	function hash.Hash
}

func (t SYToken) String() string {
	if t.tt == OperatorToken {
		return fmt.Sprintf("%v", t.op)
	}
	return fmt.Sprintf("'%s'", t.data)
}

var OpPrec = map[Operator]int{
	FuncOp:     1,
	MultiplyOp: 2,
	DivideOp:   2,
	PowerOp:    3,
	MinusOp:    4,
}

var OpRightAssoc = map[Operator]bool{
	PowerOp: true,
	MinusOp: true,
	FuncOp:  true,
}

type ParseError struct {
	pos int
	msg string
}

func (pe ParseError) Pos() int {
	return pe.pos
}

func (pe ParseError) Error() string {
	return pe.msg
}

func ParseErrorf(pos int, format string, args ...interface{}) ParseError {
	return ParseError{
		pos,
		fmt.Sprintf(format, args...),
	}
}

type Parser struct {
	output        []SYToken
	operatorStack []SYToken
}

func Parse(in string) (*Function, []error) {
	var errs []error
	l := NewLexer(strings.NewReader(in))
	p := Parser{}
LOOP:
	for {
		tt, data := l.Next()
		if tt == WhitespaceToken {
			tt, data = l.Next()
		}
		switch tt {
		case ErrorToken:
			if l.Err() != io.EOF {
				errs = append(errs, l.Err())
			}
			break LOOP
		case UnknownToken:
			errs = append(errs, ParseErrorf(l.Pos(), "bad input"))
			break LOOP
		case NumericToken:
			p.output = append(p.output, SYToken{l.Pos(), tt, data, 0, 0})
		case IdentifierToken:
			p.output = append(p.output, SYToken{l.Pos(), tt, data, 0, 0})
		case OperatorToken:
			op := l.Operator()
			sytoken := SYToken{l.Pos(), tt, data, op, 0}
			switch op {
			case FuncOp:
				sytoken.function = l.Function()
				p.operatorStack = append(p.operatorStack, sytoken)
			case OpenOp:
				p.operatorStack = append(p.operatorStack, sytoken)
			case CloseOp:
				for len(p.operatorStack) > 0 && p.operatorStack[len(p.operatorStack)-1].op != OpenOp {
					p.popOperation()
				}
				n := len(p.operatorStack)
				if n == 0 || p.operatorStack[n-1].op != OpenOp {
					errs = append(errs, ParseErrorf(l.Pos(), "mismatched closing parentheses"))
					break LOOP
				} else if n > 1 && p.operatorStack[n-2].op == FuncOp {
					p.popOperation()
				}
				p.popOperation() // pop OpenOp
			default:
				for n := len(p.operatorStack); n > 0; n-- {
					stack := p.operatorStack[n-1].op
					if !(OpPrec[stack] > OpPrec[op] || !OpRightAssoc[stack] && OpPrec[stack] == OpPrec[op]) || stack == OpenOp {
						break
					}
					p.popOperation()
				}
				p.operatorStack = append(p.operatorStack, sytoken)
			}
		default:
			panic("bad token type: " + tt.String())
		}
	}
	for len(p.operatorStack) > 0 {
		p.popOperation()
	}
	if len(errs) > 0 {
		return nil, errs
	}
	if len(p.output) == 0 {
		return nil, []error{fmt.Errorf("empty formula")}
	}

	root, err := p.popNode()
	if err != nil {
		return nil, []error{err}
	}
	if len(p.output) > 0 {
		return nil, []error{fmt.Errorf("some operands remain unparsed")}
	}

	vars := DefaultVars.Duplicate()
	return &Function{root: root, Vars: vars}, nil
}

func (p *Parser) popOperation() {
	p.output = append(p.output, p.operatorStack[len(p.operatorStack)-1])
	p.operatorStack = p.operatorStack[:len(p.operatorStack)-1]
}

var ErrNoOperand = fmt.Errorf("no operand")

func (p *Parser) popNode() (Node, error) {
	if len(p.output) == 0 {
		return nil, ErrNoOperand
	}

	tok := p.output[len(p.output)-1]
	p.output = p.output[:len(p.output)-1]

	switch tok.tt {
	case NumericToken:
		hasReal := true
		hasImag := tok.data[len(tok.data)-1] == 'i'
		iPlus := len(tok.data)
		if hasImag {
			hasReal = false
			iPlus = -1
			for i := 0; i < len(tok.data); i++ {
				if tok.data[i] == '+' {
					hasReal = true
					iPlus = i
					break
				} else if tok.data[i] == 'e' || tok.data[i] == 'E' {
					i++
				}
			}
		}

		var err error
		fr, fi := 0.0, 0.0
		if hasReal {
			fr, err = strconv.ParseFloat(string(tok.data[:iPlus]), 64)
			if err != nil {
				return nil, ParseErrorf(tok.pos, "could not parse number: %v", err)
			}
		}
		if hasImag {
			if len(tok.data) == 1 {
				fi = 1.0
			} else {
				fi, err = strconv.ParseFloat(string(tok.data[iPlus+1:len(tok.data)-1]), 64)
				if err != nil {
					return nil, ParseErrorf(tok.pos, "could not parse number: %v", err)
				}
			}
		}
		return &Number{val: complex(fr, fi)}, nil
	case IdentifierToken:
		if len(tok.data) == 1 && tok.data[0] == 'x' {
			return &Argument{}, nil
		} else {
			return &Variable{name: string(tok.data)}, nil
		}
	case OperatorToken:
		switch tok.op {
		case FuncOp:
			a, _ := p.popNode()
			if tok.function == hash.Exp {
				return &Expr{op: PowerOp, l: &Variable{name: "e"}, r: a}, nil
			} else if tok.function == hash.Ln {
				tok.function = hash.Log
			}
			return &Func{name: tok.function, a: a}, nil
		case OpenOp:
			return p.popNode()
		case MinusOp:
			a, _ := p.popNode()
			return &UnaryExpr{op: tok.op, a: a}, nil
		default:
			r, err := p.popNode()
			if err != nil && err != ErrNoOperand {
				return nil, err
			}
			l, err := p.popNode()
			if err != nil && err != ErrNoOperand {
				return nil, err
			}
			if l == nil || r == nil {
				return nil, ParseErrorf(tok.pos, "operator has no operands")
			}
			return &Expr{op: tok.op, l: l, r: r}, nil
		}
	default:
		panic("bad token type '" + tok.tt.String() + "'")
	}
}
