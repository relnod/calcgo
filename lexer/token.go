package lexer

// TokenType describes the type of a token
type TokenType byte

// Token types
const (
	TEOF TokenType = iota

	literal_beg
	// Numbers
	TInt // [0-9]+
	TDec // [0-9]+\.[0-9]+
	TBin // 0b[01]+
	THex // 0x[0-9A-F]+
	TExp // [0-9]+\^[0-9]+

	// Variable
	TVar // [a-zA-Z]+
	literal_end

	operator_beg
	// Operators
	TOpPlus  // "+"
	TOpMinus // "-"
	TOpMult  // "*"
	TOpDiv   // "/"
	operator_end

	function_beg
	// Functions
	TFnSqrt // "sqrt("
	TFnSin  // "sin("
	TFnCos  // "cos("
	TFnTan  // "tan("
	function_end

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

	TFnSqrt: "sqrt",
	TFnSin:  "sin",
	TFnCos:  "cos",
	TFnTan:  "tan",

	TLParen: "(",
	TRParen: ")",

	TInvalidCharacter:           "Invalid Character",
	TInvalidCharacterInNumber:   "Invalid character in number",
	TInvalidCharacterInVariable: "Invalid character in Variabl",
}

// Token represents a token returned by the lexer
type Token struct {
	Type  TokenType
	Value string
}

// IsLitereal returns true if the type of t is a literal.
func (t Token) IsLiteral() bool {
	return literal_beg < t.Type && t.Type < literal_end
}

// IsOperator returns true if the type of t is a operator.
func (t Token) IsOperator() bool {
	return operator_beg < t.Type && t.Type < operator_end
}

// IsFunction returns true if the type of t is a function.
func (t Token) IsFunction() bool {
	return function_beg < t.Type && t.Type < function_end
}

// String converts a token to a string.
func (t Token) String() string {
	return "{\"" + t.Value + "\", " + t.Type.String() + "}"
}

// String converts a token type to a string.
func (t TokenType) String() string {
	if 0 <= t && t < TokenType(len(tokens)) {
		return tokens[t]
	}

	return "Unknown token"
}
