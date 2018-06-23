package formulae // import "github.com/tdewolff/formulae"

import (
	"io"
	"strconv"
	"unicode"

	"github.com/tdewolff/formulae/hash"
	"github.com/tdewolff/parse"
	"github.com/tdewolff/parse/buffer"
)

var identifierStart = []*unicode.RangeTable{unicode.Lu, unicode.Ll, unicode.Lt, unicode.Lm, unicode.Lo, unicode.Nl, unicode.Other_ID_Start}
var identifierContinue = []*unicode.RangeTable{unicode.Lu, unicode.Ll, unicode.Lt, unicode.Lm, unicode.Lo, unicode.Nl, unicode.Mn, unicode.Mc, unicode.Nd, unicode.Pc, unicode.Other_ID_Continue}

type Operator int

const (
	UnknownOp Operator = iota
	FuncOp
	OpenOp
	CloseOp
	AddOp
	SubtractOp
	MinusOp
	MultiplyOp
	DivideOp
	PowerOp
)

func (op Operator) String() string {
	switch op {
	case FuncOp:
		return "func"
	case OpenOp:
		return "("
	case CloseOp:
		return ")"
	case AddOp:
		return "+"
	case SubtractOp:
		return "-"
	case MinusOp:
		return "-"
	case MultiplyOp:
		return "*"
	case DivideOp:
		return "/"
	case PowerOp:
		return "^"
	}
	return "Invalid(" + strconv.Itoa(int(op)) + ")"
}

////////////////////////////////////////////////////////////////

// TokenType determines the type of token, eg. a number or a semicolon.
type TokenType uint32

// TokenType values.
const (
	ErrorToken TokenType = iota // extra token when errors occur
	UnknownToken
	WhitespaceToken // space \t \v \f \r \n
	NumericToken
	IdentifierToken
	OperatorToken
)

// String returns the string representation of a TokenType.
func (tt TokenType) String() string {
	switch tt {
	case ErrorToken:
		return "Error"
	case UnknownToken:
		return "Unknown"
	case WhitespaceToken:
		return "Whitespace"
	case NumericToken:
		return "Numeric"
	case IdentifierToken:
		return "Identifier"
	case OperatorToken:
		return "Operator"
	}
	return "Invalid(" + strconv.Itoa(int(tt)) + ")"
}

////////////////////////////////////////////////////////////////

// Lexer is the state for the lexer.
type Lexer struct {
	r        *buffer.Lexer
	lastTT   TokenType
	lastOp   Operator
	lastFunc hash.Hash
}

// NewLexer returns a new Lexer for a given io.Reader.
func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		r: buffer.NewLexer(r),
	}
}

// Err returns the error encountered during lexing, this is often io.EOF but also other errors can be returned.
func (l *Lexer) Err() error {
	return l.r.Err()
}

// Restore restores the NULL byte at the end of the buffer.
func (l *Lexer) Restore() {
	l.r.Restore()
}

func (l *Lexer) Pos() int {
	return l.r.Offset()
}

// Next returns the next Token. It returns ErrorToken when an error was encountered. Using Err() one can retrieve the error message.
func (l *Lexer) Next() (TokenType, []byte) {
	// Add in extra multiplier
	isNumeric := l.isNumeric()
	isIdentifier, nIdent := l.isIdentifierStart()
	if isNumeric || isIdentifier || l.r.Peek(0) == '(' {
		if l.lastTT == NumericToken || l.lastTT == IdentifierToken || l.lastTT == OperatorToken && l.lastOp == CloseOp {
			l.lastTT = OperatorToken
			l.lastOp = MultiplyOp
			return OperatorToken, []byte("*")
		}
	}

	var tt TokenType
	if isNumeric {
		l.r.Move(1)
		l.consumeNumericToken()
		tt = NumericToken
	} else if isIdentifier {
		l.r.Move(nIdent)
		tt = l.consumeIdentifierToken()
	} else if l.consumeOperatorToken() {
		tt = OperatorToken
	} else if l.consumeWhitespace() {
		for l.consumeWhitespace() {
		}
		tt = WhitespaceToken
	} else if l.Err() != nil {
		return ErrorToken, nil
	} else {
		l.r.Move(1)
		tt = UnknownToken
	}
	l.lastTT = tt
	return tt, l.r.Shift()
}

func (l *Lexer) Operator() Operator {
	return l.lastOp
}

func (l *Lexer) Function() hash.Hash {
	return l.lastFunc
}

////////////////////////////////////////////////////////////////

func (l *Lexer) isNumeric() bool {
	c := l.r.Peek(0)
	return c >= '0' && c <= '9' || c == '.'
}

func (l *Lexer) isIdentifierStart() (bool, int) {
	c := l.r.Peek(0)
	if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
		return true, 1
	}
	if r, n := l.r.PeekRune(0); unicode.IsOneOf(identifierStart, r) {
		return true, n
	}
	return false, 0
}

func (l *Lexer) consumeWhitespace() bool {
	c := l.r.Peek(0)
	if c == ' ' || c == '\t' {
		l.r.Move(1)
		return true
	} else if c >= 0xC0 {
		if r, n := l.r.PeekRune(0); r == '\uFEFF' || unicode.Is(unicode.Zs, r) {
			l.r.Move(n)
			return true
		}
	}
	return false
}

func (l *Lexer) consumeIdentifierToken() TokenType {
	// Already on second identifier character
	for {
		c := l.r.Peek(0)
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
			l.r.Move(1)
		} else if c >= 0xC0 {
			if r, n := l.r.PeekRune(0); r == '\u200C' || r == '\u200D' || unicode.IsOneOf(identifierContinue, r) {
				l.r.Move(n)
			} else {
				break
			}
		} else {
			break
		}
	}

	ident := parse.ToLower(parse.Copy(l.r.Lexeme()))
	h := hash.ToHash(ident)
	if h != 0 {
		l.lastOp = FuncOp
		l.lastFunc = h
		return OperatorToken
	}
	if _, ok := DefaultVars[string(ident)]; ok {
		parse.ToLower(l.r.Lexeme())
	}
	return IdentifierToken
}

func (l *Lexer) consumeOperatorToken() bool {
	op := UnknownOp
	c := l.r.Peek(0)
	switch c {
	case '(':
		op = OpenOp
	case ')':
		op = CloseOp
	case '+':
		op = AddOp
	case '-':
		if l.lastTT == ErrorToken || l.lastTT == OperatorToken && l.lastOp != CloseOp {
			op = MinusOp
		} else {
			op = SubtractOp
		}
	case '*':
		op = MultiplyOp
	case '/':
		op = DivideOp
	case '^':
		op = PowerOp
	}

	if op == UnknownOp {
		return false
	}

	l.r.Move(1)
	l.lastOp = op
	return true
}

func (l *Lexer) consumeDigit() bool {
	if c := l.r.Peek(0); c >= '0' && c <= '9' {
		l.r.Move(1)
		return true
	}
	return false
}

func (l *Lexer) consumeNumericToken() {
	// Already on second numeric character
	for l.consumeDigit() {
	}
	if l.r.Peek(0) == '.' {
		l.r.Move(1)
		for l.consumeDigit() {
		}
	}
	mark := l.r.Pos()
	c := l.r.Peek(0)
	if c == 'e' || c == 'E' {
		l.r.Move(1)
		c = l.r.Peek(0)
		if c == '+' || c == '-' {
			l.r.Move(1)
		}
		if !l.consumeDigit() {
			// e could belong to the next token
			l.r.Rewind(mark)
			return
		}
		for l.consumeDigit() {
		}
	}
}
