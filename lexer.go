package calcgo

type stateFn func(*Lexer) stateFn

// TokenType describes the type of a token
type TokenType uint

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
	var tokens []Token

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
	return &Lexer{str: str, token: make(chan Token), pos: -1}
}

// Start runs the lexer in a go routine
func (l *Lexer) Start() {
	go l.run()
}

func (l *Lexer) NextToken() Token {
	return <-l.token
}

func (l *Lexer) hasNext() bool {
	return l.pos+1 < len(l.str)
}

func (l *Lexer) current() byte {
	return l.str[l.pos]
}

func (l *Lexer) next() byte {
	l.pos++

	return l.current()
}

func (l *Lexer) backup() {
	l.pos--
}

func (l *Lexer) peek() byte {
	s := l.next()

	l.backup()
	return s
}

func (l *Lexer) emit(tokenType TokenType) {
	l.token <- Token{Type: tokenType, Value: l.str[l.lastPos : l.pos+1]}
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

func lexAll(l *Lexer) stateFn {
	var tokenType TokenType
	if !l.hasNext() {
		return nil
	}

	l.lastPos = l.pos + 1

	switch l.next() {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return lexNumber
	case '+':
		tokenType = TOperatorPlus
	case '-':
		tokenType = TOperatorMinus
	case '*':
		tokenType = TOperatorMult
	case '/':
		tokenType = TOperatorDiv
	case '(':
		tokenType = TLeftBracket
	case ')':
		tokenType = TRightBracket
	case ' ':
		return lexAll
	default:
		l.emit(TInvalidCharacter)
		return lexAll
	}

	l.emitEmpty(tokenType)
	return lexAll
}

func lexNumber(l *Lexer) stateFn {
	tokenType := TInteger

loop:
	for l.hasNext() {
		switch l.next() {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		case '.':
			tokenType = TDecimal
		case ' ', ')':
			l.backup()
			break loop
		default:
			l.lastPos = l.pos
			l.emit(TInvalidCharacterInNumber)
			return lexAll
		}
	}

	l.emit(tokenType)
	return lexAll
}
