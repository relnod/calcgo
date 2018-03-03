package token

import "strconv"

// Type describes the type of a token
type Type byte

// Token types
const (
	EOF Type = iota

	literalBeg
	// Numbers
	Int // [0-9]+
	Dec // [0-9]+\.[0-9]+
	Bin // 0b[01]+
	Hex // 0x[0-9A-F]+
	Exp // [0-9]+\^[0-9]+

	// Variable
	Var // [a-zA-Z]+
	literalEnd

	operatorBeg
	// Operators
	Plus  // "+"
	Minus // "-"
	Mult  // "*"
	Div   // "/"
	Mod   // "%"
	Or    // "|"
	Xor   // "^"
	And   // "&"
	operatorEnd

	functionBeg
	// Functions
	Sqrt // "sqrt("
	Sin  // "sin("
	Cos  // "cos("
	Tan  // "tan("
	UnkownFunktion
	functionEnd

	// Parens
	ParenL // "("
	ParenR // ")"

	// Errors
	InvalidCharacter
	InvalidCharacterInNumber
	InvalidCharacterInVariable
)

var tokens = [...]string{
	EOF: "EOF",

	Int: "Integer",
	Dec: "Decimal",
	Bin: "Binary",
	Hex: "HexaDecimal",
	Exp: "Exponential",

	Var: "Variable",

	Plus:  "+",
	Minus: "-",
	Mult:  "*",
	Div:   "/",
	Mod:   "%",
	Or:    "|",
	Xor:   "^",
	And:   "&",

	Sqrt: "sqrt",
	Sin:  "sin",
	Cos:  "cos",
	Tan:  "tan",

	ParenL: "(",
	ParenR: ")",

	InvalidCharacter:           "Invalid Character",
	InvalidCharacterInNumber:   "Invalid character in number",
	InvalidCharacterInVariable: "Invalid character in Variabl",
	UnkownFunktion:             "Unknown function",
}

// Token represents a token returned by the lexer
type Token struct {
	Type  Type
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
func (t Type) String() string {
	if 0 <= t && t < Type(len(tokens)) {
		return tokens[t]
	}

	return "Unknown token"
}
