package interpreter_test

import (
	"math"
	"testing"

	"github.com/relnod/calcgo/interpreter"
	"github.com/relnod/calcgo/interpreter/calculator"
	"github.com/relnod/calcgo/parser"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

type newInterpreterFnc func(string) *interpreter.Interpreter

func TestInterpreter(t *testing.T) {
	Describe("Interpreter", func() {
		DescribeInterpreter(func(str string) *interpreter.Interpreter {
			return interpreter.NewInterpreter(str)
		})
	})

	Describe("Interpreter (Optimized)", func() {
		DescribeInterpreter(func(str string) *interpreter.Interpreter {
			i := interpreter.NewInterpreter(str)
			i.EnableOptimizer()

			return i
		})
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Interpreter Suite")
}

func DescribeInterpreter(newInterpreter newInterpreterFnc) {
	test := func(in string, out float64, errors []error) {
		i := newInterpreter(in)
		result, errs := i.GetResult()
		Ω(result).Should(BeNumerically("==", out))
		Expect(errs).To(Equal(errors))
	}

	Describe("numbers", func() {
		DescribeTable("integer", test,
			Entry("positive 1", "1", 1.0, nil),
			Entry("positive 2", "123456789", 123456789.0, nil),
			Entry("negative 1", "-1", -1.0, nil),
			Entry("negative 2", "-123456789", -123456789.0, nil),
		)

		DescribeTable("decimal", test,
			Entry("positive 1", "1.0", 1.0, nil),
			Entry("positive 2", "12345.6789", 12345.67890, nil),
			Entry("negative 1", "-1.0", -1.0, nil),
			Entry("negative 2", "-12345.6789", -12345.67890, nil),
		)

		DescribeTable("binary", test,
			Entry("positive 1", "0b0", 0.0, nil),
			Entry("positive 2", "0b1", 1.0, nil),
			Entry("positive 3", "0b101", 5.0, nil),
			Entry("negative 1", "-0b1", -1.0, nil),
			Entry("negative 2", "-0b101", -5.0, nil),
		)

		DescribeTable("hex", test,
			Entry("positive 1", "0x1", 1.0, nil),
			Entry("positive 2", "0xA", 10.0, nil),
			Entry("positive 3", "0x1A", 26.0, nil),
			Entry("negative 1", "-0x1", -1.0, nil),
			Entry("negative 2", "-0xA", -10.0, nil),
			Entry("negative 3", "-0x1A", -26.0, nil),
		)

		DescribeTable("exponential", test,
			Entry("positive 1", "1^1", 1.0, nil),
			Entry("positive 2", "2^2", 4.0, nil),
			Entry("negative 1", "-1^1", -1.0, nil),
			Entry("negative 2", "-2^2", 4.0, nil),
			Entry("negative 3", "-2^3", -8.0, nil),
		)

		DescribeTable("errors", test,
			Entry("not a number", "$", 0.0, []error{parser.ErrorExpectedNumberOrVariable}),
		)
	})

	Describe("operators", func() {
		DescribeTable("addition", test,
			Entry("1", "1 + 1", 2.0, nil),
			Entry("2", "1 + 2 + 3", 6.0, nil),
			Entry("negative numbers", "-1 + -2", -3.0, nil),
		)

		DescribeTable("subtraction", test,
			Entry("1", "1 - 1", 0.0, nil),
			Entry("2", "1 - 2 - 3", -4.0, nil),
			Entry("negative numbers", "-1 - -2", 1.0, nil),
		)

		DescribeTable("multiplication", test,
			Entry("1", "1 * 1", 1.0, nil),
			Entry("2", "1 * 2 * 3", 6.0, nil),
			Entry("negative numbers", "-1 * -2", 2.0, nil),
		)

		DescribeTable("division", test,
			Entry("1", "1 / 1", 1.0, nil),
			Entry("2", "1 / 2 / 3", 1.0/2.0/3.0, nil),
			Entry("negative numbers", "-1 / -2", 0.5, nil),
			Entry("error when dividing by 0", "1 / 0", 0.0, []error{calculator.ErrorDivisionByZero}),
		)

		DescribeTable("modulo", test,
			Entry("1", "1 % 1", 0.0, nil),
			Entry("2", "5 % 6", 5.0, nil),
			Entry("3", "12 % 6", 0.0, nil),
			Entry("4", "13 % 6", 1.0, nil),
		)

		DescribeTable("binary or", test,
			Entry("1", "1 | 1", 1.0, nil),
			Entry("2", "5 | 1", 5.0, nil),
		)

		DescribeTable("binary xor", test,
			Entry("1", "1 ^ 0", 1.0, nil),
			Entry("2", "5 ^ 1", 4.0, nil),
		)

		DescribeTable("binary and", test,
			Entry("1", "1 & 1", 1.0, nil),
			Entry("2", "5 & 0", 0.0, nil),
		)

		DescribeTable("'multiplication and division before addition and subtraction' rule", test,
			Entry("addition, then multiplication", "1 + 2 * 3", 7.0, nil),
			Entry("addition, then division", "1 + 4 / 2", 3.0, nil),
			Entry("subtraction, then multiplication", "1 - 2 * 3", -5.0, nil),
			Entry("subtraction, then division", "1 - 4 / 2", -1.0, nil),
		)

		DescribeTable("errors", test,
			Entry("expected operator", "1 $ 1", 0.0, []error{parser.ErrorExpectedOperator}),
			Entry("expected number", "1 + $", 0.0, []error{parser.ErrorExpectedNumberOrVariable}),
		)
	})

	Describe("brackets", func() {
		DescribeTable("with operators", test,
			Entry("simple 1", "(2 + 1) * 3", 9.0, nil),
			Entry("simple 2", "(2 + 1) / 3", 1.0, nil),
			Entry("simple 3", "(2 - 1) * 3", 3.0, nil),
			Entry("simple 4", "(2 - 1) / 3", 1.0/3.0, nil),

			Entry("simple 5", "1 + (1 + 2) * 3", 10.0, nil),
			Entry("simple 6", "1 + (1 + 2) / 3", 2.0, nil),
			Entry("simple 7", "1 - (1 + 2) * 3", -8.0, nil),
			Entry("simple 8", "1 - (1 + 2) / 3", 0.0, nil),

			Entry("nested 1", "((1 + 2) / 3) + 1", 2.0, nil),
			Entry("nested 2", "((2 + 3) / (1 + 2)) * 3", 5.0, nil),
			Entry("nested 3", "(1 - 2) * (3 - 2) / (1 + 4)", -1.0/5.0, nil),
		)

		DescribeTable("errors", test,
			Entry("missing closing paren", "(1 + 2", 0.0, []error{parser.ErrorMissingClosingBracket}),
		)
	})

	Describe("functions", func() {
		DescribeTable("basics", test,
			Entry("simple", "sqrt(1)", math.Sqrt(1), nil),
			Entry("expression in function arg", "sqrt(1 + 1)", math.Sqrt(2), nil),
			Entry("brackets in function arg", "sqrt((1))", math.Sqrt(1), nil),
			Entry("combined with addition", "sqrt((1)) + 1.0", math.Sqrt(1)+1.0, nil),
			Entry("error in function arg", "sqrt(1 / 0)", 0.0, []error{calculator.ErrorDivisionByZero}),
			Entry("undefined function", "bla(1)", 0.0, []error{parser.ErrorUnkownFunction}),
		)

		DescribeTable("functions", test,
			Entry("sqrt", "sqrt(1)", math.Sqrt(1), nil),
			Entry("sin", "sin(1)", math.Sin(1), nil),
			Entry("cos", "cos(1)", math.Cos(1), nil),
			Entry("tan", "tan(1)", math.Tan(1), nil),
		)
	})
	Describe("variables", func() {
		testVar := func(in string, inVars map[string]float64, out float64, errors []error) {
			i := newInterpreter(in)
			for key, val := range inVars {
				i.SetVar(key, val)
			}
			result, errs := i.GetResult()
			Ω(result).Should(BeNumerically("==", out))
			Expect(errs).To(Equal(errors))
		}

		DescribeTable("table", testVar,
			Entry("simple var", "a", map[string]float64{"a": 1.0}, 1.0, nil),
			Entry("operation with var", "1 + a", map[string]float64{"a": 1.0}, 2.0, nil),
			Entry("multiple vars", "a + b", map[string]float64{"a": 1.0, "b": 2.0}, 3.0, nil),
			Entry("var in function", "sqrt(a)", map[string]float64{"a": 1.0}, math.Sqrt(1.0), nil),
			Entry("error when var is not set", "a", map[string]float64{}, 0.0,
				[]error{interpreter.ErrorVariableNotDefined}),
			Entry("erro when dividing by var, which is 0", "1 / a", map[string]float64{"a": 0.0}, 0.0,
				[]error{calculator.ErrorDivisionByZero}),
		)
	})
}

var _ = DescribeTable("InterpretAST()",
	func(in *parser.AST, expOut float64, expErr error) {
		result, err := interpreter.InterpretAST(in)
		Ω(result).Should(BeNumerically("==", expOut))
		if expErr != nil {
			Expect(err).To(Equal(expErr))
		} else {
			Expect(err).To(BeNil())
		}
	},
	Entry("works", &parser.AST{
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
	}, 3.0, nil),

	Entry("errors, with missing left child", &parser.AST{
		Node: &parser.Node{
			Type:      parser.NAdd,
			Value:     "",
			LeftChild: nil,
			RightChild: &parser.Node{
				Type:       parser.NInt,
				Value:      "2",
				LeftChild:  nil,
				RightChild: nil,
			},
		},
	}, 0.0, interpreter.ErrorMissingLeftChild),

	Entry("errors, with missing right child", &parser.AST{
		Node: &parser.Node{
			Type:  parser.NAdd,
			Value: "",
			LeftChild: &parser.Node{
				Type:       parser.NInt,
				Value:      "2",
				LeftChild:  nil,
				RightChild: nil,
			},
			RightChild: nil,
		},
	}, 0.0, interpreter.ErrorMissingRightChild),

	Entry("errors, when a wrong integer is given", &parser.AST{
		Node: &parser.Node{
			Type:       parser.NInt,
			Value:      "a",
			LeftChild:  nil,
			RightChild: nil,
		},
	}, 0.0, calculator.ErrorInvalidInteger),

	Entry("errors, when a wrong decimal is given", &parser.AST{
		Node: &parser.Node{
			Type:       parser.NDec,
			Value:      "a",
			LeftChild:  nil,
			RightChild: nil,
		},
	}, 0.0, calculator.ErrorInvalidDecimal),

	Entry("errors, when an invalid node type is given", &parser.AST{
		Node: &parser.Node{
			Type:  parser.NAdd,
			Value: "",
			LeftChild: &parser.Node{
				Type:       30000,
				Value:      "a",
				LeftChild:  nil,
				RightChild: nil,
			},
			RightChild: &parser.Node{
				Type:       parser.NInt,
				Value:      "1",
				LeftChild:  nil,
				RightChild: nil,
			},
		},
	}, 0.0, interpreter.ErrorInvalidNodeType),

	Entry("errors, when an error happens in a nested node", &parser.AST{
		Node: &parser.Node{
			Type:  parser.NAdd,
			Value: "",
			LeftChild: &parser.Node{
				Type:       parser.NInt,
				Value:      "a",
				LeftChild:  nil,
				RightChild: nil,
			},
			RightChild: &parser.Node{
				Type:       parser.NInt,
				Value:      "1",
				LeftChild:  nil,
				RightChild: nil,
			},
		},
	}, 0.0, calculator.ErrorInvalidInteger),

	Entry("errors, when an error happens in a nested node", &parser.AST{
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
				Value:      "a",
				LeftChild:  nil,
				RightChild: nil,
			},
		},
	}, 0.0, calculator.ErrorInvalidInteger),
)
