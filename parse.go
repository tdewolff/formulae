package formulae

import (
	"fmt"
	"io"
	"math"
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
	MultiplyOp: 1,
	DivideOp:   1,
	PowerOp:    2,
	MinusOp:    3,
}

var OpRightAssoc = map[Operator]bool{
	PowerOp: true,
	MinusOp: true,
}

type ParseError struct {
	pos int
	err string
}

func (pe *ParseError) Error() string {
	return fmt.Sprintf("%d: %s", pe.pos, pe.err)
}

func ParseErrorf(pos int, format string, args ...interface{}) *ParseError {
	return &ParseError{
		pos,
		fmt.Sprintf(format, args...),
	}
}

type Parser struct {
	output        []SYToken
	operatorStack []SYToken
}

func Parse(in string) (*Formula, []error) {
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
				if n := len(p.operatorStack); n == 0 || p.operatorStack[n-1].op != OpenOp {
					errs = append(errs, ParseErrorf(l.Pos(), "mismatched closing parentheses"))
					break LOOP
				}
			default:
				for len(p.operatorStack) > 0 {
					stack := p.operatorStack[len(p.operatorStack)-1].op
					if !(stack == FuncOp || OpPrec[stack] > OpPrec[op] || !OpRightAssoc[stack] && OpPrec[stack] == OpPrec[op]) || stack == OpenOp {
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

	vars := Vars{}
	for _, tok := range p.output {
		if tok.tt == IdentifierToken {
			vars[string(tok.data)] = math.NaN()
		}
	}

	root, err := p.popNode()
	if err != nil {
		return nil, []error{err}
	}
	if len(p.output) > 0 {
		return nil, []error{fmt.Errorf("some operands remain unparsed")}
	}

	return &Formula{
		root: root,
		vars: vars,
	}, nil
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
		f, err := strconv.ParseFloat(string(tok.data), 64)
		if err != nil {
			return nil, ParseErrorf(tok.pos, "could not parse number: %v", err)
		}
		return &Number{tok.pos, f}, nil
	case IdentifierToken:
		return &Variable{tok.pos, string(tok.data)}, nil
	case OperatorToken:
		switch tok.op {
		case FuncOp:
			n, _ := p.popNode()
			return &Func{tok.pos, tok.function, n}, nil
		case OpenOp:
			n, _ := p.popNode()
			return &Group{n}, nil
		case MinusOp:
			n, _ := p.popNode()
			return &UnaryExpr{tok.pos, tok.op, n}, nil
		default:
			y, err := p.popNode()
			if err != nil && err != ErrNoOperand {
				return nil, err
			}
			x, err := p.popNode()
			if err != nil && err != ErrNoOperand {
				return nil, err
			}
			if x == nil || y == nil {
				return nil, ParseErrorf(tok.pos, "operator has no operands")
			}
			return &Expr{tok.pos, tok.op, x, y}, nil
		}
	default:
		panic("bad token type: " + tok.tt.String())
	}
}
