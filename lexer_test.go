package calcgo_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/relnod/calcgo"
)

type Tokens []calcgo.Token

func tokenTypeToString(tokenType calcgo.TokenType) string {
	switch calcgo.TokenType(tokenType) {
	case calcgo.TInteger:
		return "Integer"
	case calcgo.TDecimal:
		return "Decimal"
	case calcgo.TVariable:
		return "Variable"
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
	case calcgo.TInvalidCharacter:
		return "Invalid Character"
	case calcgo.TInvalidCharacterInNumber:
		return "Invalid Character in Number"
	case calcgo.TInvalidCharacterInVariable:
		return "Invalid Character in Variable"
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
	Convey("Lexer works with empty string", t, func() {
		So(calcgo.Lex(""), ShouldBeNil)
	})

	Convey("Lexer works with numbers", t, func() {
		Convey("single digit", func() {
			So(calcgo.Lex("0"), shouldEqualToken, []calcgo.Token{
				{Value: "0", Type: calcgo.TInteger},
			})
			So(calcgo.Lex("1"), shouldEqualToken, []calcgo.Token{
				{Value: "1", Type: calcgo.TInteger},
			})
			So(calcgo.Lex("2"), shouldEqualToken, []calcgo.Token{
				{Value: "2", Type: calcgo.TInteger},
			})
			So(calcgo.Lex("3"), shouldEqualToken, []calcgo.Token{
				{Value: "3", Type: calcgo.TInteger},
			})
			So(calcgo.Lex("4"), shouldEqualToken, []calcgo.Token{
				{Value: "4", Type: calcgo.TInteger},
			})
			So(calcgo.Lex("5"), shouldEqualToken, []calcgo.Token{
				{Value: "5", Type: calcgo.TInteger},
			})
			So(calcgo.Lex("6"), shouldEqualToken, []calcgo.Token{
				{Value: "6", Type: calcgo.TInteger},
			})
			So(calcgo.Lex("7"), shouldEqualToken, []calcgo.Token{
				{Value: "7", Type: calcgo.TInteger},
			})
			So(calcgo.Lex("8"), shouldEqualToken, []calcgo.Token{
				{Value: "8", Type: calcgo.TInteger},
			})
			So(calcgo.Lex("9"), shouldEqualToken, []calcgo.Token{
				{Value: "9", Type: calcgo.TInteger},
			})
		})

		Convey("multiple digits", func() {
			So(calcgo.Lex("10"), shouldEqualToken, []calcgo.Token{
				{Value: "10", Type: calcgo.TInteger},
			})
			So(calcgo.Lex("10123"), shouldEqualToken, []calcgo.Token{
				{Value: "10123", Type: calcgo.TInteger},
			})
		})

		Convey("decimals", func() {
			So(calcgo.Lex("1.0"), shouldEqualToken, []calcgo.Token{
				{Value: "1.0", Type: calcgo.TDecimal},
			})
			So(calcgo.Lex("10.1"), shouldEqualToken, []calcgo.Token{
				{Value: "10.1", Type: calcgo.TDecimal},
			})
			So(calcgo.Lex("12.3456"), shouldEqualToken, []calcgo.Token{
				{Value: "12.3456", Type: calcgo.TDecimal},
			})
			So(calcgo.Lex("0.3456"), shouldEqualToken, []calcgo.Token{
				{Value: "0.3456", Type: calcgo.TDecimal},
			})
		})

		Convey("negative numbers", func() {
			So(calcgo.Lex("-1"), shouldEqualToken, []calcgo.Token{
				{Value: "-1", Type: calcgo.TInteger},
			})
			So(calcgo.Lex("-10"), shouldEqualToken, []calcgo.Token{
				{Value: "-10", Type: calcgo.TInteger},
			})
			So(calcgo.Lex("-10.12"), shouldEqualToken, []calcgo.Token{
				{Value: "-10.12", Type: calcgo.TDecimal},
			})
			So(calcgo.Lex("(-1)"), shouldEqualToken, []calcgo.Token{
				{Value: "", Type: calcgo.TLeftBracket},
				{Value: "-1", Type: calcgo.TInteger},
				{Value: "", Type: calcgo.TRightBracket},
			})
		})
	})

	Convey("Lexer works with operators", t, func() {
		Convey("plus", func() {
			So(calcgo.Lex("+"), shouldEqualToken, []calcgo.Token{
				{Value: "", Type: calcgo.TOperatorPlus},
			})
		})
		Convey("minus", func() {
			So(calcgo.Lex("-"), shouldEqualToken, []calcgo.Token{
				{Value: "", Type: calcgo.TOperatorMinus},
			})
		})
		Convey("mult", func() {
			So(calcgo.Lex("*"), shouldEqualToken, []calcgo.Token{
				{Value: "", Type: calcgo.TOperatorMult},
			})
		})
		Convey("div", func() {
			So(calcgo.Lex("/"), shouldEqualToken, []calcgo.Token{
				{Value: "", Type: calcgo.TOperatorDiv},
			})
		})
	})

	Convey("Lexer works with brackets", t, func() {
		Convey("left bracket", func() {
			So(calcgo.Lex("("), shouldEqualToken, []calcgo.Token{
				{Value: "", Type: calcgo.TLeftBracket},
			})
		})
		Convey("right bracket", func() {
			So(calcgo.Lex(")"), shouldEqualToken, []calcgo.Token{
				{Value: "", Type: calcgo.TRightBracket},
			})
		})

		Convey("brackets and numbers", func() {
			So(calcgo.Lex("(1)"), shouldEqualToken, []calcgo.Token{
				{Value: "", Type: calcgo.TLeftBracket},
				{Value: "1", Type: calcgo.TInteger},
				{Value: "", Type: calcgo.TRightBracket},
			})
		})
	})

	Convey("Lexer works with mixed types", t, func() {
		So(calcgo.Lex("1 + 2"), shouldEqualToken, []calcgo.Token{
			{Value: "1", Type: calcgo.TInteger},
			{Value: "", Type: calcgo.TOperatorPlus},
			{Value: "2", Type: calcgo.TInteger},
		})
		So(calcgo.Lex("1 + 2 + 3 + 4"), shouldEqualToken, []calcgo.Token{
			{Value: "1", Type: calcgo.TInteger},
			{Value: "", Type: calcgo.TOperatorPlus},
			{Value: "2", Type: calcgo.TInteger},
			{Value: "", Type: calcgo.TOperatorPlus},
			{Value: "3", Type: calcgo.TInteger},
			{Value: "", Type: calcgo.TOperatorPlus},
			{Value: "4", Type: calcgo.TInteger},
		})
		So(calcgo.Lex("(1 + 2) * 2"), shouldEqualToken, []calcgo.Token{
			{Value: "", Type: calcgo.TLeftBracket},
			{Value: "1", Type: calcgo.TInteger},
			{Value: "", Type: calcgo.TOperatorPlus},
			{Value: "2", Type: calcgo.TInteger},
			{Value: "", Type: calcgo.TRightBracket},
			{Value: "", Type: calcgo.TOperatorMult},
			{Value: "2", Type: calcgo.TInteger},
		})
		So(calcgo.Lex("(2 * (1 + 2)) / 2"), shouldEqualToken, []calcgo.Token{
			{Value: "", Type: calcgo.TLeftBracket},
			{Value: "2", Type: calcgo.TInteger},
			{Value: "", Type: calcgo.TOperatorMult},
			{Value: "", Type: calcgo.TLeftBracket},
			{Value: "1", Type: calcgo.TInteger},
			{Value: "", Type: calcgo.TOperatorPlus},
			{Value: "2", Type: calcgo.TInteger},
			{Value: "", Type: calcgo.TRightBracket},
			{Value: "", Type: calcgo.TRightBracket},
			{Value: "", Type: calcgo.TOperatorDiv},
			{Value: "2", Type: calcgo.TInteger},
		})
	})

	Convey("Lexer handles invalid input correctly", t, func() {
		Convey("returns error with", func() {
			Convey("invalid character in number", func() {
				So(calcgo.Lex("1%"), shouldEqualToken, []calcgo.Token{
					{Value: "%", Type: calcgo.TInvalidCharacterInNumber},
				})
				So(calcgo.Lex("10123o"), shouldEqualToken, []calcgo.Token{
					{Value: "o", Type: calcgo.TInvalidCharacterInNumber},
				})
				So(calcgo.Lex("10123? "), shouldEqualToken, []calcgo.Token{
					{Value: "?", Type: calcgo.TInvalidCharacterInNumber},
				})
			})

			Convey("invalid characters", func() {
				So(calcgo.Lex("%"), shouldEqualToken, []calcgo.Token{
					{Value: "%", Type: calcgo.TInvalidCharacter},
				})
				So(calcgo.Lex("a$"), shouldEqualToken, []calcgo.Token{
					{Value: "$", Type: calcgo.TInvalidCharacterInVariable},
				})
				So(calcgo.Lex("a1"), shouldEqualToken, []calcgo.Token{
					{Value: "1", Type: calcgo.TInvalidCharacterInVariable},
				})
				So(calcgo.Lex("1 + ~"), shouldEqualToken, []calcgo.Token{
					{Value: "1", Type: calcgo.TInteger},
					{Value: "", Type: calcgo.TOperatorPlus},
					{Value: "~", Type: calcgo.TInvalidCharacter},
				})
			})
		})

		Convey("doesn't abort after error", func() {
			So(calcgo.Lex("# + 1"), shouldEqualToken, []calcgo.Token{
				{Value: "#", Type: calcgo.TInvalidCharacter},
				{Value: "", Type: calcgo.TOperatorPlus},
				{Value: "1", Type: calcgo.TInteger},
			})
		})

		Convey("handles multiple errors", func() {
			So(calcgo.Lex("' &"), shouldEqualToken, []calcgo.Token{
				{Value: "'", Type: calcgo.TInvalidCharacter},
				{Value: "&", Type: calcgo.TInvalidCharacter},
			})
			So(calcgo.Lex("# + '"), shouldEqualToken, []calcgo.Token{
				{Value: "#", Type: calcgo.TInvalidCharacter},
				{Value: "", Type: calcgo.TOperatorPlus},
				{Value: "'", Type: calcgo.TInvalidCharacter},
			})
		})
	})

	Convey("Lexer handles whitespace", t, func() {
		Convey("at the beginning", func() {
			So(calcgo.Lex(" 1 + 2"), shouldEqualToken, []calcgo.Token{
				{Value: "1", Type: calcgo.TInteger},
				{Value: "", Type: calcgo.TOperatorPlus},
				{Value: "2", Type: calcgo.TInteger},
			})
			So(calcgo.Lex("   1 + 2"), shouldEqualToken, []calcgo.Token{
				{Value: "1", Type: calcgo.TInteger},
				{Value: "", Type: calcgo.TOperatorPlus},
				{Value: "2", Type: calcgo.TInteger},
			})
		})

		Convey("at the end", func() {
			So(calcgo.Lex("1 + 2 "), shouldEqualToken, []calcgo.Token{
				{Value: "1", Type: calcgo.TInteger},
				{Value: "", Type: calcgo.TOperatorPlus},
				{Value: "2", Type: calcgo.TInteger},
			})
			So(calcgo.Lex("1 + 2     "), shouldEqualToken, []calcgo.Token{
				{Value: "1", Type: calcgo.TInteger},
				{Value: "", Type: calcgo.TOperatorPlus},
				{Value: "2", Type: calcgo.TInteger},
			})
		})

		Convey("multiple whitespace characters", func() {
			So(calcgo.Lex("1  +  2"), shouldEqualToken, []calcgo.Token{
				{Value: "1", Type: calcgo.TInteger},
				{Value: "", Type: calcgo.TOperatorPlus},
				{Value: "2", Type: calcgo.TInteger},
			})
			So(calcgo.Lex("  (  1 +   2 )  * 2 "), shouldEqualToken, []calcgo.Token{
				{Value: "", Type: calcgo.TLeftBracket},
				{Value: "1", Type: calcgo.TInteger},
				{Value: "", Type: calcgo.TOperatorPlus},
				{Value: "2", Type: calcgo.TInteger},
				{Value: "", Type: calcgo.TRightBracket},
				{Value: "", Type: calcgo.TOperatorMult},
				{Value: "2", Type: calcgo.TInteger},
			})
		})
	})

	Convey("Lexer works with variables", t, func() {
		Convey("single letter variables", func() {
			So(calcgo.Lex("a"), shouldEqualToken, []calcgo.Token{
				{Value: "a", Type: calcgo.TVariable},
			})
		})

		Convey("multi letter variables", func() {
			So(calcgo.Lex("ab"), shouldEqualToken, []calcgo.Token{
				{Value: "ab", Type: calcgo.TVariable},
			})

			So(calcgo.Lex("abcdefghiklmnopqrstvxyz"), shouldEqualToken, []calcgo.Token{
				{Value: "abcdefghiklmnopqrstvxyz", Type: calcgo.TVariable},
			})
		})

		Convey("variables in combination with operators and brackets", func() {
			So(calcgo.Lex("a  +  2"), shouldEqualToken, []calcgo.Token{
				{Value: "a", Type: calcgo.TVariable},
				{Value: "", Type: calcgo.TOperatorPlus},
				{Value: "2", Type: calcgo.TInteger},
			})

			So(calcgo.Lex("ab  +  bc"), shouldEqualToken, []calcgo.Token{
				{Value: "ab", Type: calcgo.TVariable},
				{Value: "", Type: calcgo.TOperatorPlus},
				{Value: "bc", Type: calcgo.TVariable},
			})

			So(calcgo.Lex("(a  +  b) - c"), shouldEqualToken, []calcgo.Token{
				{Value: "", Type: calcgo.TLeftBracket},
				{Value: "a", Type: calcgo.TVariable},
				{Value: "", Type: calcgo.TOperatorPlus},
				{Value: "b", Type: calcgo.TVariable},
				{Value: "", Type: calcgo.TRightBracket},
				{Value: "", Type: calcgo.TOperatorMinus},
				{Value: "c", Type: calcgo.TVariable},
			})
		})
	})

	Convey("Lexer works", t, func() {
		lexer := calcgo.NewLexer("1")
		lexer.Start()

		var tokens []calcgo.Token
		for {
			token := lexer.NextToken()
			if token.Type == calcgo.TEOF {
				break
			}
			tokens = append(tokens, token)
		}

		So(tokens, shouldEqualToken, []calcgo.Token{
			{Value: "1", Type: calcgo.TInteger},
		})
	})
}
