package lexer_test

import (
	"testing"

	"github.com/relnod/calcgo/lexer"
	"github.com/relnod/calcgo/token"
	. "github.com/smartystreets/goconvey/convey"
)

type Tokens []token.Token

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

func eq(actual Tokens, expected Tokens, checkPosition bool) bool {
	if len(actual) != len(expected) {
		return false
	}

	for i := 0; i < len(actual); i++ {
		if actual[i].Value != expected[i].Value {
			return false
		}

		if actual[i].Type != expected[i].Type {
			return false
		}

		if checkPosition {
			if actual[i].Start != expected[i].Start {
				return false
			}

			if actual[i].End != expected[i].End {
				return false
			}
		}
	}

	return true
}

func shouldEqualToken(actual interface{}, expected ...interface{}) string {
	act := actual.([]token.Token)
	exp := expected[0].([]token.Token)
	checkPosition := false
	if len(expected) == 2 {
		checkPosition = expected[1].(bool)
	}

	if eq(act, exp, checkPosition) {
		return ""
	}

	return tokenError(act, exp) + "(Should be Equal)"
}

func TestLexer(t *testing.T) {
	Convey("Lexer works with empty string", t, func() {
		So(lexer.Lex(""), ShouldBeNil)
	})

	Convey("Lexer works with numbers", t, func() {
		Convey("positive", func() {
			Convey("single digit", func() {
				So(lexer.Lex("0"), shouldEqualToken, []token.Token{
					{Value: "0", Type: token.TInt, Start: 0, End: 1},
				}, true)
				So(lexer.Lex("1"), shouldEqualToken, []token.Token{
					{Value: "1", Type: token.TInt},
				})
				So(lexer.Lex("2"), shouldEqualToken, []token.Token{
					{Value: "2", Type: token.TInt},
				})
				So(lexer.Lex("3"), shouldEqualToken, []token.Token{
					{Value: "3", Type: token.TInt},
				})
				So(lexer.Lex("4"), shouldEqualToken, []token.Token{
					{Value: "4", Type: token.TInt},
				})
				So(lexer.Lex("5"), shouldEqualToken, []token.Token{
					{Value: "5", Type: token.TInt},
				})
				So(lexer.Lex("6"), shouldEqualToken, []token.Token{
					{Value: "6", Type: token.TInt},
				})
				So(lexer.Lex("7"), shouldEqualToken, []token.Token{
					{Value: "7", Type: token.TInt},
				})
				So(lexer.Lex("8"), shouldEqualToken, []token.Token{
					{Value: "8", Type: token.TInt},
				})
				So(lexer.Lex("9"), shouldEqualToken, []token.Token{
					{Value: "9", Type: token.TInt},
				})
			})

			Convey("multiple digits", func() {
				So(lexer.Lex("10"), shouldEqualToken, []token.Token{
					{Value: "10", Type: token.TInt},
				})
				So(lexer.Lex("10123"), shouldEqualToken, []token.Token{
					{Value: "10123", Type: token.TInt},
				})
				So(lexer.Lex("1A"), shouldEqualToken, []token.Token{
					{Value: "A", Type: token.TInvalidCharacterInNumber},
				})
			})

			Convey("decimals", func() {
				So(lexer.Lex("1.0"), shouldEqualToken, []token.Token{
					{Value: "1.0", Type: token.TDec},
				})
				So(lexer.Lex("10.1"), shouldEqualToken, []token.Token{
					{Value: "10.1", Type: token.TDec},
				})
				So(lexer.Lex("12.3456"), shouldEqualToken, []token.Token{
					{Value: "12.3456", Type: token.TDec},
				})
				So(lexer.Lex("0.3456"), shouldEqualToken, []token.Token{
					{Value: "0.3456", Type: token.TDec},
				})
				So(lexer.Lex("0.1 1"), shouldEqualToken, []token.Token{
					{Value: "0.1", Type: token.TDec},
					{Value: "1", Type: token.TInt},
				})
				So(lexer.Lex("0.1A"), shouldEqualToken, []token.Token{
					{Value: "A", Type: token.TInvalidCharacterInNumber},
				})
			})

			Convey("binary", func() {
				So(lexer.Lex("0b1"), shouldEqualToken, []token.Token{
					{Value: "0b1", Type: token.TBin},
				})
				So(lexer.Lex("0b0101"), shouldEqualToken, []token.Token{
					{Value: "0b0101", Type: token.TBin},
				})
				So(lexer.Lex("0b1 1"), shouldEqualToken, []token.Token{
					{Value: "0b1", Type: token.TBin},
					{Value: "1", Type: token.TInt},
				})
				So(lexer.Lex("0b012"), shouldEqualToken, []token.Token{
					{Value: "2", Type: token.TInvalidCharacterInNumber},
				})
			})

			Convey("hex", func() {
				So(lexer.Lex("0x1"), shouldEqualToken, []token.Token{
					{Value: "0x1", Type: token.THex},
				})
				So(lexer.Lex("0xA"), shouldEqualToken, []token.Token{
					{Value: "0xA", Type: token.THex},
				})
				So(lexer.Lex("0x1A"), shouldEqualToken, []token.Token{
					{Value: "0x1A", Type: token.THex},
				})
				So(lexer.Lex("0x1A 1"), shouldEqualToken, []token.Token{
					{Value: "0x1A", Type: token.THex},
					{Value: "1", Type: token.TInt},
				})
				So(lexer.Lex("0xN"), shouldEqualToken, []token.Token{
					{Value: "N", Type: token.TInvalidCharacterInNumber},
				})
			})

			Convey("exponential", func() {
				So(lexer.Lex("1^1"), shouldEqualToken, []token.Token{
					{Value: "1^1", Type: token.TExp},
				})
				So(lexer.Lex("10^1"), shouldEqualToken, []token.Token{
					{Value: "10^1", Type: token.TExp},
				})
				So(lexer.Lex("1^10"), shouldEqualToken, []token.Token{
					{Value: "1^10", Type: token.TExp},
				})
				So(lexer.Lex("10^10"), shouldEqualToken, []token.Token{
					{Value: "10^10", Type: token.TExp},
				})
				So(lexer.Lex("10^10 1"), shouldEqualToken, []token.Token{
					{Value: "10^10", Type: token.TExp},
					{Value: "1", Type: token.TInt},
				})
				So(lexer.Lex("1^A"), shouldEqualToken, []token.Token{
					{Value: "A", Type: token.TInvalidCharacterInNumber},
				})
			})

		})

		Convey("negative", func() {
			Convey("integers", func() {
				So(lexer.Lex("-1"), shouldEqualToken, []token.Token{
					{Value: "-1", Type: token.TInt},
				})
				So(lexer.Lex("-10"), shouldEqualToken, []token.Token{
					{Value: "-10", Type: token.TInt},
				})
				So(lexer.Lex("(-1)"), shouldEqualToken, []token.Token{
					{Value: "", Type: token.TLParen},
					{Value: "-1", Type: token.TInt},
					{Value: "", Type: token.TRParen},
				})
			})
			Convey("decimals", func() {
				So(lexer.Lex("-10.12"), shouldEqualToken, []token.Token{
					{Value: "-10.12", Type: token.TDec},
				})
				So(lexer.Lex("-1.12"), shouldEqualToken, []token.Token{
					{Value: "-1.12", Type: token.TDec},
				})
			})
			Convey("hex", func() {
				So(lexer.Lex("-0x1B"), shouldEqualToken, []token.Token{
					{Value: "-0x1B", Type: token.THex},
				})
			})
			Convey("binary", func() {
				So(lexer.Lex("-0b1"), shouldEqualToken, []token.Token{
					{Value: "-0b1", Type: token.TBin},
				})
			})
			Convey("exponential", func() {
				So(lexer.Lex("-12^14"), shouldEqualToken, []token.Token{
					{Value: "-12^14", Type: token.TExp},
				})
			})
		})

	})

	Convey("Lexer works with operators", t, func() {
		Convey("plus", func() {
			So(lexer.Lex("+"), shouldEqualToken, []token.Token{
				{Value: "", Type: token.TOpPlus},
			})
		})
		Convey("minus", func() {
			So(lexer.Lex("-"), shouldEqualToken, []token.Token{
				{Value: "", Type: token.TOpMinus},
			})
		})
		Convey("mult", func() {
			So(lexer.Lex("*"), shouldEqualToken, []token.Token{
				{Value: "", Type: token.TOpMult},
			})
		})
		Convey("div", func() {
			So(lexer.Lex("/"), shouldEqualToken, []token.Token{
				{Value: "", Type: token.TOpDiv},
			})
		})
		Convey("modulo", func() {
			So(lexer.Lex("%"), shouldEqualToken, []token.Token{
				{Value: "", Type: token.TOpMod},
			})
		})
		Convey("or", func() {
			So(lexer.Lex("|"), shouldEqualToken, []token.Token{
				{Value: "", Type: token.TOpOr},
			})
		})
		Convey("xor", func() {
			So(lexer.Lex("^"), shouldEqualToken, []token.Token{
				{Value: "", Type: token.TOpXor},
			})
		})
		Convey("and", func() {
			So(lexer.Lex("&"), shouldEqualToken, []token.Token{
				{Value: "", Type: token.TOpAnd},
			})
		})
	})

	Convey("Lexer works with brackets", t, func() {
		Convey("left bracket", func() {
			So(lexer.Lex("("), shouldEqualToken, []token.Token{
				{Value: "", Type: token.TLParen},
			})
		})
		Convey("right bracket", func() {
			So(lexer.Lex(")"), shouldEqualToken, []token.Token{
				{Value: "", Type: token.TRParen},
			})
		})

		Convey("brackets and numbers", func() {
			So(lexer.Lex("(1)"), shouldEqualToken, []token.Token{
				{Value: "", Type: token.TLParen},
				{Value: "1", Type: token.TInt},
				{Value: "", Type: token.TRParen},
			})
		})
	})

	Convey("Lexer works with mixed types", t, func() {
		So(lexer.Lex("1 + 2"), shouldEqualToken, []token.Token{
			{Value: "1", Type: token.TInt},
			{Value: "", Type: token.TOpPlus},
			{Value: "2", Type: token.TInt},
		})
		So(lexer.Lex("1 + 2 + 3 + 4"), shouldEqualToken, []token.Token{
			{Value: "1", Type: token.TInt},
			{Value: "", Type: token.TOpPlus},
			{Value: "2", Type: token.TInt},
			{Value: "", Type: token.TOpPlus},
			{Value: "3", Type: token.TInt},
			{Value: "", Type: token.TOpPlus},
			{Value: "4", Type: token.TInt},
		})
		So(lexer.Lex("(1 + 2) * 2"), shouldEqualToken, []token.Token{
			{Value: "", Type: token.TLParen},
			{Value: "1", Type: token.TInt},
			{Value: "", Type: token.TOpPlus},
			{Value: "2", Type: token.TInt},
			{Value: "", Type: token.TRParen},
			{Value: "", Type: token.TOpMult},
			{Value: "2", Type: token.TInt},
		})
		So(lexer.Lex("(2 * (1 + 2)) / 2"), shouldEqualToken, []token.Token{
			{Value: "", Type: token.TLParen},
			{Value: "2", Type: token.TInt},
			{Value: "", Type: token.TOpMult},
			{Value: "", Type: token.TLParen},
			{Value: "1", Type: token.TInt},
			{Value: "", Type: token.TOpPlus},
			{Value: "2", Type: token.TInt},
			{Value: "", Type: token.TRParen},
			{Value: "", Type: token.TRParen},
			{Value: "", Type: token.TOpDiv},
			{Value: "2", Type: token.TInt},
		})
	})

	Convey("Lexer handles invalid input correctly", t, func() {
		Convey("returns error with", func() {
			Convey("invalid character in number", func() {
				So(lexer.Lex("1%"), shouldEqualToken, []token.Token{
					{Value: "%", Type: token.TInvalidCharacterInNumber},
				})
				So(lexer.Lex("10123o"), shouldEqualToken, []token.Token{
					{Value: "o", Type: token.TInvalidCharacterInNumber},
				})
				So(lexer.Lex("10123? "), shouldEqualToken, []token.Token{
					{Value: "?", Type: token.TInvalidCharacterInNumber},
				})
				So(lexer.Lex("1.1.1"), shouldEqualToken, []token.Token{
					{Value: ".", Type: token.TInvalidCharacterInNumber},
					{Value: "1", Type: token.TInt},
				})
				So(lexer.Lex("1^1^1"), shouldEqualToken, []token.Token{
					{Value: "^", Type: token.TInvalidCharacterInNumber},
					{Value: "1", Type: token.TInt},
				})
			})

			Convey("invalid characters", func() {
				So(lexer.Lex("$"), shouldEqualToken, []token.Token{
					{Value: "$", Type: token.TInvalidCharacter},
				})
				So(lexer.Lex("a$"), shouldEqualToken, []token.Token{
					{Value: "$", Type: token.TInvalidCharacterInVariable},
				})
				So(lexer.Lex("a1"), shouldEqualToken, []token.Token{
					{Value: "1", Type: token.TInvalidCharacterInVariable},
				})
				So(lexer.Lex("1 + ~"), shouldEqualToken, []token.Token{
					{Value: "1", Type: token.TInt},
					{Value: "", Type: token.TOpPlus},
					{Value: "~", Type: token.TInvalidCharacter},
				})
			})

			Convey("unkown function", func() {
				So(lexer.Lex("abcdef("), shouldEqualToken, []token.Token{
					{Value: "abcdef(", Type: token.TFnUnkown},
				})
			})
		})

		Convey("doesn't abort after error", func() {
			So(lexer.Lex("# + 1"), shouldEqualToken, []token.Token{
				{Value: "#", Type: token.TInvalidCharacter},
				{Value: "", Type: token.TOpPlus},
				{Value: "1", Type: token.TInt},
			})
		})

		Convey("handles multiple errors", func() {
			So(lexer.Lex("' $"), shouldEqualToken, []token.Token{
				{Value: "'", Type: token.TInvalidCharacter},
				{Value: "$", Type: token.TInvalidCharacter},
			})
			So(lexer.Lex("# + '"), shouldEqualToken, []token.Token{
				{Value: "#", Type: token.TInvalidCharacter},
				{Value: "", Type: token.TOpPlus},
				{Value: "'", Type: token.TInvalidCharacter},
			})
		})
	})

	Convey("Lexer handles whitespace", t, func() {
		Convey("at the beginning", func() {
			So(lexer.Lex(" 1 + 2"), shouldEqualToken, []token.Token{
				{Value: "1", Type: token.TInt},
				{Value: "", Type: token.TOpPlus},
				{Value: "2", Type: token.TInt},
			})
			So(lexer.Lex("   1 + 2"), shouldEqualToken, []token.Token{
				{Value: "1", Type: token.TInt},
				{Value: "", Type: token.TOpPlus},
				{Value: "2", Type: token.TInt},
			})
		})

		Convey("at the end", func() {
			So(lexer.Lex("1 + 2 "), shouldEqualToken, []token.Token{
				{Value: "1", Type: token.TInt},
				{Value: "", Type: token.TOpPlus},
				{Value: "2", Type: token.TInt},
			})
			So(lexer.Lex("1 + 2     "), shouldEqualToken, []token.Token{
				{Value: "1", Type: token.TInt},
				{Value: "", Type: token.TOpPlus},
				{Value: "2", Type: token.TInt},
			})
		})

		Convey("multiple whitespace characters", func() {
			So(lexer.Lex("1  +  2"), shouldEqualToken, []token.Token{
				{Value: "1", Type: token.TInt},
				{Value: "", Type: token.TOpPlus},
				{Value: "2", Type: token.TInt},
			})
			So(lexer.Lex("  (  1 +   2 )  * 2 "), shouldEqualToken, []token.Token{
				{Value: "", Type: token.TLParen},
				{Value: "1", Type: token.TInt},
				{Value: "", Type: token.TOpPlus},
				{Value: "2", Type: token.TInt},
				{Value: "", Type: token.TRParen},
				{Value: "", Type: token.TOpMult},
				{Value: "2", Type: token.TInt},
			})
		})
	})

	Convey("Lexer works with variables", t, func() {
		Convey("single letter variables", func() {
			So(lexer.Lex("a"), shouldEqualToken, []token.Token{
				{Value: "a", Type: token.TVar},
			})
		})

		Convey("multi letter variables", func() {
			So(lexer.Lex("ab"), shouldEqualToken, []token.Token{
				{Value: "ab", Type: token.TVar},
			})

			So(lexer.Lex("abcdefghiklmnopqrstvxyz"), shouldEqualToken, []token.Token{
				{Value: "abcdefghiklmnopqrstvxyz", Type: token.TVar},
			})
		})

		Convey("variables in combination with operators and brackets", func() {
			So(lexer.Lex("a  +  2"), shouldEqualToken, []token.Token{
				{Value: "a", Type: token.TVar},
				{Value: "", Type: token.TOpPlus},
				{Value: "2", Type: token.TInt},
			})

			So(lexer.Lex("ab  +  bc"), shouldEqualToken, []token.Token{
				{Value: "ab", Type: token.TVar},
				{Value: "", Type: token.TOpPlus},
				{Value: "bc", Type: token.TVar},
			})

			So(lexer.Lex("(a  +  b) - c"), shouldEqualToken, []token.Token{
				{Value: "", Type: token.TLParen},
				{Value: "a", Type: token.TVar},
				{Value: "", Type: token.TOpPlus},
				{Value: "b", Type: token.TVar},
				{Value: "", Type: token.TRParen},
				{Value: "", Type: token.TOpMinus},
				{Value: "c", Type: token.TVar},
			})
		})
	})

	Convey("Lexer works with functions", t, func() {
		Convey("general", func() {
			So(lexer.Lex("sqrt"), shouldEqualToken, []token.Token{
				{Value: "sqrt", Type: token.TVar},
			})

			So(lexer.Lex("sqrt("), shouldEqualToken, []token.Token{
				{Value: "", Type: token.TFnSqrt},
			})

			So(lexer.Lex("sqrt()"), shouldEqualToken, []token.Token{
				{Value: "", Type: token.TFnSqrt},
				{Value: "", Type: token.TRParen},
			})

			So(lexer.Lex("sqrt( 1 )"), shouldEqualToken, []token.Token{
				{Value: "", Type: token.TFnSqrt},
				{Value: "1", Type: token.TInt},
				{Value: "", Type: token.TRParen},
			})
		})
		Convey("sqrt()", func() {
			So(lexer.Lex("sqrt("), shouldEqualToken, []token.Token{
				{Value: "", Type: token.TFnSqrt},
			})
		})
		Convey("sin()", func() {
			So(lexer.Lex("sin("), shouldEqualToken, []token.Token{
				{Value: "", Type: token.TFnSin},
			})
		})
		Convey("cos()", func() {
			So(lexer.Lex("cos("), shouldEqualToken, []token.Token{
				{Value: "", Type: token.TFnCos},
			})
		})
		Convey("tan()", func() {
			So(lexer.Lex("tan("), shouldEqualToken, []token.Token{
				{Value: "", Type: token.TFnTan},
			})
		})
	})

	Convey("Lexer works", t, func() {
		l := lexer.NewLexer("1")
		l.Start()

		var tokens []token.Token
		for {
			t := l.NextToken()
			if t.Type == token.TEOF {
				break
			}
			tokens = append(tokens, t)
		}

		So(tokens, shouldEqualToken, []token.Token{
			{Value: "1", Type: token.TInt, Start: 0, End: 1},
		}, true)
	})
}
