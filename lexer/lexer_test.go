package lexer_test

import (
	"testing"

	"github.com/relnod/calcgo/lexer"
	. "github.com/smartystreets/goconvey/convey"
)

type Tokens []lexer.Token

func (tokens Tokens) String() string {
	str := ""

	for i := 0; i < len(tokens); i++ {
		str += tokens[i].String() + "\n"
	}

	return str
}

func tokenError(actual Tokens, expected Tokens) string {
	return "Expected: \n" +
		expected.String() +
		"Actual: \n" +
		actual.String()
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
	t1 := actual.([]lexer.Token)
	t2 := expected[0].([]lexer.Token)

	if eq(t1, t2) {
		return ""
	}

	return tokenError(t1, t2) + "(Should be Equal)"
}

func shouldNotEqualToken(actual interface{}, expected ...interface{}) string {
	t1 := actual.([]lexer.Token)
	t2 := expected[0].([]lexer.Token)

	if !eq(t1, t2) {
		return ""
	}

	return tokenError(t1, t2) + "(Should not be Equal)"
}

func TestLexer(t *testing.T) {
	Convey("Lexer works with empty string", t, func() {
		So(lexer.Lex(""), ShouldBeNil)
	})

	Convey("Lexer works with numbers", t, func() {
		Convey("single digit", func() {
			So(lexer.Lex("0"), shouldEqualToken, []lexer.Token{
				{Value: "0", Type: lexer.TInteger},
			})
			So(lexer.Lex("1"), shouldEqualToken, []lexer.Token{
				{Value: "1", Type: lexer.TInteger},
			})
			So(lexer.Lex("2"), shouldEqualToken, []lexer.Token{
				{Value: "2", Type: lexer.TInteger},
			})
			So(lexer.Lex("3"), shouldEqualToken, []lexer.Token{
				{Value: "3", Type: lexer.TInteger},
			})
			So(lexer.Lex("4"), shouldEqualToken, []lexer.Token{
				{Value: "4", Type: lexer.TInteger},
			})
			So(lexer.Lex("5"), shouldEqualToken, []lexer.Token{
				{Value: "5", Type: lexer.TInteger},
			})
			So(lexer.Lex("6"), shouldEqualToken, []lexer.Token{
				{Value: "6", Type: lexer.TInteger},
			})
			So(lexer.Lex("7"), shouldEqualToken, []lexer.Token{
				{Value: "7", Type: lexer.TInteger},
			})
			So(lexer.Lex("8"), shouldEqualToken, []lexer.Token{
				{Value: "8", Type: lexer.TInteger},
			})
			So(lexer.Lex("9"), shouldEqualToken, []lexer.Token{
				{Value: "9", Type: lexer.TInteger},
			})
		})

		Convey("multiple digits", func() {
			So(lexer.Lex("10"), shouldEqualToken, []lexer.Token{
				{Value: "10", Type: lexer.TInteger},
			})
			So(lexer.Lex("10123"), shouldEqualToken, []lexer.Token{
				{Value: "10123", Type: lexer.TInteger},
			})
		})

		Convey("decimals", func() {
			So(lexer.Lex("1.0"), shouldEqualToken, []lexer.Token{
				{Value: "1.0", Type: lexer.TDecimal},
			})
			So(lexer.Lex("10.1"), shouldEqualToken, []lexer.Token{
				{Value: "10.1", Type: lexer.TDecimal},
			})
			So(lexer.Lex("12.3456"), shouldEqualToken, []lexer.Token{
				{Value: "12.3456", Type: lexer.TDecimal},
			})
			So(lexer.Lex("0.3456"), shouldEqualToken, []lexer.Token{
				{Value: "0.3456", Type: lexer.TDecimal},
			})
		})

		Convey("negative numbers", func() {
			So(lexer.Lex("-1"), shouldEqualToken, []lexer.Token{
				{Value: "-1", Type: lexer.TInteger},
			})
			So(lexer.Lex("-10"), shouldEqualToken, []lexer.Token{
				{Value: "-10", Type: lexer.TInteger},
			})
			So(lexer.Lex("-10.12"), shouldEqualToken, []lexer.Token{
				{Value: "-10.12", Type: lexer.TDecimal},
			})
			So(lexer.Lex("(-1)"), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TLeftBracket},
				{Value: "-1", Type: lexer.TInteger},
				{Value: "", Type: lexer.TRightBracket},
			})
		})
	})

	Convey("Lexer works with operators", t, func() {
		Convey("plus", func() {
			So(lexer.Lex("+"), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TOperatorPlus},
			})
		})
		Convey("minus", func() {
			So(lexer.Lex("-"), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TOperatorMinus},
			})
		})
		Convey("mult", func() {
			So(lexer.Lex("*"), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TOperatorMult},
			})
		})
		Convey("div", func() {
			So(lexer.Lex("/"), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TOperatorDiv},
			})
		})
	})

	Convey("Lexer works with brackets", t, func() {
		Convey("left bracket", func() {
			So(lexer.Lex("("), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TLeftBracket},
			})
		})
		Convey("right bracket", func() {
			So(lexer.Lex(")"), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TRightBracket},
			})
		})

		Convey("brackets and numbers", func() {
			So(lexer.Lex("(1)"), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TLeftBracket},
				{Value: "1", Type: lexer.TInteger},
				{Value: "", Type: lexer.TRightBracket},
			})
		})
	})

	Convey("Lexer works with mixed types", t, func() {
		So(lexer.Lex("1 + 2"), shouldEqualToken, []lexer.Token{
			{Value: "1", Type: lexer.TInteger},
			{Value: "", Type: lexer.TOperatorPlus},
			{Value: "2", Type: lexer.TInteger},
		})
		So(lexer.Lex("1 + 2 + 3 + 4"), shouldEqualToken, []lexer.Token{
			{Value: "1", Type: lexer.TInteger},
			{Value: "", Type: lexer.TOperatorPlus},
			{Value: "2", Type: lexer.TInteger},
			{Value: "", Type: lexer.TOperatorPlus},
			{Value: "3", Type: lexer.TInteger},
			{Value: "", Type: lexer.TOperatorPlus},
			{Value: "4", Type: lexer.TInteger},
		})
		So(lexer.Lex("(1 + 2) * 2"), shouldEqualToken, []lexer.Token{
			{Value: "", Type: lexer.TLeftBracket},
			{Value: "1", Type: lexer.TInteger},
			{Value: "", Type: lexer.TOperatorPlus},
			{Value: "2", Type: lexer.TInteger},
			{Value: "", Type: lexer.TRightBracket},
			{Value: "", Type: lexer.TOperatorMult},
			{Value: "2", Type: lexer.TInteger},
		})
		So(lexer.Lex("(2 * (1 + 2)) / 2"), shouldEqualToken, []lexer.Token{
			{Value: "", Type: lexer.TLeftBracket},
			{Value: "2", Type: lexer.TInteger},
			{Value: "", Type: lexer.TOperatorMult},
			{Value: "", Type: lexer.TLeftBracket},
			{Value: "1", Type: lexer.TInteger},
			{Value: "", Type: lexer.TOperatorPlus},
			{Value: "2", Type: lexer.TInteger},
			{Value: "", Type: lexer.TRightBracket},
			{Value: "", Type: lexer.TRightBracket},
			{Value: "", Type: lexer.TOperatorDiv},
			{Value: "2", Type: lexer.TInteger},
		})
	})

	Convey("Lexer handles invalid input correctly", t, func() {
		Convey("returns error with", func() {
			Convey("invalid character in number", func() {
				So(lexer.Lex("1%"), shouldEqualToken, []lexer.Token{
					{Value: "%", Type: lexer.TInvalidCharacterInNumber},
				})
				So(lexer.Lex("10123o"), shouldEqualToken, []lexer.Token{
					{Value: "o", Type: lexer.TInvalidCharacterInNumber},
				})
				So(lexer.Lex("10123? "), shouldEqualToken, []lexer.Token{
					{Value: "?", Type: lexer.TInvalidCharacterInNumber},
				})
			})

			Convey("invalid characters", func() {
				So(lexer.Lex("%"), shouldEqualToken, []lexer.Token{
					{Value: "%", Type: lexer.TInvalidCharacter},
				})
				So(lexer.Lex("a$"), shouldEqualToken, []lexer.Token{
					{Value: "$", Type: lexer.TInvalidCharacterInVariable},
				})
				So(lexer.Lex("a1"), shouldEqualToken, []lexer.Token{
					{Value: "1", Type: lexer.TInvalidCharacterInVariable},
				})
				So(lexer.Lex("1 + ~"), shouldEqualToken, []lexer.Token{
					{Value: "1", Type: lexer.TInteger},
					{Value: "", Type: lexer.TOperatorPlus},
					{Value: "~", Type: lexer.TInvalidCharacter},
				})
			})
		})

		Convey("doesn't abort after error", func() {
			So(lexer.Lex("# + 1"), shouldEqualToken, []lexer.Token{
				{Value: "#", Type: lexer.TInvalidCharacter},
				{Value: "", Type: lexer.TOperatorPlus},
				{Value: "1", Type: lexer.TInteger},
			})
		})

		Convey("handles multiple errors", func() {
			So(lexer.Lex("' &"), shouldEqualToken, []lexer.Token{
				{Value: "'", Type: lexer.TInvalidCharacter},
				{Value: "&", Type: lexer.TInvalidCharacter},
			})
			So(lexer.Lex("# + '"), shouldEqualToken, []lexer.Token{
				{Value: "#", Type: lexer.TInvalidCharacter},
				{Value: "", Type: lexer.TOperatorPlus},
				{Value: "'", Type: lexer.TInvalidCharacter},
			})
		})
	})

	Convey("Lexer handles whitespace", t, func() {
		Convey("at the beginning", func() {
			So(lexer.Lex(" 1 + 2"), shouldEqualToken, []lexer.Token{
				{Value: "1", Type: lexer.TInteger},
				{Value: "", Type: lexer.TOperatorPlus},
				{Value: "2", Type: lexer.TInteger},
			})
			So(lexer.Lex("   1 + 2"), shouldEqualToken, []lexer.Token{
				{Value: "1", Type: lexer.TInteger},
				{Value: "", Type: lexer.TOperatorPlus},
				{Value: "2", Type: lexer.TInteger},
			})
		})

		Convey("at the end", func() {
			So(lexer.Lex("1 + 2 "), shouldEqualToken, []lexer.Token{
				{Value: "1", Type: lexer.TInteger},
				{Value: "", Type: lexer.TOperatorPlus},
				{Value: "2", Type: lexer.TInteger},
			})
			So(lexer.Lex("1 + 2     "), shouldEqualToken, []lexer.Token{
				{Value: "1", Type: lexer.TInteger},
				{Value: "", Type: lexer.TOperatorPlus},
				{Value: "2", Type: lexer.TInteger},
			})
		})

		Convey("multiple whitespace characters", func() {
			So(lexer.Lex("1  +  2"), shouldEqualToken, []lexer.Token{
				{Value: "1", Type: lexer.TInteger},
				{Value: "", Type: lexer.TOperatorPlus},
				{Value: "2", Type: lexer.TInteger},
			})
			So(lexer.Lex("  (  1 +   2 )  * 2 "), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TLeftBracket},
				{Value: "1", Type: lexer.TInteger},
				{Value: "", Type: lexer.TOperatorPlus},
				{Value: "2", Type: lexer.TInteger},
				{Value: "", Type: lexer.TRightBracket},
				{Value: "", Type: lexer.TOperatorMult},
				{Value: "2", Type: lexer.TInteger},
			})
		})
	})

	Convey("Lexer works with variables", t, func() {
		Convey("single letter variables", func() {
			So(lexer.Lex("a"), shouldEqualToken, []lexer.Token{
				{Value: "a", Type: lexer.TVariable},
			})
		})

		Convey("multi letter variables", func() {
			So(lexer.Lex("ab"), shouldEqualToken, []lexer.Token{
				{Value: "ab", Type: lexer.TVariable},
			})

			So(lexer.Lex("abcdefghiklmnopqrstvxyz"), shouldEqualToken, []lexer.Token{
				{Value: "abcdefghiklmnopqrstvxyz", Type: lexer.TVariable},
			})
		})

		Convey("variables in combination with operators and brackets", func() {
			So(lexer.Lex("a  +  2"), shouldEqualToken, []lexer.Token{
				{Value: "a", Type: lexer.TVariable},
				{Value: "", Type: lexer.TOperatorPlus},
				{Value: "2", Type: lexer.TInteger},
			})

			So(lexer.Lex("ab  +  bc"), shouldEqualToken, []lexer.Token{
				{Value: "ab", Type: lexer.TVariable},
				{Value: "", Type: lexer.TOperatorPlus},
				{Value: "bc", Type: lexer.TVariable},
			})

			So(lexer.Lex("(a  +  b) - c"), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TLeftBracket},
				{Value: "a", Type: lexer.TVariable},
				{Value: "", Type: lexer.TOperatorPlus},
				{Value: "b", Type: lexer.TVariable},
				{Value: "", Type: lexer.TRightBracket},
				{Value: "", Type: lexer.TOperatorMinus},
				{Value: "c", Type: lexer.TVariable},
			})
		})
	})

	Convey("Lexer works with functions", t, func() {
		Convey("general", func() {
			So(lexer.Lex("sqrt"), shouldEqualToken, []lexer.Token{
				{Value: "sqrt", Type: lexer.TVariable},
			})

			So(lexer.Lex("sqrt("), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TFuncSqrt},
			})

			So(lexer.Lex("sqrt()"), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TFuncSqrt},
				{Value: "", Type: lexer.TRightBracket},
			})

			So(lexer.Lex("sqrt( 1 )"), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TFuncSqrt},
				{Value: "1", Type: lexer.TInteger},
				{Value: "", Type: lexer.TRightBracket},
			})
		})
		Convey("sqrt()", func() {
			So(lexer.Lex("sqrt("), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TFuncSqrt},
			})
		})
		Convey("sin()", func() {
			So(lexer.Lex("sin("), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TFuncSin},
			})
		})
		Convey("cos()", func() {
			So(lexer.Lex("cos("), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TFuncCos},
			})
		})
		Convey("tan()", func() {
			So(lexer.Lex("tan("), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TFuncTan},
			})
		})
	})

	Convey("Lexer works", t, func() {
		l := lexer.NewLexer("1")
		l.Start()

		var tokens []lexer.Token
		for {
			token := l.NextToken()
			if token.Type == lexer.TEOF {
				break
			}
			tokens = append(tokens, token)
		}

		So(tokens, shouldEqualToken, []lexer.Token{
			{Value: "1", Type: lexer.TInteger},
		})
	})
}
