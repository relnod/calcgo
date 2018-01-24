package lexer_test

import (
	"regexp"
	"testing"

	"github.com/relnod/calcgo/lexer"
)

// |         | in       |          | out      |
// |---------|----------|----------|----------|
// | START   | 0        | NUMBER   |          |
// | START   | 1-9      | INT      |          |
// | START   | (        | START    | TLPAREN  |
// | START   | )        | START    | TRPAREN  |
// | START   | +        | START    | TOpPlus  |
// | START   | -        | START    | TOpMinus |
// | START   | *        | START    | TOpMult  |
// | START   | /        | START    | TOpDiv   |
// | START   | %        | START    | TOpMod   |
// | START   | |        | START    | TOpOr    |
// | START   | ^        | START    | TOpXor   |
// | START   | &        | START    | TOpAnd   |
// | START   | s        | S        |          |
// | START   | c        | C        |          |
// | START   | t        | T        |          |
// | START   | a-z      | VAR      |          |
// |         |          |          |          |
// | NUMBER  | b        | BYTE     |          |
// | NUMBER  | x        | HEX      |          |
// | NUMBER  | .        | DEC      |          |
// |         |          |          |          |
// | BIN     | 0,1      | BIN2     |          |
// |         |          |          |          |
// | BIN2    | 0,1      | BIN2     |          |
// | BIN2    | " "      | START    | TBin     |
// |         |          |          |          |
// | HEX     | 1-9, A-F | HEX2     |          |
// |         |          |          |          |
// | HEX2    | 1-0, A-F | HEX2     |          |
// | HEX2    | " "      | START    | THex     |
// |         |          |          |          |
// | INT     | 1-9      | INT      |          |
// | INT     | .        | DEC      |          |
// | INT     | ^        | EXP      |          |
// | INT     | " "      | START    | TInt     |
// |         |          |          |          |
// | DEC     | 1-9      | DEC      |          |
// | DEC     | " "      | START    | TDec     |
// |         |          |          |          |
// | Exp     | 1-9      | Exp      |          |
// | Exp     | " "      | START    | TExp     |
// |         |          |          |          |
// | S       | q        | SQ       |          |
// | S       | i        | SI       |          |
// | S       | a-z      | VAR      |          |
// | S       | " "      | STAR     | TVar     |
// |         |          |          |          |
// | SQ      | r        | SQR      |          |
// | SQ      | a-z      | VAR      |          |
// | SQ      | " "      | STAR     | TVar     |
// |         |          |          |          |
// | SQR     | t        | SQRT     |          |
// | SQR     | a-z      | VAR      |          |
// | SQR     | " "      | STAR     | TVar     |
// |         |          |          |          |
// | SQRT    | (        | START    | TFnSqrt  |
// | SQRT    | a-z      | VAR      |          |
// | SQRT    | " "      | STAR     | TVar     |
// |         |          |          |          |
// | SI      | n        | SIN      |          |
// | SI      | a-z      | VAR      |          |
// | SI      | " "      | STAR     | TVar     |
// |         |          |          |          |
// | SIN     | (        | START    | TFnSin   |
// | SIN     | a-z      | VAR      |          |
// | SIN     | " "      | STAR     | TVar     |
// |         |          |          |          |
// | C       | o        | CO       |          |
// | C       | a-z      | VAR      |          |
// | C       | " "      | STAR     | TVar     |
// |         |          |          |          |
// | CO      | s        | COS      |          |
// | CO      | a-z      | VAR      |          |
// | CO      | " "      | STAR     | TVar     |
// |         |          |          |          |
// | Cos     | (        | START    | TFnCos   |
// | COS     | a-z      | VAR      |          |
// | COS     | " "      | STAR     | TVar     |
// |         |          |          |          |
// | T       | a        | TA       |          |
// | T       | a-z      | VAR      |          |
// | T       | " "      | STAR     | TVar     |
// |         |          |          |          |
// | TA      | n        | TAN      |          |
// | TA      | a-z      | VAR      |          |
// | TA      | " "      | STAR     | TVar     |
// |         |          |          |          |
// | TAN     | (        | START    | TFnTan   |
// | TAN     | a-z      | VAR      |          |
// | TAN     | " "      | STAR     | TVar     |
// |         |          |          |          |
// | VAR     | a-z      | VAR      |          |
// | VAR     | " "      | START    | TVar     |

type Spec struct {
	start       string
	transitions []Transition

	current string
}

func NewSpec() *Spec {
	return &Spec{}
}

func (s *Spec) SetStart(start string) {
	s.start = start
}

func (s *Spec) AddTransition(from, in, to string, out lexer.TokenType) {
	s.transitions = append(transitions, &Transition{
		from: start,
		in:   regexp.Compile(in),
		to:   to,
		out:  out,
	})
}

type Transition struct {
	from string
	in   regexp

	to  string
	out *lexer.TokenType
}

func TestSpec(t *testing.T) {
	spec := NewSpec()
	spec.SetStart("START")

	spec.AddTransition("START", "(", "START", lexer.TLParen)
	spec.AddTransition("START", ")", "START", lexer.TLParen)
	spec.AddTransition("START", "+", "START", lexer.TOpPlus)
	spec.AddTransition("START", "+", "START", lexer.TOpPlus)
	spec.AddTransition("START", "-", "START", lexer.TOpMinus)
	spec.AddTransition("START", "*", "START", lexer.TOpMult)
	spec.AddTransition("START", "/", "START", lexer.TOpDiv)
	spec.AddTransition("START", "%", "START", lexer.TOpMod)
	spec.AddTransition("START", "|", "START", lexer.TOpOr)
	spec.AddTransition("START", "^", "START", lexer.TOpXor)
	spec.AddTransition("START", "&", "START", lexer.TOpAnd)

	spec.AddTransition("START", "0", "NUMBER", nil)
	spec.AddTransition("START", "[1-9]", "INT", nil)

	spec.AddTransition("START", "s", "S", nil)
	spec.AddTransition("START", "c", "C", nil)
	spec.AddTransition("START", "t", "T", nil)
	spec.AddTransition("START", "[a-z]", "VAR", nil)

	spec.AddTransition("NUMBER", "b", "BIN", nil)
	spec.AddTransition("NUMBER", "x", "HEX", nil)
	spec.AddTransition("NUMBER", ".", "DEC", nil)

	spec.AddTransition("BIN", "0|1", "BIN2", nil)
	spec.AddTransition("BIN2", "0|1", "BIN2", nil)
	spec.AddTransition("BIN2", " ", "START", lexer.TBin)

	spec.AddTransition("HEX", "[1-9A-F]", "HEX2", nil)
	spec.AddTransition("HEX2", "[1-9A-F]", "HEX2", nil)
	spec.AddTransition("HEX2", " ", "START", lexer.THex)

	spec.AddTransition("INT", "[1-9]", "INT", nil)
	spec.AddTransition("INT", ".", "DEC", nil)
	spec.AddTransition("INT", "^", "Exp", nil)
	spec.AddTransition("INT", " ", "START", lexer.TInt)

	spec.AddTransition("DEC", "[1-9]", "DEC2", nil)
	spec.AddTransition("DEC2", "[1-9]", "DEC2", nil)
	spec.AddTransition("DEC", " ", "START", lexer.TDec)

	spec.AddTransition("Exp", "[1-9]", "Exp2", nil)
	spec.AddTransition("Exp2", "[1-9]", "Exp2", nil)
	spec.AddTransition("Exp2", " ", "START", lexer.TDec)

	spec.AddTransition("S", "q", "SQ", nil)
	spec.AddTransition("S", "i", "SI", nil)
	spec.AddTransition("S", "[a-z]", "VAR", nil)
	spec.AddTransition("S", " ", "START", lexer.TVar)

	spec.AddTransition("SQ", "r", "SQR", nil)
	spec.AddTransition("SQ", "[a-z]", "VAR", nil)
	spec.AddTransition("SQ", " ", "START", lexer.TVAR)

	spec.AddTransition("SQR", "t", "SQRT", nil)
	spec.AddTransition("SQR", "[a-z]", "VAR", nil)
	spec.AddTransition("SQR", " ", "START", lexer.TVAR)

	spec.AddTransition("SQRT", "(", "START", lexer.TFnSqrt)
	spec.AddTransition("SQRT", "[a-z]", "VAR", nil)
	spec.AddTransition("SQRT", " ", "START", lexer.TVAR)

	spec.AddTransition("SI", "n", "SIN", nil)
	spec.AddTransition("SI", "[a-z]", "VAR", nil)
	spec.AddTransition("SI", " ", "START", lexer.TVAR)

	spec.AddTransition("SIN", "(", "START", lexer.TFnSin)
	spec.AddTransition("SIN", "[a-z]", "VAR", nil)
	spec.AddTransition("SIN", " ", "START", lexer.TVAR)

	spec.AddTransition("C", "o", "CO", nil)
	spec.AddTransition("C", "[a-z]", "VAR", nil)
	spec.AddTransition("C", " ", "START", lexer.TVar)

	spec.AddTransition("CO", "s", "COS", nil)
	spec.AddTransition("CO", "[a-z]", "VAR", nil)
	spec.AddTransition("CO", " ", "START", lexer.TVar)

	spec.AddTransition("COS", "(", "START", lexer.TCOS)
	spec.AddTransition("COS", "[a-z]", "VAR", nil)
	spec.AddTransition("COS", " ", "START", lexer.TVar)

	spec.AddTransition("T", "a", "TA", nil)
	spec.AddTransition("T", "[a-z]", "VAR", nil)
	spec.AddTransition("T", " ", "START", lexer.TVar)

	spec.AddTransition("TA", "n", "TAN", nil)
	spec.AddTransition("TA", "[a-z]", "VAR", nil)
	spec.AddTransition("TA", " ", "START", lexer.TVar)

	spec.AddTransition("TAN", "(", "START", lexer.TFnTan)
	spec.AddTransition("TAN", "[a-z]", "VAR", nil)
	spec.AddTransition("TAN", " ", "START", lexer.TVar)
}
