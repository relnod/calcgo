package calcgo

type stateFn func(*Lexer) stateFn

// TokenType describes the type of a token
type TokenType byte

// Token types
const (
	TEOF TokenType = iota
	TInteger
	TDecimal
	TOperatorPlus
	TOperatorMinus
	TOperatorMult
	TOperatorDiv
	TLeftBracket
	TRightBracket
	TInvalidCharacter
	TInvalidCharacterInNumber
)

// Token represents a token returned by the lexer
type Token struct {
	Type  TokenType
	Value string
}

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

func (l *Lexer) peek() (byte, bool) {
	s, ok := l.next()

	if !ok {
		return 0, false
	}

	l.backup()
	return s, true
}

func (l *Lexer) emit(tokenType TokenType) {
	l.token <- Token{Type: tokenType, Value: l.str[l.lastPos+1 : l.pos+1]}
}

func (l *Lexer) emitEmpty(tokenType TokenType) {
	l.token <- Token{Type: tokenType, Value: ""}
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
	if isWhiteSpace(b) {
		return lexAll
	}

	switch b {
	case '+':
		tokenType = TOperatorPlus
	case '-':
		if b, ok := l.peek(); ok && isDigit(b) {
			return lexNumber
		}
		tokenType = TOperatorMinus
	case '*':
		tokenType = TOperatorMult
	case '/':
		tokenType = TOperatorDiv
	case '(':
		tokenType = TLeftBracket
	case ')':
		tokenType = TRightBracket
	default:
		l.emit(TInvalidCharacter)
		return lexAll
	}

	l.emitEmpty(tokenType)
	return lexAll
}

func lexNumber(l *Lexer) stateFn {
	tokenType := TInteger

	for {
		b, ok := l.next()
		if !ok {
			break
		}

		if isDigit(b) {
			continue
		} else if b == '.' {
			tokenType = TDecimal
		} else if isWhiteSpace(b) || b == ')' {
			l.backup()
			break
		} else {
			l.lastPos = l.pos - 1
			l.emit(TInvalidCharacterInNumber)
			return lexAll
		}
	}

	l.emit(tokenType)
	return lexAll
}
