package token_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/relnod/calcgo/token"
)

func TestToken(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Token Suite")
}

var _ = DescribeTable("IsLiteral()",
	func(tokenType token.Type, expected bool) {
		t := token.Token{Type: tokenType}
		Expect(t.IsLiteral()).To(Equal(expected))
	},
	Entry("1", token.EOF, false),
	Entry("2", token.Int, true),
	Entry("3", token.Dec, true),
	Entry("4", token.Hex, true),
	Entry("5", token.Bin, true),
	Entry("6", token.Exp, true),
	Entry("7", token.Var, true),
	Entry("8", token.Plus, false),
)

var _ = DescribeTable("IsOperator()",
	func(tokenType token.Type, expected bool) {
		t := token.Token{Type: tokenType}
		Expect(t.IsOperator()).To(Equal(expected))
	},
	Entry("1", token.Var, false),
	Entry("2", token.Plus, true),
	Entry("3", token.Minus, true),
	Entry("4", token.Mult, true),
	Entry("5", token.Div, true),
	Entry("6", token.Or, true),
	Entry("7", token.Xor, true),
	Entry("7", token.And, true),
	Entry("8", token.ParenL, false),
)

var _ = DescribeTable("IsFunction()",
	func(tokenType token.Type, expected bool) {
		t := token.Token{Type: tokenType}
		Expect(t.IsFunction()).To(Equal(expected))
	},
	Entry("1", token.Div, false),
	Entry("2", token.Sqrt, true),
	Entry("3", token.Sin, true),
	Entry("4", token.Cos, true),
	Entry("5", token.Tan, true),
	Entry("8", token.ParenL, false),
)

var _ = DescribeTable("String()",
	func(t token.Token, expected string) {
		Expect(t.String()).To(Equal(expected))
	},
	Entry("valid token", token.Token{Type: token.Div, Value: "/"},
		`{Value: '/', Type: '/', Start: '0', End: '0', }`),
	Entry("unkown token", token.Token{Type: 255, Value: "foo"},
		`{Value: 'foo', Type: 'Unknown token', Start: '0', End: '0', }`),
)
