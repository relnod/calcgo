package lexer

// TokenType describes the type of a token
type TokenType byte

// Token types
const (
	TEOF TokenType = iota
	TInteger
	TDecimal
	TVariable
	TOperatorPlus
	TOperatorMinus
	TOperatorMult
	TOperatorDiv
	TFuncSqrt
	TFuncSin
	TFuncCos
	TFuncTan
	TLeftBracket
	TRightBracket
	TInvalidCharacter
	TInvalidCharacterInNumber
	TInvalidCharacterInVariable
)

// Token represents a token returned by the lexer
type Token struct {
	Type  TokenType
	Value string
}

// IsFunction returns true if the type of t is a function.
func (t Token) IsFunction() bool {
	return t.Type == TFuncSqrt ||
		t.Type == TFuncSin ||
		t.Type == TFuncCos ||
		t.Type == TFuncTan
}

func (t Token) String() string {
	return "{\"" + t.Value + "\", " + t.Type.String() + "}"
}

func (t TokenType) String() string {
	switch t {
	case TInteger:
		return "Integer"
	case TDecimal:
		return "Decimal"
	case TVariable:
		return "Variable"
	case TOperatorPlus:
		return "Plus"
	case TOperatorMinus:
		return "Minus"
	case TOperatorMult:
		return "Mult"
	case TOperatorDiv:
		return "Div"
	case TFuncSqrt:
		return "Sqrt"
	case TFuncSin:
		return "Sin"
	case TFuncCos:
		return "Cos"
	case TFuncTan:
		return "Tan"
	case TLeftBracket:
		return "Left Bracket"
	case TRightBracket:
		return "RightBracket"
	case TInvalidCharacter:
		return "Invalid Character"
	case TInvalidCharacterInNumber:
		return "Invalid Character in Number"
	case TInvalidCharacterInVariable:
		return "Invalid Character in Variable"
	default:
		return "Undefined Token Type"
	}
}
