package calcgo_test

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	
	"gitlab.com/relnod/calcgo"
)

type Tokens []calcgo.Token

func tokenTypeToString(tokenType calcgo.TokenType) string {
	switch calcgo.TokenType(tokenType) {
	case calcgo.TInteger: 
		return "Number"
	case calcgo.TOperatorPlus: 
		return "Plus"
	case calcgo.TOperatorMinus: 
		return "Minus"
	case calcgo.TOperatorMult: 
		return "Mult"
	case calcgo.TOperatorDiv: 
		return "Div"
	case calcgo.TLeftBracket: 
		return "Left Bracket"
	case calcgo.TRightBracket: 
		return "RightBracket"
	case calcgo.TInvalidCharacterInNumber:
		return "Invalid Character in Number"
	case calcgo.TInvalidCharacter:
		return "Invalid Character"
	default:
		return "Undefined Token Type"
	}
}

func tokenToString(t calcgo.Token) string {
	return "{\"" + t.Value + "\", " + tokenTypeToString(t.Type) + "}"
}

func tokensToString(tokens []calcgo.Token) string {
	str := ""

	for i := 0; i < len(tokens); i++ {
		str += tokenToString(tokens[i]) + "\n"
	}

	return str
}

func tokenError(actual Tokens, expected Tokens) string {
	return "Expected: \n" +
		tokensToString(expected) +
		"Actual: \n" +
		tokensToString(actual)
}

func eq(t1 Tokens, t2 Tokens) bool {
	if len(t1) != len(t2) {
		return false
	}

	for i := 0; i < len(t1); i++ {
		if t1[i].Value != t2[i].Value {
			return false
		}

		if t1[i].Type != t2[i].Type {
			return false
		}
	}

	return true
}

func shouldEqualToken(actual interface{}, expected ...interface{}) string {
	t1 := actual.([]calcgo.Token)
	t2 := expected[0].([]calcgo.Token)

	if eq(t1, t2) {
		return ""
	}

	return tokenError(t1, t2) + "(Should be Equal)"
}

func shouldNotEqualToken(actual interface{}, expected ...interface{}) string {
	t1 := actual.([]calcgo.Token)
	t2 := expected[0].([]calcgo.Token)

	if !eq(t1, t2) {
		return ""
	}

	return tokenError(t1, t2) + "(Should not be Equal)"
}

func TestLexer(t *testing.T) {
	Convey("works with numbers", t, func() {
		Convey("single digit", func() {
			So(calcgo.Lex("0"), shouldEqualToken, []calcgo.Token{
				calcgo.Token{Value: "0", Type: calcgo.TInteger},
			})
			So(calcgo.Lex("1"), shouldEqualToken, []calcgo.Token{
				calcgo.Token{Value: "1", Type: calcgo.TInteger},
			})
			So(calcgo.Lex("2"), shouldEqualToken, []calcgo.Token{
				calcgo.Token{Value: "2", Type: calcgo.TInteger},
			})
			So(calcgo.Lex("3"), shouldEqualToken, []calcgo.Token{
				calcgo.Token{Value: "3", Type: calcgo.TInteger},
			})
			So(calcgo.Lex("4"), shouldEqualToken, []calcgo.Token{
				calcgo.Token{Value: "4", Type: calcgo.TInteger},
			})
			So(calcgo.Lex("5"), shouldEqualToken, []calcgo.Token{
				calcgo.Token{Value: "5", Type: calcgo.TInteger},
			})
			So(calcgo.Lex("6"), shouldEqualToken, []calcgo.Token{
				calcgo.Token{Value: "6", Type: calcgo.TInteger},
			})
			So(calcgo.Lex("7"), shouldEqualToken, []calcgo.Token{
				calcgo.Token{Value: "7", Type: calcgo.TInteger},
			})
			So(calcgo.Lex("8"), shouldEqualToken, []calcgo.Token{
				calcgo.Token{Value: "8", Type: calcgo.TInteger},
			})
			So(calcgo.Lex("9"), shouldEqualToken, []calcgo.Token{
				calcgo.Token{Value: "9", Type: calcgo.TInteger},
			})
		})

		Convey("multiple digits", func() {
			So(calcgo.Lex("10"), shouldEqualToken, []calcgo.Token{
				calcgo.Token{Value: "10", Type: calcgo.TInteger},
			})
			So(calcgo.Lex("10123"), shouldEqualToken, []calcgo.Token{
				calcgo.Token{Value: "10123", Type: calcgo.TInteger},
			})
		})

		Convey("decimals", func() {
			So(calcgo.Lex("1.0"), shouldEqualToken, []calcgo.Token{
				calcgo.Token{Value: "1.0", Type: calcgo.TDecimal},
			})
			So(calcgo.Lex("10.1"), shouldEqualToken, []calcgo.Token{
				calcgo.Token{Value: "10.1", Type: calcgo.TDecimal},
			})
			So(calcgo.Lex("12.3456"), shouldEqualToken, []calcgo.Token{
				calcgo.Token{Value: "12.3456", Type: calcgo.TDecimal},
			})
		})

		Convey("wrong number", func() {
			So(calcgo.Lex("1+"), shouldEqualToken, []calcgo.Token{
				calcgo.Token{Value: "+", Type: calcgo.TInvalidCharacterInNumber},
			})
			So(calcgo.Lex("10123-"), shouldEqualToken, []calcgo.Token{
				calcgo.Token{Value: "-", Type: calcgo.TInvalidCharacterInNumber},
			})
			So(calcgo.Lex("10123- "), shouldEqualToken, []calcgo.Token{
				calcgo.Token{Value: "-", Type: calcgo.TInvalidCharacterInNumber},
			})
		})
	})

	Convey("works with operators", t, func() {
		Convey("plus", func() {
			So(calcgo.Lex("+"), shouldEqualToken, []calcgo.Token{
				calcgo.Token{Value: "+", Type: calcgo.TOperatorPlus},
			})
		})
		Convey("minus", func() {
			So(calcgo.Lex("-"), shouldEqualToken, []calcgo.Token{
				calcgo.Token{Value: "-", Type: calcgo.TOperatorMinus},
			})
		})
		Convey("mult", func() {
			So(calcgo.Lex("*"), shouldEqualToken, []calcgo.Token{
				calcgo.Token{Value: "*", Type: calcgo.TOperatorMult},
			})
		})
		Convey("div", func() {
			So(calcgo.Lex("/"), shouldEqualToken, []calcgo.Token{
				calcgo.Token{Value: "/", Type: calcgo.TOperatorDiv},
			})
		})
	})

	Convey("works with brackets", t, func() {
		Convey("left bracket", func() {
			So(calcgo.Lex("("), shouldEqualToken, []calcgo.Token{
				calcgo.Token{Value: "(", Type: calcgo.TLeftBracket},
			})
		})
		Convey("right bracket", func() {
			So(calcgo.Lex(")"), shouldEqualToken, []calcgo.Token{
				calcgo.Token{Value: ")", Type: calcgo.TRightBracket},
			})
		})
	})

	Convey("works with invalid characters", t, func() {
		So(calcgo.Lex("%"), shouldEqualToken, []calcgo.Token{
			calcgo.Token{Value: "%", Type: calcgo.TInvalidCharacter},			
		})
		So(calcgo.Lex("1 + a"), shouldEqualToken, []calcgo.Token{
			calcgo.Token{Value: "1", Type: calcgo.TInteger},
			calcgo.Token{Value: "+", Type: calcgo.TOperatorPlus},			
			calcgo.Token{Value: "a", Type: calcgo.TInvalidCharacter},						
		})
	})

	Convey("works with mixed types", t, func() {
		So(calcgo.Lex("1 + 2"), shouldEqualToken, []calcgo.Token{
			calcgo.Token{Value: "1", Type: calcgo.TInteger},
			calcgo.Token{Value: "+", Type: calcgo.TOperatorPlus},
			calcgo.Token{Value: "2", Type: calcgo.TInteger},		
		})
		So(calcgo.Lex("1 + 2 + 3 + 4"), shouldEqualToken, []calcgo.Token{
			calcgo.Token{Value: "1", Type: calcgo.TInteger},
			calcgo.Token{Value: "+", Type: calcgo.TOperatorPlus},
			calcgo.Token{Value: "2", Type: calcgo.TInteger},
			calcgo.Token{Value: "+", Type: calcgo.TOperatorPlus},
			calcgo.Token{Value: "3", Type: calcgo.TInteger},
			calcgo.Token{Value: "+", Type: calcgo.TOperatorPlus},
			calcgo.Token{Value: "4", Type: calcgo.TInteger},
		})
		So(calcgo.Lex("(1 + 2) * 2"), shouldEqualToken, []calcgo.Token{
			calcgo.Token{Value: "(", Type: calcgo.TLeftBracket},
			calcgo.Token{Value: "1", Type: calcgo.TInteger},
			calcgo.Token{Value: "+", Type: calcgo.TOperatorPlus},
			calcgo.Token{Value: "2", Type: calcgo.TInteger},
			calcgo.Token{Value: ")", Type: calcgo.TRightBracket},
			calcgo.Token{Value: "*", Type: calcgo.TOperatorMult},
			calcgo.Token{Value: "2", Type: calcgo.TInteger},			
		})
		So(calcgo.Lex("(2 * (1 + 2)) / 2"), shouldEqualToken, []calcgo.Token{
			calcgo.Token{Value: "(", Type: calcgo.TLeftBracket},
			calcgo.Token{Value: "2", Type: calcgo.TInteger},
			calcgo.Token{Value: "*", Type: calcgo.TOperatorMult},
			calcgo.Token{Value: "(", Type: calcgo.TLeftBracket},
			calcgo.Token{Value: "1", Type: calcgo.TInteger},
			calcgo.Token{Value: "+", Type: calcgo.TOperatorPlus},
			calcgo.Token{Value: "2", Type: calcgo.TInteger},
			calcgo.Token{Value: ")", Type: calcgo.TRightBracket},
			calcgo.Token{Value: ")", Type: calcgo.TRightBracket},			
			calcgo.Token{Value: "/", Type: calcgo.TOperatorDiv},
			calcgo.Token{Value: "2", Type: calcgo.TInteger},			
		})
	})
}