package calcgo

// TokenType describes the type of a token
type TokenType uint

// Token types
const (
	TInteger TokenType = iota
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

	for i := 0; i < len(str); i++ {
		var token Token
		token, i = getNextToken(str, i)
		if i == -1 {
			break
		}
		tokens = append(tokens, token)
	}

	return tokens
}

func getNextToken(str string, i int) (Token, int) {
	var token Token

	switch str[i] {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		token, i = getNumberToken(str, i)
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
		i++
		if i == len(str) {
			return token, -1
		}
		token, i = getNextToken(str, i)
	default:
		token = Token{TInvalidCharacter, string(str[i])}
	}

	return token, i
}

func getNumberToken(str string, i int) (Token, int) {
	var token Token

	number := string(str[i])
	i++
	isDecimal := false
	for i < len(str) {
		switch str[i] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			number += string(str[i])
		case '.':
			number += string(str[i])
			isDecimal = true
		case ' ':
			goto endNumberToken
		case '+', '-', '*', '/', ')':
			i--
			goto endNumberToken
		default:
			return Token{TInvalidCharacterInNumber, string(str[i])}, i
		}
		i++
	}

endNumberToken:
	var tokenType TokenType
	if isDecimal {
		tokenType = TDecimal
	} else {
		tokenType = TInteger
	}

	token = Token{tokenType, number}

	return token, i
}
