package lexer_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/relnod/calcgo/lexer"
	"github.com/relnod/calcgo/token"
)

func TestLexer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lexer Suite")
}

var _ = Describe("LexString()", func() {
	test := func(in string, expected []token.Token) {
		Expect(lexer.LexString(in)).To(Equal(expected))
		// tokens := lexer.Lex(in)
		// Expect(tokens).To(HaveLen(len(expected)))
		// for i := range tokens {
		// 	Expect(tokens[i]).To(Equal(expected[i]))
		// }
	}

	Describe("numbers", func() {
		DescribeTable("integers", test,
			Entry("single digit 0", "0", []token.Token{{Value: "0", Type: token.Int, Start: 0, End: 1}}),
			Entry("single digit 1", "1", []token.Token{{Value: "1", Type: token.Int, Start: 0, End: 1}}),
			Entry("single digit 2", "2", []token.Token{{Value: "2", Type: token.Int, Start: 0, End: 1}}),
			Entry("single digit 3", "3", []token.Token{{Value: "3", Type: token.Int, Start: 0, End: 1}}),
			Entry("single digit 4", "4", []token.Token{{Value: "4", Type: token.Int, Start: 0, End: 1}}),
			Entry("single digit 5", "5", []token.Token{{Value: "5", Type: token.Int, Start: 0, End: 1}}),
			Entry("single digit 6", "6", []token.Token{{Value: "6", Type: token.Int, Start: 0, End: 1}}),
			Entry("single digit 7", "7", []token.Token{{Value: "7", Type: token.Int, Start: 0, End: 1}}),
			Entry("single digit 8", "8", []token.Token{{Value: "8", Type: token.Int, Start: 0, End: 1}}),
			Entry("single digit 9", "9", []token.Token{{Value: "9", Type: token.Int, Start: 0, End: 1}}),

			Entry("multiple digit 1", "10", []token.Token{
				{Value: "10", Type: token.Int, Start: 0, End: 2},
			}),
			Entry("multiple digit 2", "10123", []token.Token{
				{Value: "10123", Type: token.Int, Start: 0, End: 5},
			}),

			Entry("negative", "-9", []token.Token{{Value: "-9", Type: token.Int, Start: 0, End: 2}}),

			Entry("multiple integers", "1 1", []token.Token{
				{Value: "1", Type: token.Int, Start: 0, End: 1},
				{Value: "1", Type: token.Int, Start: 2, End: 3},
			}),

			Entry("invalid character in number", "1a", []token.Token{
				{Value: "a", Type: token.InvalidCharacterInNumber, Start: 0, End: 2},
			}),
		)

		DescribeTable("decimals", test,
			Entry("1", "1.0", []token.Token{{Value: "1.0", Type: token.Dec, Start: 0, End: 3}}),
			Entry("2", "10.1", []token.Token{{Value: "10.1", Type: token.Dec, Start: 0, End: 4}}),
			Entry("3", "12.3456", []token.Token{{Value: "12.3456", Type: token.Dec, Start: 0, End: 7}}),
			Entry("4", "0.1", []token.Token{{Value: "0.1", Type: token.Dec, Start: 0, End: 3}}),

			Entry("negative", "-9.1", []token.Token{{Value: "-9.1", Type: token.Dec, Start: 0, End: 4}}),

			Entry("multiple decimals", "1.0 1.0", []token.Token{
				{Value: "1.0", Type: token.Dec, Start: 0, End: 3},
				{Value: "1.0", Type: token.Dec, Start: 4, End: 7},
			}),

			Entry("invalid character", "1.a", []token.Token{
				{Value: "a", Type: token.InvalidCharacterInNumber, Start: 0, End: 3},
			}),
		)

		DescribeTable("binary", test,
			Entry("1", "0b1", []token.Token{{Value: "0b1", Type: token.Bin, Start: 0, End: 3}}),
			Entry("2", "0b0101", []token.Token{{Value: "0b0101", Type: token.Bin, Start: 0, End: 6}}),

			Entry("negative", "-0b1", []token.Token{{Value: "-0b1", Type: token.Bin, Start: 0, End: 4}}),

			Entry("multiple numbers", "0b1 1", []token.Token{
				{Value: "0b1", Type: token.Bin, Start: 0, End: 3},
				{Value: "1", Type: token.Int, Start: 4, End: 5},
			}),

			Entry("invalid character", "0b012", []token.Token{
				{Value: "2", Type: token.InvalidCharacterInNumber, Start: 0, End: 5},
			}),
		)

		DescribeTable("hex", test,
			Entry("1", "0x1", []token.Token{{Value: "0x1", Type: token.Hex, Start: 0, End: 3}}),
			Entry("2", "0xA1F9", []token.Token{{Value: "0xA1F9", Type: token.Hex, Start: 0, End: 6}}),

			Entry("negative", "-0x1", []token.Token{{Value: "-0x1", Type: token.Hex, Start: 0, End: 4}}),

			Entry("multiple numbers", "0x1 1", []token.Token{
				{Value: "0x1", Type: token.Hex, Start: 0, End: 3},
				{Value: "1", Type: token.Int, Start: 4, End: 5},
			}),

			Entry("invalid character", "0xN", []token.Token{
				{Value: "N", Type: token.InvalidCharacterInNumber, Start: 0, End: 3},
			}),
		)

		DescribeTable("exponential", test,
			Entry("1", "1^1", []token.Token{{Value: "1^1", Type: token.Exp, Start: 0, End: 3}}),
			Entry("2", "10^1", []token.Token{{Value: "10^1", Type: token.Exp, Start: 0, End: 4}}),
			Entry("3", "1^10", []token.Token{{Value: "1^10", Type: token.Exp, Start: 0, End: 4}}),

			Entry("negative", "-0^1", []token.Token{{Value: "-0^1", Type: token.Exp, Start: 0, End: 4}}),

			Entry("multiple numbers", "1^1 1", []token.Token{
				{Value: "1^1", Type: token.Exp, Start: 0, End: 3},
				{Value: "1", Type: token.Int, Start: 4, End: 5},
			}),

			Entry("invalid character", "0^a", []token.Token{
				{Value: "a", Type: token.InvalidCharacterInNumber, Start: 0, End: 3},
			}),
		)
	})

	DescribeTable("operators", test,
		Entry("plus", "+", []token.Token{{Value: "", Type: token.Plus, Start: 0, End: 1}}),
		Entry("minus", "-", []token.Token{{Value: "", Type: token.Minus, Start: 0, End: 1}}),
		Entry("mult", "*", []token.Token{{Value: "", Type: token.Mult, Start: 0, End: 1}}),
		Entry("div", "/", []token.Token{{Value: "", Type: token.Div, Start: 0, End: 1}}),
		Entry("modulo", "%", []token.Token{{Value: "", Type: token.Mod, Start: 0, End: 1}}),
		Entry("or", "|", []token.Token{{Value: "", Type: token.Or, Start: 0, End: 1}}),
		Entry("xor", "^", []token.Token{{Value: "", Type: token.Xor, Start: 0, End: 1}}),
		Entry("and", "&", []token.Token{{Value: "", Type: token.And, Start: 0, End: 1}}),
	)

	DescribeTable("parens", test,
		Entry("left", "(", []token.Token{{Value: "", Type: token.ParenL, Start: 0, End: 1}}),
		Entry("right", ")", []token.Token{{Value: "", Type: token.ParenR, Start: 0, End: 1}}),
		Entry("works with numbers", "(1)", []token.Token{
			{Value: "", Type: token.ParenL, Start: 0, End: 1},
			{Value: "1", Type: token.Int, Start: 1, End: 2},
			{Value: "", Type: token.ParenR, Start: 2, End: 3},
		}),
	)

	DescribeTable("mixed token types", test,

		Entry("1 + 2", "1 + 2", []token.Token{
			{Value: "1", Type: token.Int, Start: 0, End: 1},
			{Value: "", Type: token.Plus, Start: 2, End: 3},
			{Value: "2", Type: token.Int, Start: 4, End: 5},
		}),
		Entry("1 + 2 + 3 + 4", "1 + 2 + 3 + 4", []token.Token{
			{Value: "1", Type: token.Int, Start: 0, End: 1},
			{Value: "", Type: token.Plus, Start: 2, End: 3},
			{Value: "2", Type: token.Int, Start: 4, End: 5},
			{Value: "", Type: token.Plus, Start: 6, End: 7},
			{Value: "3", Type: token.Int, Start: 8, End: 9},
			{Value: "", Type: token.Plus, Start: 10, End: 11},
			{Value: "4", Type: token.Int, Start: 12, End: 13},
		}),
		Entry("(1 + 2) * 2", "(1 + 2) * 2", []token.Token{
			{Value: "", Type: token.ParenL, Start: 0, End: 1},
			{Value: "1", Type: token.Int, Start: 1, End: 2},
			{Value: "", Type: token.Plus, Start: 3, End: 4},
			{Value: "2", Type: token.Int, Start: 5, End: 6},
			{Value: "", Type: token.ParenR, Start: 6, End: 7},
			{Value: "", Type: token.Mult, Start: 8, End: 9},
			{Value: "2", Type: token.Int, Start: 10, End: 11},
		}),
		Entry("(2 * (1 + 2)) / 2", "(2 * (1 + 2)) / 2", []token.Token{
			{Value: "", Type: token.ParenL, Start: 0, End: 1},
			{Value: "2", Type: token.Int, Start: 1, End: 2},
			{Value: "", Type: token.Mult, Start: 3, End: 4},
			{Value: "", Type: token.ParenL, Start: 5, End: 6},
			{Value: "1", Type: token.Int, Start: 6, End: 7},
			{Value: "", Type: token.Plus, Start: 8, End: 9},
			{Value: "2", Type: token.Int, Start: 10, End: 11},
			{Value: "", Type: token.ParenR, Start: 11, End: 12},
			{Value: "", Type: token.ParenR, Start: 12, End: 13},
			{Value: "", Type: token.Div, Start: 14, End: 15},
			{Value: "2", Type: token.Int, Start: 16, End: 17},
		}),
	)

	DescribeTable("Lexer handles whitespace", test,
		Entry("at the beginning (single)", " 1 + 2", []token.Token{
			{Value: "1", Type: token.Int, Start: 1, End: 2},
			{Value: "", Type: token.Plus, Start: 3, End: 4},
			{Value: "2", Type: token.Int, Start: 5, End: 6},
		}),
		Entry("at the beginning (multiple)", "   1 + 2", []token.Token{
			{Value: "1", Type: token.Int, Start: 3, End: 4},
			{Value: "", Type: token.Plus, Start: 5, End: 6},
			{Value: "2", Type: token.Int, Start: 7, End: 8},
		}),

		Entry("at the end (single)", "1 + 2 ", []token.Token{
			{Value: "1", Type: token.Int, Start: 0, End: 1},
			{Value: "", Type: token.Plus, Start: 2, End: 3},
			{Value: "2", Type: token.Int, Start: 4, End: 5},
		}),
		Entry("at the end (multiplle)", "1 + 2     ", []token.Token{
			{Value: "1", Type: token.Int, Start: 0, End: 1},
			{Value: "", Type: token.Plus, Start: 2, End: 3},
			{Value: "2", Type: token.Int, Start: 4, End: 5},
		}),

		Entry("multiple whitespace characters in between", "1  +  2", []token.Token{
			{Value: "1", Type: token.Int, Start: 0, End: 1},
			{Value: "", Type: token.Plus, Start: 3, End: 4},
			{Value: "2", Type: token.Int, Start: 6, End: 7},
		}),
		Entry("whitespace everywhere", "  (  1 +   2 )  * 2 ", []token.Token{
			{Value: "", Type: token.ParenL, Start: 2, End: 3},
			{Value: "1", Type: token.Int, Start: 5, End: 6},
			{Value: "", Type: token.Plus, Start: 7, End: 8},
			{Value: "2", Type: token.Int, Start: 11, End: 12},
			{Value: "", Type: token.ParenR, Start: 13, End: 14},
			{Value: "", Type: token.Mult, Start: 16, End: 17},
			{Value: "2", Type: token.Int, Start: 18, End: 19},
		}),
	)

	DescribeTable("variables", test,
		Entry("single letter", "a", []token.Token{{Value: "a", Type: token.Var, Start: 0, End: 1}}),

		Entry("multi letter 1", "ab", []token.Token{{Value: "ab", Type: token.Var, Start: 0, End: 2}}),
		Entry("multi letter 2", "abcdefghiklmnopqrstvxyz", []token.Token{
			{Value: "abcdefghiklmnopqrstvxyz", Type: token.Var, Start: 0, End: 23},
		}),

		Entry("works with other tokens 1", "a  +  2", []token.Token{
			{Value: "a", Type: token.Var, Start: 0, End: 1},
			{Value: "", Type: token.Plus, Start: 3, End: 4},
			{Value: "2", Type: token.Int, Start: 6, End: 7},
		}),
		Entry("works with other tokens 1", "ab  +  bc", []token.Token{
			{Value: "ab", Type: token.Var, Start: 0, End: 2},
			{Value: "", Type: token.Plus, Start: 4, End: 5},
			{Value: "bc", Type: token.Var, Start: 7, End: 9},
		}),
		Entry("works with other tokens 1", "(a  +  b) - c", []token.Token{
			{Value: "", Type: token.ParenL, Start: 0, End: 1},
			{Value: "a", Type: token.Var, Start: 1, End: 2},
			{Value: "", Type: token.Plus, Start: 4, End: 5},
			{Value: "b", Type: token.Var, Start: 7, End: 8},
			{Value: "", Type: token.ParenR, Start: 8, End: 9},
			{Value: "", Type: token.Minus, Start: 10, End: 12},
			{Value: "c", Type: token.Var, Start: 12, End: 13},
		}),
	)

	DescribeTable("Lexer works with functions", test,
		Entry("sqrt", "sqrt(", []token.Token{{Value: "", Type: token.Sqrt, Start: 0, End: 5}}),
		Entry("sin", "sin(", []token.Token{{Value: "", Type: token.Sin, Start: 0, End: 4}}),
		Entry("cos", "cos(", []token.Token{{Value: "", Type: token.Cos, Start: 0, End: 4}}),
		Entry("tan", "tan(", []token.Token{{Value: "", Type: token.Tan, Start: 0, End: 4}}),

		Entry("is var without paren", "sqrt", []token.Token{{Value: "sqrt", Type: token.Var, Start: 0, End: 4}}),

		Entry("function with empty body", "sqrt()", []token.Token{
			{Value: "", Type: token.Sqrt, Start: 0, End: 5},
			{Value: "", Type: token.ParenR, Start: 5, End: 6},
		}),

		Entry("function with not empty body", "sqrt( 1 )", []token.Token{
			{Value: "", Type: token.Sqrt, Start: 0, End: 5},
			{Value: "1", Type: token.Int, Start: 6, End: 7},
			{Value: "", Type: token.ParenR, Start: 8, End: 9},
		}),
	)
})
