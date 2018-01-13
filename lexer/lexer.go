package lexer

type stateFn func(*Lexer) stateFn

// Lexer holds the state of the lexer
type Lexer struct {
	token   chan Token
	str     string
	pos     int
	lastPos int
}

// Lex takes a string as input and returns a list of tokens
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
func Lex(str string) []Token {
	if len(str) == 0 {
		return nil
	}

	tokens := make([]Token, 0, len(str)/2)

	lexer := NewLexer(str)
	lexer.Start()

	for {
		token := lexer.NextToken()
		if token.Type == TEOF {
			break
		}
		tokens = append(tokens, token)
	}

	return tokens
}

// NewLexer returns a new lexer object
func NewLexer(str string) *Lexer {
	return &Lexer{str: str, token: make(chan Token, len(str)/3), pos: -1}
}

// Start runs the lexer in a go routine
func (l *Lexer) Start() {
	go l.run()
}

// NextToken returns the next token from the token chanel
func (l *Lexer) NextToken() Token {
	return <-l.token
}

// GetChanel returns the token chanel
func (l *Lexer) GetChanel() chan Token {
	return l.token
}

func (l *Lexer) current() byte {
	return l.str[l.pos]
}

func (l *Lexer) stored() string {
	return l.str[l.lastPos+1 : l.pos+1]
}

func (l *Lexer) next() (byte, bool) {
	if l.pos+1 >= len(l.str) {
		return 0, false
	}
	l.pos++

	return l.current(), true
}

func (l *Lexer) backup() {
	l.pos--
}

func (l *Lexer) emit(tokenType TokenType) {
	l.token <- Token{Type: tokenType, Value: l.stored()}
}

func (l *Lexer) emitEmpty(tokenType TokenType) {
	l.token <- Token{Type: tokenType, Value: ""}
}

func (l *Lexer) emitSingle(tokenType TokenType) {
	l.lastPos = l.pos - 1
	l.token <- Token{Type: tokenType, Value: l.stored()}
}

func (l *Lexer) run() {
	for state := lexAll; state != nil; {
		state = state(l)
	}
	close(l.token)
}

func isWhiteSpace(b byte) bool {
	return b == ' '
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func isHexDigit(b byte) bool {
	return isDigit(b) || b >= 'A' && b <= 'F'
}

func isLetter(b byte) bool {
	return b >= 'a' && b <= 'z'
}

func lexAll(l *Lexer) stateFn {
	var tokenType TokenType

	l.lastPos = l.pos

	b, ok := l.next()
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
		tokenType = TOpPlus
	case '-':
		if b, ok := l.next(); ok && isDigit(b) {
			return lexNumber
		}
		tokenType = TOpMinus
	case '*':
		tokenType = TOpMult
	case '/':
		tokenType = TOpDiv
	case '(':
		tokenType = TLParen
	case ')':
		tokenType = TRParen
	default:
		l.emit(TInvalidCharacter)
		return lexAll
	}

	l.emitEmpty(tokenType)
	return lexAll
}

func lexNumber(l *Lexer) stateFn {
	if l.current() == '0' {
		if b, ok := l.next(); ok && b == 'x' {
			return lexHex
		}
		l.backup()
	}

	for {
		b, ok := l.next()
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
			l.backup()
			break
		}

		l.emitSingle(TInvalidCharacterInNumber)
		return lexAll
	}

	l.emit(TInt)
	return lexAll
}

func lexDecimal(l *Lexer) stateFn {
	for {
		b, ok := l.next()
		if !ok {
			break
		}

		if isDigit(b) {
			continue
		}

		if isWhiteSpace(b) || b == ')' {
			l.backup()
			break
		}

		l.emitSingle(TInvalidCharacterInNumber)
		return lexAll
	}

	l.emit(TDec)
	return lexAll
}

func lexExponential(l *Lexer) stateFn {
	for {
		b, ok := l.next()
		if !ok {
			break
		}

		if isDigit(b) {
			continue
		}

		if isWhiteSpace(b) || b == ')' {
			l.backup()
			break
		}

		l.emitSingle(TInvalidCharacterInNumber)
		return lexAll
	}

	l.emit(TExp)
	return lexAll
}

func lexHex(l *Lexer) stateFn {
	for {
		b, ok := l.next()
		if !ok {
			break
		}

		if isHexDigit(b) {
			continue
		}

		if isWhiteSpace(b) || b == ')' {
			l.backup()
			break
		}

		l.emitSingle(TInvalidCharacterInNumber)
		return lexAll
	}

	l.emit(THex)
	return lexAll
}

func lexVariableOrFunction(l *Lexer) stateFn {
	for {
		b, ok := l.next()
		if !ok {
			break
		}

		if isLetter(b) {
			continue
		}

		if isWhiteSpace(b) || b == ')' {
			l.backup()
			break
		}

		if b == '(' {
			switch l.stored() {
			case "sqrt(":
				l.emitEmpty(TFnSqrt)
				return lexAll
			case "sin(":
				l.emitEmpty(TFnSin)
				return lexAll
			case "cos(":
				l.emitEmpty(TFnCos)
				return lexAll
			case "tan(":
				l.emitEmpty(TFnTan)
				return lexAll
			}
		}

		l.emitSingle(TInvalidCharacterInVariable)
		return lexAll
	}

	l.emit(TVar)
	return lexAll
}
