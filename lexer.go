package calcgo

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

// Lexer
type Lexer struct {
	token chan Token
	str   string
	index int
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
	return &Lexer{str: str, token: make(chan Token)}
}

// Start runs the lexer in a go routine
func (l *Lexer) Start() {
	go l.run()
}

func (l *Lexer) NextToken() Token {
	return <-l.token
}

func (l *Lexer) run() {
	for l.index = 0; l.index < len(l.str); l.index++ {
		var token Token
		token = l.getNextToken()
		l.token <- token
	}
	close(l.token)
}

func (l *Lexer) current() string {
	return string(l.str[l.index])
}

func (l *Lexer) getNextToken() Token {
	var token Token

	switch l.str[l.index] {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		token = l.getNumberToken()
	case '+':
		token = Token{TOperatorPlus, ""}
	case '-':
		token = Token{TOperatorMinus, ""}
	case '*':
		token = Token{TOperatorMult, ""}
	case '/':
		token = Token{TOperatorDiv, ""}
	case '(':
		token = Token{TLeftBracket, ""}
	case ')':
		token = Token{TRightBracket, ""}
	case ' ':
		l.index++
		if l.index == len(l.str) {
			return Token{TEOF, ""}
		}
		token = l.getNextToken()
	default:
		token = Token{TInvalidCharacter, l.current()}
	}

	return token
}

func (l *Lexer) getNumberToken() Token {
	var token Token

	number := l.current()
	l.index++
	isDecimal := false
	for l.index < len(l.str) {
		switch l.str[l.index] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			number += l.current()
		case '.':
			number += l.current()
			isDecimal = true
		case ' ':
			goto endNumberToken
		case '+', '-', '*', '/', ')':
			l.index--
			goto endNumberToken
		default:
			return Token{TInvalidCharacterInNumber, l.current()}
		}
		l.index++
	}

endNumberToken:
	var tokenType TokenType
	if isDecimal {
		tokenType = TDecimal
	} else {
		tokenType = TInteger
	}

	token = Token{tokenType, number}

	return token
}
