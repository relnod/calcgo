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
		Convey("positive", func() {
			Convey("single digit", func() {
				So(lexer.Lex("0"), shouldEqualToken, []lexer.Token{
					{Value: "0", Type: lexer.TInt},
				})
				So(lexer.Lex("1"), shouldEqualToken, []lexer.Token{
					{Value: "1", Type: lexer.TInt},
				})
				So(lexer.Lex("2"), shouldEqualToken, []lexer.Token{
					{Value: "2", Type: lexer.TInt},
				})
				So(lexer.Lex("3"), shouldEqualToken, []lexer.Token{
					{Value: "3", Type: lexer.TInt},
				})
				So(lexer.Lex("4"), shouldEqualToken, []lexer.Token{
					{Value: "4", Type: lexer.TInt},
				})
				So(lexer.Lex("5"), shouldEqualToken, []lexer.Token{
					{Value: "5", Type: lexer.TInt},
				})
				So(lexer.Lex("6"), shouldEqualToken, []lexer.Token{
					{Value: "6", Type: lexer.TInt},
				})
				So(lexer.Lex("7"), shouldEqualToken, []lexer.Token{
					{Value: "7", Type: lexer.TInt},
				})
				So(lexer.Lex("8"), shouldEqualToken, []lexer.Token{
					{Value: "8", Type: lexer.TInt},
				})
				So(lexer.Lex("9"), shouldEqualToken, []lexer.Token{
					{Value: "9", Type: lexer.TInt},
				})
			})

			Convey("multiple digits", func() {
				So(lexer.Lex("10"), shouldEqualToken, []lexer.Token{
					{Value: "10", Type: lexer.TInt},
				})
				So(lexer.Lex("10123"), shouldEqualToken, []lexer.Token{
					{Value: "10123", Type: lexer.TInt},
				})
				So(lexer.Lex("1A"), shouldEqualToken, []lexer.Token{
					{Value: "A", Type: lexer.TInvalidCharacterInNumber},
				})
			})

			Convey("decimals", func() {
				So(lexer.Lex("1.0"), shouldEqualToken, []lexer.Token{
					{Value: "1.0", Type: lexer.TDec},
				})
				So(lexer.Lex("10.1"), shouldEqualToken, []lexer.Token{
					{Value: "10.1", Type: lexer.TDec},
				})
				So(lexer.Lex("12.3456"), shouldEqualToken, []lexer.Token{
					{Value: "12.3456", Type: lexer.TDec},
				})
				So(lexer.Lex("0.3456"), shouldEqualToken, []lexer.Token{
					{Value: "0.3456", Type: lexer.TDec},
				})
				So(lexer.Lex("0.1 1"), shouldEqualToken, []lexer.Token{
					{Value: "0.1", Type: lexer.TDec},
					{Value: "1", Type: lexer.TInt},
				})
				So(lexer.Lex("0.1A"), shouldEqualToken, []lexer.Token{
					{Value: "A", Type: lexer.TInvalidCharacterInNumber},
				})
			})

			Convey("binary", func() {
				So(lexer.Lex("0b1"), shouldEqualToken, []lexer.Token{
					{Value: "0b1", Type: lexer.TBin},
				})
				So(lexer.Lex("0b0101"), shouldEqualToken, []lexer.Token{
					{Value: "0b0101", Type: lexer.TBin},
				})
				So(lexer.Lex("0b1 1"), shouldEqualToken, []lexer.Token{
					{Value: "0b1", Type: lexer.TBin},
					{Value: "1", Type: lexer.TInt},
				})
				So(lexer.Lex("0b012"), shouldEqualToken, []lexer.Token{
					{Value: "2", Type: lexer.TInvalidCharacterInNumber},
				})
			})

			Convey("hex", func() {
				So(lexer.Lex("0x1"), shouldEqualToken, []lexer.Token{
					{Value: "0x1", Type: lexer.THex},
				})
				So(lexer.Lex("0xA"), shouldEqualToken, []lexer.Token{
					{Value: "0xA", Type: lexer.THex},
				})
				So(lexer.Lex("0x1A"), shouldEqualToken, []lexer.Token{
					{Value: "0x1A", Type: lexer.THex},
				})
				So(lexer.Lex("0x1A 1"), shouldEqualToken, []lexer.Token{
					{Value: "0x1A", Type: lexer.THex},
					{Value: "1", Type: lexer.TInt},
				})
				So(lexer.Lex("0xN"), shouldEqualToken, []lexer.Token{
					{Value: "N", Type: lexer.TInvalidCharacterInNumber},
				})
			})

			Convey("exponential", func() {
				So(lexer.Lex("1^1"), shouldEqualToken, []lexer.Token{
					{Value: "1^1", Type: lexer.TExp},
				})
				So(lexer.Lex("10^1"), shouldEqualToken, []lexer.Token{
					{Value: "10^1", Type: lexer.TExp},
				})
				So(lexer.Lex("1^10"), shouldEqualToken, []lexer.Token{
					{Value: "1^10", Type: lexer.TExp},
				})
				So(lexer.Lex("10^10"), shouldEqualToken, []lexer.Token{
					{Value: "10^10", Type: lexer.TExp},
				})
				So(lexer.Lex("10^10 1"), shouldEqualToken, []lexer.Token{
					{Value: "10^10", Type: lexer.TExp},
					{Value: "1", Type: lexer.TInt},
				})
				So(lexer.Lex("1^A"), shouldEqualToken, []lexer.Token{
					{Value: "A", Type: lexer.TInvalidCharacterInNumber},
				})
			})

		})

		Convey("negative", func() {
			Convey("integers", func() {
				So(lexer.Lex("-1"), shouldEqualToken, []lexer.Token{
					{Value: "-1", Type: lexer.TInt},
				})
				So(lexer.Lex("-10"), shouldEqualToken, []lexer.Token{
					{Value: "-10", Type: lexer.TInt},
				})
				So(lexer.Lex("(-1)"), shouldEqualToken, []lexer.Token{
					{Value: "", Type: lexer.TLParen},
					{Value: "-1", Type: lexer.TInt},
					{Value: "", Type: lexer.TRParen},
				})
			})
			Convey("decimals", func() {
				So(lexer.Lex("-10.12"), shouldEqualToken, []lexer.Token{
					{Value: "-10.12", Type: lexer.TDec},
				})
				So(lexer.Lex("-0.12"), shouldEqualToken, []lexer.Token{
					{Value: "-0.12", Type: lexer.TDec},
				})
			})
			Convey("hex", func() {
				So(lexer.Lex("-0x1B"), shouldEqualToken, []lexer.Token{
					{Value: "-0x1B", Type: lexer.THex},
				})
			})
			Convey("binary", func() {
				So(lexer.Lex("-0b1"), shouldEqualToken, []lexer.Token{
					{Value: "-0b1", Type: lexer.TBin},
				})
			})
			Convey("exponential", func() {
				So(lexer.Lex("-12^14"), shouldEqualToken, []lexer.Token{
					{Value: "-12^14", Type: lexer.TExp},
				})
			})
		})

	})

	Convey("Lexer works with operators", t, func() {
		Convey("plus", func() {
			So(lexer.Lex("+"), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TOpPlus},
			})
		})
		Convey("minus", func() {
			So(lexer.Lex("-"), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TOpMinus},
			})
		})
		Convey("mult", func() {
			So(lexer.Lex("*"), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TOpMult},
			})
		})
		Convey("div", func() {
			So(lexer.Lex("/"), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TOpDiv},
			})
		})
	})

	Convey("Lexer works with brackets", t, func() {
		Convey("left bracket", func() {
			So(lexer.Lex("("), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TLParen},
			})
		})
		Convey("right bracket", func() {
			So(lexer.Lex(")"), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TRParen},
			})
		})

		Convey("brackets and numbers", func() {
			So(lexer.Lex("(1)"), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TLParen},
				{Value: "1", Type: lexer.TInt},
				{Value: "", Type: lexer.TRParen},
			})
		})
	})

	Convey("Lexer works with mixed types", t, func() {
		So(lexer.Lex("1 + 2"), shouldEqualToken, []lexer.Token{
			{Value: "1", Type: lexer.TInt},
			{Value: "", Type: lexer.TOpPlus},
			{Value: "2", Type: lexer.TInt},
		})
		So(lexer.Lex("1 + 2 + 3 + 4"), shouldEqualToken, []lexer.Token{
			{Value: "1", Type: lexer.TInt},
			{Value: "", Type: lexer.TOpPlus},
			{Value: "2", Type: lexer.TInt},
			{Value: "", Type: lexer.TOpPlus},
			{Value: "3", Type: lexer.TInt},
			{Value: "", Type: lexer.TOpPlus},
			{Value: "4", Type: lexer.TInt},
		})
		So(lexer.Lex("(1 + 2) * 2"), shouldEqualToken, []lexer.Token{
			{Value: "", Type: lexer.TLParen},
			{Value: "1", Type: lexer.TInt},
			{Value: "", Type: lexer.TOpPlus},
			{Value: "2", Type: lexer.TInt},
			{Value: "", Type: lexer.TRParen},
			{Value: "", Type: lexer.TOpMult},
			{Value: "2", Type: lexer.TInt},
		})
		So(lexer.Lex("(2 * (1 + 2)) / 2"), shouldEqualToken, []lexer.Token{
			{Value: "", Type: lexer.TLParen},
			{Value: "2", Type: lexer.TInt},
			{Value: "", Type: lexer.TOpMult},
			{Value: "", Type: lexer.TLParen},
			{Value: "1", Type: lexer.TInt},
			{Value: "", Type: lexer.TOpPlus},
			{Value: "2", Type: lexer.TInt},
			{Value: "", Type: lexer.TRParen},
			{Value: "", Type: lexer.TRParen},
			{Value: "", Type: lexer.TOpDiv},
			{Value: "2", Type: lexer.TInt},
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
				So(lexer.Lex("1.1.1"), shouldEqualToken, []lexer.Token{
					{Value: ".", Type: lexer.TInvalidCharacterInNumber},
					{Value: "1", Type: lexer.TInt},
				})
				So(lexer.Lex("1^1^1"), shouldEqualToken, []lexer.Token{
					{Value: "^", Type: lexer.TInvalidCharacterInNumber},
					{Value: "1", Type: lexer.TInt},
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
					{Value: "1", Type: lexer.TInt},
					{Value: "", Type: lexer.TOpPlus},
					{Value: "~", Type: lexer.TInvalidCharacter},
				})
			})
		})

		Convey("doesn't abort after error", func() {
			So(lexer.Lex("# + 1"), shouldEqualToken, []lexer.Token{
				{Value: "#", Type: lexer.TInvalidCharacter},
				{Value: "", Type: lexer.TOpPlus},
				{Value: "1", Type: lexer.TInt},
			})
		})

		Convey("handles multiple errors", func() {
			So(lexer.Lex("' &"), shouldEqualToken, []lexer.Token{
				{Value: "'", Type: lexer.TInvalidCharacter},
				{Value: "&", Type: lexer.TInvalidCharacter},
			})
			So(lexer.Lex("# + '"), shouldEqualToken, []lexer.Token{
				{Value: "#", Type: lexer.TInvalidCharacter},
				{Value: "", Type: lexer.TOpPlus},
				{Value: "'", Type: lexer.TInvalidCharacter},
			})
		})
	})

	Convey("Lexer handles whitespace", t, func() {
		Convey("at the beginning", func() {
			So(lexer.Lex(" 1 + 2"), shouldEqualToken, []lexer.Token{
				{Value: "1", Type: lexer.TInt},
				{Value: "", Type: lexer.TOpPlus},
				{Value: "2", Type: lexer.TInt},
			})
			So(lexer.Lex("   1 + 2"), shouldEqualToken, []lexer.Token{
				{Value: "1", Type: lexer.TInt},
				{Value: "", Type: lexer.TOpPlus},
				{Value: "2", Type: lexer.TInt},
			})
		})

		Convey("at the end", func() {
			So(lexer.Lex("1 + 2 "), shouldEqualToken, []lexer.Token{
				{Value: "1", Type: lexer.TInt},
				{Value: "", Type: lexer.TOpPlus},
				{Value: "2", Type: lexer.TInt},
			})
			So(lexer.Lex("1 + 2     "), shouldEqualToken, []lexer.Token{
				{Value: "1", Type: lexer.TInt},
				{Value: "", Type: lexer.TOpPlus},
				{Value: "2", Type: lexer.TInt},
			})
		})

		Convey("multiple whitespace characters", func() {
			So(lexer.Lex("1  +  2"), shouldEqualToken, []lexer.Token{
				{Value: "1", Type: lexer.TInt},
				{Value: "", Type: lexer.TOpPlus},
				{Value: "2", Type: lexer.TInt},
			})
			So(lexer.Lex("  (  1 +   2 )  * 2 "), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TLParen},
				{Value: "1", Type: lexer.TInt},
				{Value: "", Type: lexer.TOpPlus},
				{Value: "2", Type: lexer.TInt},
				{Value: "", Type: lexer.TRParen},
				{Value: "", Type: lexer.TOpMult},
				{Value: "2", Type: lexer.TInt},
			})
		})
	})

	Convey("Lexer works with variables", t, func() {
		Convey("single letter variables", func() {
			So(lexer.Lex("a"), shouldEqualToken, []lexer.Token{
				{Value: "a", Type: lexer.TVar},
			})
		})

		Convey("multi letter variables", func() {
			So(lexer.Lex("ab"), shouldEqualToken, []lexer.Token{
				{Value: "ab", Type: lexer.TVar},
			})

			So(lexer.Lex("abcdefghiklmnopqrstvxyz"), shouldEqualToken, []lexer.Token{
				{Value: "abcdefghiklmnopqrstvxyz", Type: lexer.TVar},
			})
		})

		Convey("variables in combination with operators and brackets", func() {
			So(lexer.Lex("a  +  2"), shouldEqualToken, []lexer.Token{
				{Value: "a", Type: lexer.TVar},
				{Value: "", Type: lexer.TOpPlus},
				{Value: "2", Type: lexer.TInt},
			})

			So(lexer.Lex("ab  +  bc"), shouldEqualToken, []lexer.Token{
				{Value: "ab", Type: lexer.TVar},
				{Value: "", Type: lexer.TOpPlus},
				{Value: "bc", Type: lexer.TVar},
			})

			So(lexer.Lex("(a  +  b) - c"), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TLParen},
				{Value: "a", Type: lexer.TVar},
				{Value: "", Type: lexer.TOpPlus},
				{Value: "b", Type: lexer.TVar},
				{Value: "", Type: lexer.TRParen},
				{Value: "", Type: lexer.TOpMinus},
				{Value: "c", Type: lexer.TVar},
			})
		})
	})

	Convey("Lexer works with functions", t, func() {
		Convey("general", func() {
			So(lexer.Lex("sqrt"), shouldEqualToken, []lexer.Token{
				{Value: "sqrt", Type: lexer.TVar},
			})

			So(lexer.Lex("sqrt("), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TFnSqrt},
			})

			So(lexer.Lex("sqrt()"), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TFnSqrt},
				{Value: "", Type: lexer.TRParen},
			})

			So(lexer.Lex("sqrt( 1 )"), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TFnSqrt},
				{Value: "1", Type: lexer.TInt},
				{Value: "", Type: lexer.TRParen},
			})
		})
		Convey("sqrt()", func() {
			So(lexer.Lex("sqrt("), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TFnSqrt},
			})
		})
		Convey("sin()", func() {
			So(lexer.Lex("sin("), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TFnSin},
			})
		})
		Convey("cos()", func() {
			So(lexer.Lex("cos("), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TFnCos},
			})
		})
		Convey("tan()", func() {
			So(lexer.Lex("tan("), shouldEqualToken, []lexer.Token{
				{Value: "", Type: lexer.TFnTan},
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
			{Value: "1", Type: lexer.TInt},
		})
	})
}
