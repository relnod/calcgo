package calcgo

// TokenType describes the type of a token
type TokenType uint

const (
	TInteger                  TokenType = iota
	TDecimal                  TokenType = iota	
	TOperatorPlus             TokenType = iota
	TOperatorMinus            TokenType = iota
	TOperatorMult             TokenType = iota
	TOperatorDiv              TokenType = iota
	TLeftBracket              TokenType = iota
	TRightBracket             TokenType = iota
	TInvalidCharacter         TokenType = iota
	TInvalidCharacterInNumber TokenType = iota
)

type Token struct {
	Type  TokenType
	Value string
}

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

	switch string(str[i]) {
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		token, i = getNummberToken(str, i)
	case "+":
		token = Token{TOperatorPlus, "+"}
	case "-":
		token = Token{TOperatorMinus, "-"}
	case "*":
		token = Token{TOperatorMult, "*"}
	case "/":
		token = Token{TOperatorDiv, "/"}
	case "(":
		token = Token{TLeftBracket, "("}
	case ")":
		token = Token{TRightBracket, ")"}
	case " ":
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

func getNummberToken(str string, i int) (Token, int) {
	var token Token

	var number string
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
		case ')':
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
