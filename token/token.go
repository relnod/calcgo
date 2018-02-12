package token

import "strconv"

// TokenType describes the type of a token
type TokenType byte

// Token types
const (
	TEOF TokenType = iota

	literalBeg
	// Numbers
	TInt // [0-9]+
	TDec // [0-9]+\.[0-9]+
	TBin // 0b[01]+
	THex // 0x[0-9A-F]+
	TExp // [0-9]+\^[0-9]+

	// Variable
	TVar // [a-zA-Z]+
	literalEnd

	operatorBeg
	// Operators
	TOpPlus  // "+"
	TOpMinus // "-"
	TOpMult  // "*"
	TOpDiv   // "/"
	TOpMod   // "%"
	TOpOr    // "|"
	TOpXor   // "^"
	TOpAnd   // "&"
	operatorEnd

	functionBeg
	// Functions
	TFnSqrt // "sqrt("
	TFnSin  // "sin("
	TFnCos  // "cos("
	TFnTan  // "tan("
	TFnUnkown
	functionEnd

	// Parens
	TLParen // "("
	TRParen // ")"

	// Errors
	TInvalidCharacter
	TInvalidCharacterInNumber
	TInvalidCharacterInVariable
)

var tokens = [...]string{
	TEOF: "EOF",

	TInt: "Integer",
	TDec: "Decimal",
	TBin: "Binary",
	THex: "HexaDecimal",
	TExp: "Exponential",

	TVar: "Variable",

	TOpPlus:  "+",
	TOpMinus: "-",
	TOpMult:  "*",
	TOpDiv:   "/",
	TOpMod:   "%",
	TOpOr:    "|",
	TOpXor:   "^",
	TOpAnd:   "&",

	TFnSqrt: "sqrt",
	TFnSin:  "sin",
	TFnCos:  "cos",
	TFnTan:  "tan",

	TLParen: "(",
	TRParen: ")",

	TInvalidCharacter:           "Invalid Character",
	TInvalidCharacterInNumber:   "Invalid character in number",
	TInvalidCharacterInVariable: "Invalid character in Variabl",
	TFnUnkown:                   "Unkown function",
}

// Token represents a token returned by the lexer
type Token struct {
	Type  TokenType
	Value string
	Start int
	End   int
}

// IsLiteral returns true if the type of t is a literal.
func (t Token) IsLiteral() bool {
	return literalBeg < t.Type && t.Type < literalEnd
}

// IsOperator returns true if the type of t is a operator.
func (t Token) IsOperator() bool {
	return operatorBeg < t.Type && t.Type < operatorEnd
}

// IsFunction returns true if the type of t is a function.
func (t Token) IsFunction() bool {
	return functionBeg < t.Type && t.Type < functionEnd
}

// String converts a token to a string.
func (t Token) String() string {
	return "{" +
		"Value: '" + t.Value + "', " +
		"Type: '" + t.Type.String() + "', " +
		"Start: '" + strconv.Itoa(t.Start) + "', " +
		"End: '" + strconv.Itoa(t.End) + "', " +
		"}"
}

// String converts a token type to a string.
func (t TokenType) String() string {
	if 0 <= t && t < TokenType(len(tokens)) {
		return tokens[t]
	}

	return "Unknown token"
}
