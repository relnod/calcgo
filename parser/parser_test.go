package parser_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/relnod/calcgo/parser"
)

func TestParser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Parser Suite")
}

var _ = Describe("Parser", func() {
	test := func(str string, expAST parser.AST, expErrs []error) {
		ast, errs := parser.Parse(str)
		Expect(ast).To(Equal(expAST))
		Expect(errs).To(Equal(expErrs))
	}

	DescribeTable("literals",
		func(str string, nodeType parser.NodeType, expErrs []error) {
			test(str, parser.AST{
				Node: &parser.Node{
					Type:       nodeType,
					Value:      str,
					LeftChild:  nil,
					RightChild: nil,
				},
			}, expErrs)
		},
		Entry("int", "20", parser.NInt, nil),
		Entry("dec", "20.23", parser.NDec, nil),
		Entry("bin", "0b1", parser.NBin, nil),
		Entry("hex", "0x1", parser.NHex, nil),
		Entry("exp", "1^1", parser.NExp, nil),
		PEntry("invalid number", "1#", parser.NInvalidNumber, []error{parser.ErrorExpectedNumberOrVariable}),
		Entry("variable", "a", parser.NVar, nil),
		PEntry("invalid variable", "a#", parser.NInvalidVariable, []error{parser.ErrorExpectedNumberOrVariable}),
	)

	DescribeTable("operators",
		func(op string, nodeType parser.NodeType, expErrs []error) {
			test(fmt.Sprintf("1 %s 2", op), parser.AST{
				Node: &parser.Node{
					Type:  nodeType,
					Value: "",
					LeftChild: &parser.Node{
						Type:       parser.NInt,
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &parser.Node{
						Type:       parser.NInt,
						Value:      "2",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			}, expErrs)
			test(fmt.Sprintf("1 %s 2 %s 3", op, op), parser.AST{
				Node: &parser.Node{
					Type:  nodeType,
					Value: "",
					LeftChild: &parser.Node{
						Type:  nodeType,
						Value: "",
						LeftChild: &parser.Node{
							Type:       parser.NInt,
							Value:      "1",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &parser.Node{
							Type:       parser.NInt,
							Value:      "2",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
					RightChild: &parser.Node{
						Type:       parser.NInt,
						Value:      "3",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			}, expErrs)
		},
		Entry("addition", "+", parser.NAdd, nil),
		Entry("subtraction", "-", parser.NSub, nil),
		Entry("multiplication", "*", parser.NMult, nil),
		Entry("division", "/", parser.NDiv, nil),
		Entry("modulo", "%", parser.NMod, nil),
		Entry("or", "|", parser.NOr, nil),
		Entry("xor", "^", parser.NXor, nil),
		Entry("and", "&", parser.NAnd, nil),
		PEntry("invalid", "{", parser.NInvalidOperator, []error{parser.ErrorExpectedOperator}),
	)

	DescribeTable("functions",
		func(fn string, nodeType parser.NodeType, expErrs []error) {
			test(fn+"(1)", parser.AST{
				Node: &parser.Node{
					Type:  nodeType,
					Value: "",
					LeftChild: &parser.Node{
						Type:       parser.NInt,
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: nil,
				},
			}, expErrs)
		},
		Entry("sqrt", "sqrt", parser.NFnSqrt, nil),
		Entry("sin", "sin", parser.NFnSin, nil),
		Entry("cos", "cos", parser.NFnCos, nil),
		Entry("tan", "tan", parser.NFnTan, nil),
		PEntry("unknown function", "foo", parser.NInvalidFunction, []error{parser.ErrorUnkownFunction}),
		Entry("missing closing paren", "sqrt(", parser.NFnSqrt, []error{parser.ErrorMissingClosingBracket}),
	)

	DescribeTable("'multiplication and division before addition and subtraction' rule 2",
		func(op1, op2 string, nodeType1, nodeType2 parser.NodeType) {
			test(fmt.Sprintf("1 %s 2 %s 3", op1, op2), parser.AST{
				Node: &parser.Node{
					Type:  nodeType1,
					Value: "",
					LeftChild: &parser.Node{
						Type:       parser.NInt,
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &parser.Node{
						Type:  nodeType2,
						Value: "",
						LeftChild: &parser.Node{
							Type:       parser.NInt,
							Value:      "2",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &parser.Node{
							Type:       parser.NInt,
							Value:      "3",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
				},
			}, nil)
		},
		Entry("1", "+", "*", parser.NAdd, parser.NMult),
		Entry("2", "+", "/", parser.NAdd, parser.NDiv),
		Entry("3", "-", "*", parser.NSub, parser.NMult),
		Entry("4", "-", "/", parser.NSub, parser.NDiv),
	)

	DescribeTable("'multiplication and division before addition and subtraction' rule 1",
		func(op1, op2 string, nodeType1, nodeType2 parser.NodeType) {
			test(fmt.Sprintf("1 %s 2 %s 3", op1, op2), parser.AST{
				Node: &parser.Node{
					Type:  nodeType1,
					Value: "",
					LeftChild: &parser.Node{
						Type:  nodeType2,
						Value: "",
						LeftChild: &parser.Node{
							Type:       parser.NInt,
							Value:      "1",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &parser.Node{
							Type:       parser.NInt,
							Value:      "2",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
					RightChild: &parser.Node{
						Type:       parser.NInt,
						Value:      "3",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			}, nil)
		},
		Entry("5", "*", "+", parser.NAdd, parser.NMult),
		Entry("6", "*", "-", parser.NSub, parser.NMult),
		Entry("7", "/", "+", parser.NAdd, parser.NDiv),
		Entry("8", "/", "-", parser.NSub, parser.NDiv),
	)

	DescribeTable("parens", test,
		Entry("surrounding number", "(1)", parser.AST{
			Node: &parser.Node{
				Type:       parser.NInt,
				Value:      "1",
				LeftChild:  nil,
				RightChild: nil,
			},
		}, nil),
		Entry("multiple surrounding number", "((1))", parser.AST{
			Node: &parser.Node{
				Type:       parser.NInt,
				Value:      "1",
				LeftChild:  nil,
				RightChild: nil,
			},
		}, nil),
		Entry("surrounding operatrion", "(1 + 2)", parser.AST{
			Node: &parser.Node{
				Type:  parser.NAdd,
				Value: "",
				LeftChild: &parser.Node{
					Type:       parser.NInt,
					Value:      "1",
					LeftChild:  nil,
					RightChild: nil,
				},
				RightChild: &parser.Node{
					Type:       parser.NInt,
					Value:      "2",
					LeftChild:  nil,
					RightChild: nil,
				},
			},
		}, nil),
		Entry("multiple surrounding operatrion", "((1 + 2))", parser.AST{
			Node: &parser.Node{
				Type:  parser.NAdd,
				Value: "",
				LeftChild: &parser.Node{
					Type:       parser.NInt,
					Value:      "1",
					LeftChild:  nil,
					RightChild: nil,
				},
				RightChild: &parser.Node{
					Type:       parser.NInt,
					Value:      "2",
					LeftChild:  nil,
					RightChild: nil,
				},
			},
		}, nil),
		Entry("missing closing paren", "(1", parser.AST{
			Node: &parser.Node{
				Type:       parser.NInt,
				Value:      "1",
				LeftChild:  nil,
				RightChild: nil,
			},
		}, []error{parser.ErrorMissingClosingBracket}),
		Entry("unexpected closing paren", "1)", parser.AST{
			Node: &parser.Node{
				Type:       parser.NInt,
				Value:      "1",
				LeftChild:  nil,
				RightChild: nil,
			},
		}, []error{parser.ErrorUnexpectedClosingBracket}),
		Entry("unexpected closing paren", "(1))", parser.AST{
			Node: &parser.Node{
				Type:       parser.NInt,
				Value:      "1",
				LeftChild:  nil,
				RightChild: nil,
			},
		}, []error{parser.ErrorUnexpectedClosingBracket}),
	)

	DescribeTable("parens breaking 'multiplication and division before addition and subtraction' rule",
		func(op1, op2 string, nodeType1, nodeType2 parser.NodeType) {
			test(fmt.Sprintf("(1 %s 2) %s 3", op1, op2), parser.AST{
				Node: &parser.Node{
					Type:  nodeType1,
					Value: "",
					LeftChild: &parser.Node{
						Type:  nodeType2,
						Value: "",
						LeftChild: &parser.Node{
							Type:       parser.NInt,
							Value:      "1",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &parser.Node{
							Type:       parser.NInt,
							Value:      "2",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
					RightChild: &parser.Node{
						Type:       parser.NInt,
						Value:      "3",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			}, nil)
		},
		Entry("1", "+", "*", parser.NMult, parser.NAdd),
		Entry("2", "+", "/", parser.NDiv, parser.NAdd),
		Entry("3", "-", "*", parser.NMult, parser.NSub),
		Entry("4", "-", "/", parser.NDiv, parser.NSub),
	)
})
