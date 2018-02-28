package lexer

import (
	"io"

	"github.com/relnod/calcgo/token"
)

type stateFn func(*Lexer) token.Token

// Lexer holds the state of the lexer.
type Lexer struct {
	buf BufferedReader
}

// Lex takes an io.Reader and returns a list of tokens.
func Lex(r io.Reader) []token.Token {
	return lexInternal(NewLexer(r))
}

// LexString takes a string as input and returns a list of tokens.
//
// Example:
//  calcgo.Lex("(1 + 2) * 3")
//
// Result:
//  []calcgo.Token{
//    {Value: "",  Type: calcgo.TLeftBracket},
//    {Value: "1", Type: calcgo.TInteger},
//    {Value: "",  Type: calcgo.TOperatorPlus},
//    {Value: "2", Type: calcgo.TInteger},
//    {Value: "",  Type: calcgo.TRightBracket},
//    {Value: "",  Type: calcgo.TOperatorMult},
//    {Value: "2", Type: calcgo.TInteger},
//  })
func LexString(str string) []token.Token {
	return lexInternal(NewLexerFromString(str))
}

func lexInternal(l *Lexer) []token.Token {
	tokens := make([]token.Token, 0)

	for {
		t := l.Read()
		if t.Type == token.EOF {
			break
		}
		tokens = append(tokens, t)
	}

	return tokens
}

// NewLexer returns a new lexer object.
func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		buf: NewBufferedReader(r),
	}
}

// NewLexerFromString returns a new lexer object.
func NewLexerFromString(str string) *Lexer {
	return &Lexer{
		buf: NewBufferedReaderFromString(str),
	}
}

// NewLexerFromBufferedReader returns a new lexer object.
func NewLexerFromBufferedReader(r BufferedReader) *Lexer {
	return &Lexer{
		buf: r,
	}
}

// Read returns the next token.
func (l *Lexer) Read() token.Token {
	token := lexAll(l)

	return token
}

// createToken takes a tokentype and a value to create a token, which it then
// emits.
func (l *Lexer) createToken(tokenType token.Type, value string) token.Token {
	return token.Token{
		Type:  tokenType,
		Value: value,
		Start: l.buf.StartPos(),
		End:   l.buf.CurrPos(),
	}
}

// create emits a new token with type tokenType and the currently stored value.
func (l *Lexer) create(tokenType token.Type) token.Token {
	return l.createToken(tokenType, string(l.buf.All()))
}

// createEmpty emits a new token with type tokenType and an empty value.
func (l *Lexer) createEmpty(tokenType token.Type) token.Token {
	return l.createToken(tokenType, "")
}

// emitEmpty emits a new token with type tokenType and the current character.
func (l *Lexer) createSingle(tokenType token.Type) token.Token {
	return l.createToken(tokenType, string(l.buf.Current()))
}

// lexAll is the entry state of the lexer state machine and also for all tokens.
//
// Transitions:
//  - [0-9] -> lexNumber
//  - [a-z] -> lexVariableOrFunction
//  - rest  -> lexAll
func lexAll(l *Lexer) token.Token {
	var tokenType token.Type

	l.buf.Reset()

	b, ok := l.buf.Next()
	if !ok {
		// todo: return something else here. Maybe forward ok?
		return token.Token{Type: token.EOF}
	}
	if isDigit(b) {
		return lexNumber(l)
	}
	if isLetter(b) {
		return lexVariableOrFunction(l)
	}
	if isWhiteSpace(b) {
		return lexAll(l)
	}

	switch b {
	case '+':
		tokenType = token.Plus
	case '-':
		if b, ok := l.buf.Next(); ok && isDigit(b) {
			return lexNumber(l)
		}
		tokenType = token.Minus
	case '*':
		tokenType = token.Mult
	case '/':
		tokenType = token.Div
	case '%':
		tokenType = token.Mod
	case '|':
		tokenType = token.Or
	case '^':
		tokenType = token.Xor
	case '&':
		tokenType = token.And
	case '(':
		tokenType = token.ParenL
	case ')':
		tokenType = token.ParenR
	default:
		return l.create(token.InvalidCharacter)
	}

	return l.createEmpty(tokenType)
}

// lexNumber is the entry state for all number tokens.
//
// Transitions:
//  - 0x       -> lexHex
//  - 0b       -> lexBin
//  - [0-9]+\. -> lexDecimal
//  - [0-9]+\^ -> lexExponential
//  - rest     -> lexAll
func lexNumber(l *Lexer) token.Token {
	if l.buf.Current() == '0' {
		b, ok := l.buf.Next()
		if ok {
			if b == 'x' {
				return lexHex(l)
			}
			if b == 'b' {
				return lexBin(l)
			}
		}
		l.buf.Backup()
	}

	for {
		b, ok := l.buf.Next()
		if !ok {
			break
		}

		if isDigit(b) {
			continue
		}

		if b == '.' {
			return lexDecimal(l)
		}

		if b == '^' {
			return lexExponential(l)
		}

		if isWhiteSpace(b) || b == ')' {
			l.buf.Backup()
			break
		}

		return l.createSingle(token.InvalidCharacterInNumber)
	}

	return l.create(token.Int)
}

// lexDecimal creates a decimal number token.
//
// Transitions:
//  -> lexAll
func lexDecimal(l *Lexer) token.Token {
	for {
		b, ok := l.buf.Next()
		if !ok {
			break
		}

		if isDigit(b) {
			continue
		}

		if isWhiteSpace(b) || b == ')' {
			l.buf.Backup()
			break
		}

		return l.createSingle(token.InvalidCharacterInNumber)
	}

	return l.create(token.Dec)
}

// lexHex creates a hex number token.
//
// Transitions:
//  -> lexAll
func lexHex(l *Lexer) token.Token {
	for {
		b, ok := l.buf.Next()
		if !ok {
			break
		}

		if isHexDigit(b) {
			continue
		}

		if isWhiteSpace(b) || b == ')' {
			l.buf.Backup()
			break
		}

		return l.createSingle(token.InvalidCharacterInNumber)
	}

	return l.create(token.Hex)
}

// lexBin creates a binary number token.
//
// Transitions:
//  -> lexAll
func lexBin(l *Lexer) token.Token {
	for {
		b, ok := l.buf.Next()
		if !ok {
			break
		}

		if isBinDigit(b) {
			continue
		}

		if isWhiteSpace(b) || b == ')' {
			l.buf.Backup()
			break
		}

		return l.createSingle(token.InvalidCharacterInNumber)
	}

	return l.create(token.Bin)
}

// lexExponential creates an exponential number token.
//
// Transitions:
//  -> lexAll
func lexExponential(l *Lexer) token.Token {
	for {
		b, ok := l.buf.Next()
		if !ok {
			break
		}

		if isDigit(b) {
			continue
		}

		if isWhiteSpace(b) || b == ')' {
			l.buf.Backup()
			break
		}

		return l.createSingle(token.InvalidCharacterInNumber)
	}

	return l.create(token.Exp)
}

// lexVariableOrFunction creates a variable or function token.
//
// Transitions:
//  -> lexAll
func lexVariableOrFunction(l *Lexer) token.Token {
	for {
		b, ok := l.buf.Next()
		if !ok {
			break
		}

		if isLetter(b) {
			continue
		}

		if isWhiteSpace(b) || b == ')' {
			l.buf.Backup()
			break
		}

		if b == '(' {
			switch string(l.buf.All()) {
			case "sqrt(":
				return l.createEmpty(token.Sqrt)
			case "sin(":
				return l.createEmpty(token.Sin)
			case "cos(":
				return l.createEmpty(token.Cos)
			case "tan(":
				return l.createEmpty(token.Tan)
			default:
				return l.create(token.UnkownFunktion)
			}
		}

		return l.createSingle(token.InvalidCharacterInVariable)
	}

	return l.create(token.Var)
}

// isWhiteSpace checks if b is a whitespace character.
func isWhiteSpace(b byte) bool {
	return b == ' '
}

// isDigit checks if b is a digit.
func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

// isHexDigit checks if b is a hexadecimal digit.
func isHexDigit(b byte) bool {
	return isDigit(b) || b >= 'A' && b <= 'F'
}

// isBinDigit checks if b is 0 or 1.
func isBinDigit(b byte) bool {
	return b == '0' || b == '1'
}

// isLetter checks if a is a letter.
func isLetter(b byte) bool {
	return b >= 'a' && b <= 'z'
}
