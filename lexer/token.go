package lexer

// TokenType describes the type of a token
type TokenType byte

// Token types
const (
	TEOF TokenType = iota

	// Numbers
	TInt // [0-9]+
	TDec // [0-9]+\.[0-9]+
	TExp // [0-9]+\^[0-9]+

	// Variable
	TVar // [a-zA-Z]+

	// Operators
	TOpPlus  // "+"
	TOpMinus // "-"
	TOpMult  // "*"
	TOpDiv   // "/"

	// Functions
	TFnSqrt // "sqrt("
	TFnSin  // "sin("
	TFnCos  // "cos("
	TFnTan  // "tan("

	// Parens
	TLParen // "("
	TRParen // ")"

	// Errors
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
	return t.Type == TFnSqrt ||
		t.Type == TFnSin ||
		t.Type == TFnCos ||
		t.Type == TFnTan
}

func (t Token) String() string {
	return "{\"" + t.Value + "\", " + t.Type.String() + "}"
}

func (t TokenType) String() string {
	switch t {
	case TInt:
		return "Integer"
	case TDec:
		return "Decimal"
	case TVar:
		return "Variable"
	case TOpPlus:
		return "Plus"
	case TOpMinus:
		return "Minus"
	case TOpMult:
		return "Mult"
	case TOpDiv:
		return "Div"
	case TFnSqrt:
		return "Sqrt"
	case TFnSin:
		return "Sin"
	case TFnCos:
		return "Cos"
	case TFnTan:
		return "Tan"
	case TLParen:
		return "Left Bracket"
	case TRParen:
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
