package parser_test

import (
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/relnod/calcgo/parser"
)

var _ = DescribeTable("IsLiteral()",
	func(nodeType parser.NodeType, exp bool) {
		n := parser.Node{Type: nodeType}
		Expect(parser.IsLiteral(&n)).To(Equal(exp))
	},
	Entry("1", parser.NError, false),
	Entry("2", parser.NInt, true),
	Entry("3", parser.NVar, true),
	Entry("4", parser.NAdd, false),
)

var _ = DescribeTable("IsOperator()",
	func(nodeType parser.NodeType, exp bool) {
		n := parser.Node{Type: nodeType}
		Expect(parser.IsOperator(&n)).To(Equal(exp))
	},
	Entry("1", parser.NVar, false),
	Entry("2", parser.NAdd, true),
	Entry("3", parser.NSub, true),
	Entry("4", parser.NMult, true),
	Entry("5", parser.NDiv, true),
	Entry("6", parser.NMod, true),
	Entry("7", parser.NOr, true),
	Entry("8", parser.NXor, true),
	Entry("9", parser.NAnd, true),
	Entry("10", parser.NFnSqrt, false),
)

var _ = DescribeTable("IsFunction()",
	func(nodeType parser.NodeType, exp bool) {
		n := parser.Node{Type: nodeType}
		Expect(parser.IsFunction(&n)).To(Equal(exp))
	},
	Entry("1", parser.NDiv, false),
	Entry("2", parser.NFnSin, true),
	Entry("3", parser.NFnCos, true),
	Entry("4", parser.NFnTan, true),
	Entry("5", parser.NFnSqrt, true),
)
