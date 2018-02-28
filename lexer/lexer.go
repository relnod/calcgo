package lexer

import (
	"io"

	"github.com/relnod/calcgo/token"
)

type stateFn func(*Lexer) stateFn

// Lexer holds the state of the lexer.
type Lexer struct {
	token chan token.Token
	buf   BufferedReader
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

	l.Start()

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
		token: make(chan token.Token, 0),
		buf:   NewBufferedReader(r),
	}
}

// NewLexerFromString returns a new lexer object.
func NewLexerFromString(str string) *Lexer {
	return &Lexer{
		token: make(chan token.Token, len(str)/3),
		buf:   NewBufferedReaderFromString(str),
	}
}

// NewLexerFromBufferedReader returns a new lexer object.
func NewLexerFromBufferedReader(r BufferedReader) *Lexer {
	return &Lexer{
		token: make(chan token.Token, 0),
		buf:   r,
	}
}

// Start runs the lexer in a go routine.
func (l *Lexer) Start() {
	go l.run()
}

// Read returns the next token from the token chanel.
func (l *Lexer) Read() token.Token {
	return <-l.token
}

// emitInternal takes a tokentype and a value to create a token, which it then
// emits.
func (l *Lexer) emitInternal(tokenType token.Type, value string) {
	l.token <- token.Token{
		Type:  tokenType,
		Value: value,
		Start: l.buf.StartPos(),
		End:   l.buf.CurrPos(),
	}
}

// emit emits a new token with type tokenType and the currently stored value.
func (l *Lexer) emit(tokenType token.Type) {
	l.emitInternal(tokenType, string(l.buf.All()))
}

// emitEmpty emits a new token with type tokenType and an empty value.
func (l *Lexer) emitEmpty(tokenType token.Type) {
	l.emitInternal(tokenType, "")
}

// emitEmpty emits a new token with type tokenType and the current character.
func (l *Lexer) emitSingle(tokenType token.Type) {
	l.emitInternal(tokenType, string(l.buf.Current()))
}

// run runs the lexer state machine.
func (l *Lexer) run() {
	for state := lexAll; state != nil; {
		state = state(l)
	}
	close(l.token)
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

// lexAll is the entry state of the lexer state machine and also for all tokens.
//
// Transitions:
//  - [0-9] -> lexNumber
//  - [a-z] -> lexVariableOrFunction
//  - rest  -> lexAll
func lexAll(l *Lexer) stateFn {
	var tokenType token.Type

	l.buf.Reset()

	b, ok := l.buf.Next()
	if !ok {
		return nil
	}
	if isDigit(b) {
		return lexNumber
	}
	if isLetter(b) {
		return lexVariableOrFunction
	}
	if isWhiteSpace(b) {
		return lexAll
	}

	switch b {
	case '+':
		tokenType = token.Plus
	case '-':
		if b, ok := l.buf.Next(); ok && isDigit(b) {
			return lexNumber
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
		l.emit(token.InvalidCharacter)
		return lexAll
	}

	l.emitEmpty(tokenType)
	return lexAll
}

// lexNumber is the entry state for all number tokens.
//
// Transitions:
//  - 0x       -> lexHex
//  - 0b       -> lexBin
//  - [0-9]+\. -> lexDecimal
//  - [0-9]+\^ -> lexExponential
//  - rest     -> lexAll
func lexNumber(l *Lexer) stateFn {
	if l.buf.Current() == '0' {
		b, ok := l.buf.Next()
		if ok {
			if b == 'x' {
				return lexHex
			}
			if b == 'b' {
				return lexBin
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
			return lexDecimal
		}

		if b == '^' {
			return lexExponential
		}

		if isWhiteSpace(b) || b == ')' {
			l.buf.Backup()
			break
		}

		l.emitSingle(token.InvalidCharacterInNumber)
		return lexAll
	}

	l.emit(token.Int)
	return lexAll
}

// lexDecimal creates a decimal number token.
//
// Transitions:
//  -> lexAll
func lexDecimal(l *Lexer) stateFn {
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

		l.emitSingle(token.InvalidCharacterInNumber)
		return lexAll
	}

	l.emit(token.Dec)
	return lexAll
}

// lexHex creates a hex number token.
//
// Transitions:
//  -> lexAll
func lexHex(l *Lexer) stateFn {
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

		l.emitSingle(token.InvalidCharacterInNumber)
		return lexAll
	}

	l.emit(token.Hex)
	return lexAll
}

// lexBin creates a binary number token.
//
// Transitions:
//  -> lexAll
func lexBin(l *Lexer) stateFn {
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

		l.emitSingle(token.InvalidCharacterInNumber)
		return lexAll
	}

	l.emit(token.Bin)
	return lexAll
}

// lexExponential creates an exponential number token.
//
// Transitions:
//  -> lexAll
func lexExponential(l *Lexer) stateFn {
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

		l.emitSingle(token.InvalidCharacterInNumber)
		return lexAll
	}

	l.emit(token.Exp)
	return lexAll
}

// lexVariableOrFunction creates a variable or function token.
//
// Transitions:
//  -> lexAll
func lexVariableOrFunction(l *Lexer) stateFn {
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
				l.emitEmpty(token.Sqrt)
			case "sin(":
				l.emitEmpty(token.Sin)
			case "cos(":
				l.emitEmpty(token.Cos)
			case "tan(":
				l.emitEmpty(token.Tan)
			default:
				l.emit(token.UnkownFunktion)
			}
			return lexAll
		}

		l.emitSingle(token.InvalidCharacterInVariable)
		return lexAll
	}

	l.emit(token.Var)
	return lexAll
}
