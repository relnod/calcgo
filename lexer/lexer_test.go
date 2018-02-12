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
					{Value: "0", Type: token.Int, Start: 0, End: 1},
				}, true)
				So(lexer.Lex("1"), shouldEqualToken, []token.Token{
					{Value: "1", Type: token.Int},
				})
				So(lexer.Lex("2"), shouldEqualToken, []token.Token{
					{Value: "2", Type: token.Int},
				})
				So(lexer.Lex("3"), shouldEqualToken, []token.Token{
					{Value: "3", Type: token.Int},
				})
				So(lexer.Lex("4"), shouldEqualToken, []token.Token{
					{Value: "4", Type: token.Int},
				})
				So(lexer.Lex("5"), shouldEqualToken, []token.Token{
					{Value: "5", Type: token.Int},
				})
				So(lexer.Lex("6"), shouldEqualToken, []token.Token{
					{Value: "6", Type: token.Int},
				})
				So(lexer.Lex("7"), shouldEqualToken, []token.Token{
					{Value: "7", Type: token.Int},
				})
				So(lexer.Lex("8"), shouldEqualToken, []token.Token{
					{Value: "8", Type: token.Int},
				})
				So(lexer.Lex("9"), shouldEqualToken, []token.Token{
					{Value: "9", Type: token.Int},
				})
			})

			Convey("multiple digits", func() {
				So(lexer.Lex("10"), shouldEqualToken, []token.Token{
					{Value: "10", Type: token.Int},
				})
				So(lexer.Lex("10123"), shouldEqualToken, []token.Token{
					{Value: "10123", Type: token.Int},
				})
				So(lexer.Lex("1A"), shouldEqualToken, []token.Token{
					{Value: "A", Type: token.InvalidCharacterInNumber},
				})
			})

			Convey("decimals", func() {
				So(lexer.Lex("1.0"), shouldEqualToken, []token.Token{
					{Value: "1.0", Type: token.Dec},
				})
				So(lexer.Lex("10.1"), shouldEqualToken, []token.Token{
					{Value: "10.1", Type: token.Dec},
				})
				So(lexer.Lex("12.3456"), shouldEqualToken, []token.Token{
					{Value: "12.3456", Type: token.Dec},
				})
				So(lexer.Lex("0.3456"), shouldEqualToken, []token.Token{
					{Value: "0.3456", Type: token.Dec},
				})
				So(lexer.Lex("0.1 1"), shouldEqualToken, []token.Token{
					{Value: "0.1", Type: token.Dec},
					{Value: "1", Type: token.Int},
				})
				So(lexer.Lex("0.1A"), shouldEqualToken, []token.Token{
					{Value: "A", Type: token.InvalidCharacterInNumber},
				})
			})

			Convey("binary", func() {
				So(lexer.Lex("0b1"), shouldEqualToken, []token.Token{
					{Value: "0b1", Type: token.Bin},
				})
				So(lexer.Lex("0b0101"), shouldEqualToken, []token.Token{
					{Value: "0b0101", Type: token.Bin},
				})
				So(lexer.Lex("0b1 1"), shouldEqualToken, []token.Token{
					{Value: "0b1", Type: token.Bin},
					{Value: "1", Type: token.Int},
				})
				So(lexer.Lex("0b012"), shouldEqualToken, []token.Token{
					{Value: "2", Type: token.InvalidCharacterInNumber},
				})
			})

			Convey("hex", func() {
				So(lexer.Lex("0x1"), shouldEqualToken, []token.Token{
					{Value: "0x1", Type: token.Hex},
				})
				So(lexer.Lex("0xA"), shouldEqualToken, []token.Token{
					{Value: "0xA", Type: token.Hex},
				})
				So(lexer.Lex("0x1A"), shouldEqualToken, []token.Token{
					{Value: "0x1A", Type: token.Hex},
				})
				So(lexer.Lex("0x1A 1"), shouldEqualToken, []token.Token{
					{Value: "0x1A", Type: token.Hex},
					{Value: "1", Type: token.Int},
				})
				So(lexer.Lex("0xN"), shouldEqualToken, []token.Token{
					{Value: "N", Type: token.InvalidCharacterInNumber},
				})
			})

			Convey("exponential", func() {
				So(lexer.Lex("1^1"), shouldEqualToken, []token.Token{
					{Value: "1^1", Type: token.Exp},
				})
				So(lexer.Lex("10^1"), shouldEqualToken, []token.Token{
					{Value: "10^1", Type: token.Exp},
				})
				So(lexer.Lex("1^10"), shouldEqualToken, []token.Token{
					{Value: "1^10", Type: token.Exp},
				})
				So(lexer.Lex("10^10"), shouldEqualToken, []token.Token{
					{Value: "10^10", Type: token.Exp},
				})
				So(lexer.Lex("10^10 1"), shouldEqualToken, []token.Token{
					{Value: "10^10", Type: token.Exp},
					{Value: "1", Type: token.Int},
				})
				So(lexer.Lex("1^A"), shouldEqualToken, []token.Token{
					{Value: "A", Type: token.InvalidCharacterInNumber},
				})
			})

		})

		Convey("negative", func() {
			Convey("integers", func() {
				So(lexer.Lex("-1"), shouldEqualToken, []token.Token{
					{Value: "-1", Type: token.Int},
				})
				So(lexer.Lex("-10"), shouldEqualToken, []token.Token{
					{Value: "-10", Type: token.Int},
				})
				So(lexer.Lex("(-1)"), shouldEqualToken, []token.Token{
					{Value: "", Type: token.ParenL},
					{Value: "-1", Type: token.Int},
					{Value: "", Type: token.ParenR},
				})
			})
			Convey("decimals", func() {
				So(lexer.Lex("-10.12"), shouldEqualToken, []token.Token{
					{Value: "-10.12", Type: token.Dec},
				})
				So(lexer.Lex("-1.12"), shouldEqualToken, []token.Token{
					{Value: "-1.12", Type: token.Dec},
				})
			})
			Convey("hex", func() {
				So(lexer.Lex("-0x1B"), shouldEqualToken, []token.Token{
					{Value: "-0x1B", Type: token.Hex},
				})
			})
			Convey("binary", func() {
				So(lexer.Lex("-0b1"), shouldEqualToken, []token.Token{
					{Value: "-0b1", Type: token.Bin},
				})
			})
			Convey("exponential", func() {
				So(lexer.Lex("-12^14"), shouldEqualToken, []token.Token{
					{Value: "-12^14", Type: token.Exp},
				})
			})
		})

	})

	Convey("Lexer works with operators", t, func() {
		Convey("plus", func() {
			So(lexer.Lex("+"), shouldEqualToken, []token.Token{
				{Value: "", Type: token.Plus},
			})
		})
		Convey("minus", func() {
			So(lexer.Lex("-"), shouldEqualToken, []token.Token{
				{Value: "", Type: token.Minus},
			})
		})
		Convey("mult", func() {
			So(lexer.Lex("*"), shouldEqualToken, []token.Token{
				{Value: "", Type: token.Mult},
			})
		})
		Convey("div", func() {
			So(lexer.Lex("/"), shouldEqualToken, []token.Token{
				{Value: "", Type: token.Div},
			})
		})
		Convey("modulo", func() {
			So(lexer.Lex("%"), shouldEqualToken, []token.Token{
				{Value: "", Type: token.Mod},
			})
		})
		Convey("or", func() {
			So(lexer.Lex("|"), shouldEqualToken, []token.Token{
				{Value: "", Type: token.Or},
			})
		})
		Convey("xor", func() {
			So(lexer.Lex("^"), shouldEqualToken, []token.Token{
				{Value: "", Type: token.Xor},
			})
		})
		Convey("and", func() {
			So(lexer.Lex("&"), shouldEqualToken, []token.Token{
				{Value: "", Type: token.And},
			})
		})
	})

	Convey("Lexer works with brackets", t, func() {
		Convey("left bracket", func() {
			So(lexer.Lex("("), shouldEqualToken, []token.Token{
				{Value: "", Type: token.ParenL},
			})
		})
		Convey("right bracket", func() {
			So(lexer.Lex(")"), shouldEqualToken, []token.Token{
				{Value: "", Type: token.ParenR},
			})
		})

		Convey("brackets and numbers", func() {
			So(lexer.Lex("(1)"), shouldEqualToken, []token.Token{
				{Value: "", Type: token.ParenL},
				{Value: "1", Type: token.Int},
				{Value: "", Type: token.ParenR},
			})
		})
	})

	Convey("Lexer works with mixed types", t, func() {
		So(lexer.Lex("1 + 2"), shouldEqualToken, []token.Token{
			{Value: "1", Type: token.Int},
			{Value: "", Type: token.Plus},
			{Value: "2", Type: token.Int},
		})
		So(lexer.Lex("1 + 2 + 3 + 4"), shouldEqualToken, []token.Token{
			{Value: "1", Type: token.Int},
			{Value: "", Type: token.Plus},
			{Value: "2", Type: token.Int},
			{Value: "", Type: token.Plus},
			{Value: "3", Type: token.Int},
			{Value: "", Type: token.Plus},
			{Value: "4", Type: token.Int},
		})
		So(lexer.Lex("(1 + 2) * 2"), shouldEqualToken, []token.Token{
			{Value: "", Type: token.ParenL},
			{Value: "1", Type: token.Int},
			{Value: "", Type: token.Plus},
			{Value: "2", Type: token.Int},
			{Value: "", Type: token.ParenR},
			{Value: "", Type: token.Mult},
			{Value: "2", Type: token.Int},
		})
		So(lexer.Lex("(2 * (1 + 2)) / 2"), shouldEqualToken, []token.Token{
			{Value: "", Type: token.ParenL},
			{Value: "2", Type: token.Int},
			{Value: "", Type: token.Mult},
			{Value: "", Type: token.ParenL},
			{Value: "1", Type: token.Int},
			{Value: "", Type: token.Plus},
			{Value: "2", Type: token.Int},
			{Value: "", Type: token.ParenR},
			{Value: "", Type: token.ParenR},
			{Value: "", Type: token.Div},
			{Value: "2", Type: token.Int},
		})
	})

	Convey("Lexer handles invalid input correctly", t, func() {
		Convey("returns error with", func() {
			Convey("invalid character in number", func() {
				So(lexer.Lex("1%"), shouldEqualToken, []token.Token{
					{Value: "%", Type: token.InvalidCharacterInNumber},
				})
				So(lexer.Lex("10123o"), shouldEqualToken, []token.Token{
					{Value: "o", Type: token.InvalidCharacterInNumber},
				})
				So(lexer.Lex("10123? "), shouldEqualToken, []token.Token{
					{Value: "?", Type: token.InvalidCharacterInNumber},
				})
				So(lexer.Lex("1.1.1"), shouldEqualToken, []token.Token{
					{Value: ".", Type: token.InvalidCharacterInNumber},
					{Value: "1", Type: token.Int},
				})
				So(lexer.Lex("1^1^1"), shouldEqualToken, []token.Token{
					{Value: "^", Type: token.InvalidCharacterInNumber},
					{Value: "1", Type: token.Int},
				})
			})

			Convey("invalid characters", func() {
				So(lexer.Lex("$"), shouldEqualToken, []token.Token{
					{Value: "$", Type: token.InvalidCharacter},
				})
				So(lexer.Lex("a$"), shouldEqualToken, []token.Token{
					{Value: "$", Type: token.InvalidCharacterInVariable},
				})
				So(lexer.Lex("a1"), shouldEqualToken, []token.Token{
					{Value: "1", Type: token.InvalidCharacterInVariable},
				})
				So(lexer.Lex("1 + ~"), shouldEqualToken, []token.Token{
					{Value: "1", Type: token.Int},
					{Value: "", Type: token.Plus},
					{Value: "~", Type: token.InvalidCharacter},
				})
			})

			Convey("unkown function", func() {
				So(lexer.Lex("abcdef("), shouldEqualToken, []token.Token{
					{Value: "abcdef(", Type: token.UnkownFunktion},
				})
			})
		})

		Convey("doesn't abort after error", func() {
			So(lexer.Lex("# + 1"), shouldEqualToken, []token.Token{
				{Value: "#", Type: token.InvalidCharacter},
				{Value: "", Type: token.Plus},
				{Value: "1", Type: token.Int},
			})
		})

		Convey("handles multiple errors", func() {
			So(lexer.Lex("' $"), shouldEqualToken, []token.Token{
				{Value: "'", Type: token.InvalidCharacter},
				{Value: "$", Type: token.InvalidCharacter},
			})
			So(lexer.Lex("# + '"), shouldEqualToken, []token.Token{
				{Value: "#", Type: token.InvalidCharacter},
				{Value: "", Type: token.Plus},
				{Value: "'", Type: token.InvalidCharacter},
			})
		})
	})

	Convey("Lexer handles whitespace", t, func() {
		Convey("at the beginning", func() {
			So(lexer.Lex(" 1 + 2"), shouldEqualToken, []token.Token{
				{Value: "1", Type: token.Int},
				{Value: "", Type: token.Plus},
				{Value: "2", Type: token.Int},
			})
			So(lexer.Lex("   1 + 2"), shouldEqualToken, []token.Token{
				{Value: "1", Type: token.Int},
				{Value: "", Type: token.Plus},
				{Value: "2", Type: token.Int},
			})
		})

		Convey("at the end", func() {
			So(lexer.Lex("1 + 2 "), shouldEqualToken, []token.Token{
				{Value: "1", Type: token.Int},
				{Value: "", Type: token.Plus},
				{Value: "2", Type: token.Int},
			})
			So(lexer.Lex("1 + 2     "), shouldEqualToken, []token.Token{
				{Value: "1", Type: token.Int},
				{Value: "", Type: token.Plus},
				{Value: "2", Type: token.Int},
			})
		})

		Convey("multiple whitespace characters", func() {
			So(lexer.Lex("1  +  2"), shouldEqualToken, []token.Token{
				{Value: "1", Type: token.Int},
				{Value: "", Type: token.Plus},
				{Value: "2", Type: token.Int},
			})
			So(lexer.Lex("  (  1 +   2 )  * 2 "), shouldEqualToken, []token.Token{
				{Value: "", Type: token.ParenL},
				{Value: "1", Type: token.Int},
				{Value: "", Type: token.Plus},
				{Value: "2", Type: token.Int},
				{Value: "", Type: token.ParenR},
				{Value: "", Type: token.Mult},
				{Value: "2", Type: token.Int},
			})
		})
	})

	Convey("Lexer works with variables", t, func() {
		Convey("single letter variables", func() {
			So(lexer.Lex("a"), shouldEqualToken, []token.Token{
				{Value: "a", Type: token.Var},
			})
		})

		Convey("multi letter variables", func() {
			So(lexer.Lex("ab"), shouldEqualToken, []token.Token{
				{Value: "ab", Type: token.Var},
			})

			So(lexer.Lex("abcdefghiklmnopqrstvxyz"), shouldEqualToken, []token.Token{
				{Value: "abcdefghiklmnopqrstvxyz", Type: token.Var},
			})
		})

		Convey("variables in combination with operators and brackets", func() {
			So(lexer.Lex("a  +  2"), shouldEqualToken, []token.Token{
				{Value: "a", Type: token.Var},
				{Value: "", Type: token.Plus},
				{Value: "2", Type: token.Int},
			})

			So(lexer.Lex("ab  +  bc"), shouldEqualToken, []token.Token{
				{Value: "ab", Type: token.Var},
				{Value: "", Type: token.Plus},
				{Value: "bc", Type: token.Var},
			})

			So(lexer.Lex("(a  +  b) - c"), shouldEqualToken, []token.Token{
				{Value: "", Type: token.ParenL},
				{Value: "a", Type: token.Var},
				{Value: "", Type: token.Plus},
				{Value: "b", Type: token.Var},
				{Value: "", Type: token.ParenR},
				{Value: "", Type: token.Minus},
				{Value: "c", Type: token.Var},
			})
		})
	})

	Convey("Lexer works with functions", t, func() {
		Convey("general", func() {
			So(lexer.Lex("sqrt"), shouldEqualToken, []token.Token{
				{Value: "sqrt", Type: token.Var},
			})

			So(lexer.Lex("sqrt("), shouldEqualToken, []token.Token{
				{Value: "", Type: token.Sqrt},
			})

			So(lexer.Lex("sqrt()"), shouldEqualToken, []token.Token{
				{Value: "", Type: token.Sqrt},
				{Value: "", Type: token.ParenR},
			})

			So(lexer.Lex("sqrt( 1 )"), shouldEqualToken, []token.Token{
				{Value: "", Type: token.Sqrt},
				{Value: "1", Type: token.Int},
				{Value: "", Type: token.ParenR},
			})
		})
		Convey("sqrt()", func() {
			So(lexer.Lex("sqrt("), shouldEqualToken, []token.Token{
				{Value: "", Type: token.Sqrt},
			})
		})
		Convey("sin()", func() {
			So(lexer.Lex("sin("), shouldEqualToken, []token.Token{
				{Value: "", Type: token.Sin},
			})
		})
		Convey("cos()", func() {
			So(lexer.Lex("cos("), shouldEqualToken, []token.Token{
				{Value: "", Type: token.Cos},
			})
		})
		Convey("tan()", func() {
			So(lexer.Lex("tan("), shouldEqualToken, []token.Token{
				{Value: "", Type: token.Tan},
			})
		})
	})

	Convey("Lexer works", t, func() {
		l := lexer.NewLexer("1")
		l.Start()

		var tokens []token.Token
		for {
			t := l.NextToken()
			if t.Type == token.EOF {
				break
			}
			tokens = append(tokens, t)
		}

		So(tokens, shouldEqualToken, []token.Token{
			{Value: "1", Type: token.Int, Start: 0, End: 1},
		}, true)
	})
}
