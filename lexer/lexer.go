package lexer

type stateFn func(*Lexer) stateFn

// Lexer holds the state of the lexer.
type Lexer struct {
	token   chan Token
	str     string
	pos     int
	lastPos int
}

// Lex takes a string as input and returns a list of tokens.
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

// NewLexer returns a new lexer object.
func NewLexer(str string) *Lexer {
	return &Lexer{str: str, token: make(chan Token, len(str)/3), pos: -1}
}

// Start runs the lexer in a go routine.
func (l *Lexer) Start() {
	go l.run()
}

// NextToken returns the next token from the token chanel.
func (l *Lexer) NextToken() Token {
	return <-l.token
}

// GetChanel returns the token chanel.
func (l *Lexer) GetChanel() chan Token {
	return l.token
}

// current returns the character at the current postion.
func (l *Lexer) current() byte {
	return l.str[l.pos]
}

// stored returnes the string, that is currently stored.
func (l *Lexer) stored() string {
	return l.str[l.lastPos+1 : l.pos+1]
}

// next procedes to the next character and returns it. Also returns indicator,
// wether there is a next character.
func (l *Lexer) next() (byte, bool) {
	if l.pos+1 >= len(l.str) {
		return 0, false
	}
	l.pos++

	return l.current(), true
}

// backup moves the position one character backwards.
func (l *Lexer) backup() {
	l.pos--
}

// emitInternal takes a tokentype and a value to create a token, which it then
// emits.
func (l *Lexer) emitInternal(tokenType TokenType, value string) {
	l.token <- Token{
		Type:  tokenType,
		Value: value,
		Start: l.lastPos + 1,
		End:   l.pos + 1,
	}
}

// emit emits a new token with type tokenType and the currently stored value.
func (l *Lexer) emit(tokenType TokenType) {
	l.emitInternal(tokenType, l.stored())
}

// emitEmpty emits a new token with type tokenType and an empty value.
func (l *Lexer) emitEmpty(tokenType TokenType) {
	l.emitInternal(tokenType, "")
}

// emitEmpty emits a new token with type tokenType and the current character.
func (l *Lexer) emitSingle(tokenType TokenType) {
	l.emitInternal(tokenType, string(l.current()))
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
	case '%':
		tokenType = TOpMod
	case '|':
		tokenType = TOpOr
	case '^':
		tokenType = TOpXor
	case '&':
		tokenType = TOpAnd
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

// lexNumber is the entry state for all number tokens.
//
// Transitions:
//  - 0x       -> lexHex
//  - 0b       -> lexBin
//  - [0-9]+\. -> lexDecimal
//  - [0-9]+\^ -> lexExponential
//  - rest     -> lexAll
func lexNumber(l *Lexer) stateFn {
	if l.current() == '0' {
		b, ok := l.next()
		if ok {
			if b == 'x' {
				return lexHex
			}
			if b == 'b' {
				return lexBin
			}
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

// lexDecimal creates a decimal number token.
//
// Transitions:
//  -> lexAll
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

// lexHex creates a hex number token.
//
// Transitions:
//  -> lexAll
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

// lexBin creates a binary number token.
//
// Transitions:
//  -> lexAll
func lexBin(l *Lexer) stateFn {
	for {
		b, ok := l.next()
		if !ok {
			break
		}

		if isBinDigit(b) {
			continue
		}

		if isWhiteSpace(b) || b == ')' {
			l.backup()
			break
		}

		l.emitSingle(TInvalidCharacterInNumber)
		return lexAll
	}

	l.emit(TBin)
	return lexAll
}

// lexExponential creates an exponential number token.
//
// Transitions:
//  -> lexAll
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

// lexVariableOrFunction creates a variable or function token.
//
// Transitions:
//  -> lexAll
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
			case "sin(":
				l.emitEmpty(TFnSin)
			case "cos(":
				l.emitEmpty(TFnCos)
			case "tan(":
				l.emitEmpty(TFnTan)
			default:
				l.emit(TFnUnkown)
			}
			return lexAll
		}

		l.emitSingle(TInvalidCharacterInVariable)
		return lexAll
	}

	l.emit(TVar)
	return lexAll
}
