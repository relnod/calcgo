package interpreter_test

import (
	"math"
	"testing"

	. "github.com/relnod/calcgo/calcgotest"
	"github.com/relnod/calcgo/interpreter"
	"github.com/relnod/calcgo/parser"
	. "github.com/smartystreets/goconvey/convey"
)

type TestCaseWrapper struct {
	description string
	cases       []TestCase
	wrappers    []TestCaseWrapper
}

type TestCase struct {
	given          string
	expectedValue  float64
	expectedErrors []error
}

var testCases = []TestCaseWrapper{
	{"no input", []TestCase{{"", 0, nil}}, nil},
	{"numbers", nil, []TestCaseWrapper{
		{"positive integer", []TestCase{
			{"1", 1, nil},
			{"123456789", 123456789, nil},
		}, nil},
		{"positive decimals", []TestCase{
			{"1.0", 1, nil},
			{"12345.67890", 12345.67890, nil},
		}, nil},
		{"negativ integers", []TestCase{
			{"-1", -1, nil},
			{"-123456789", -123456789, nil},
		}, nil},
		{"negativ decimals", []TestCase{
			{"-2.0", -2, nil},
			{"-23456.123", -23456.123, nil},
		}, nil},
	}},
	{"operators", nil, []TestCaseWrapper{
		{"addition", []TestCase{
			{"1 + 1", 1 + 1, nil},
			{"1 + 2 + 3", 1 + 2 + 3, nil},
		}, nil},
		{"subtraction", []TestCase{
			{"1 - 1", 1 - 1, nil},
			{"1 - 2 - 3", 1 - 2 - 3, nil},
		}, nil},
		{"multiplication", []TestCase{
			{"1 * 1", 1 * 1, nil},
			{"1 * 2 * 3", 1 * 2 * 3, nil},
		}, nil},
		{"division", []TestCase{
			{"1 / 1", 1 / 1, nil},
			{"1 / 2 / 3", 1.0 / 2.0 / 3.0, nil},
		}, nil},
	}},
	{"operations with negative numbers", nil, []TestCaseWrapper{
		{"left side negative", []TestCase{
			{"-1 + 1", -1 + 1, nil},
			{"-1 - 1", -1 - 1, nil},
			{"-1 * 1", -1 * 1, nil},
			{"-1 / 1", -1 / 1, nil},
		}, nil},
		{"right side negative", []TestCase{
			{"1 + -1", 1 + -1, nil},
			{"1 - -1", 1 - -1, nil},
			{"1 * -1", 1 * -1, nil},
			{"1 / -1", 1 / -1, nil},
		}, nil},
		{"both sides negative", []TestCase{
			{"-1 + -1", -1 + -1, nil},
			{"-1 - -1", -1 - -1, nil},
			{"-1 * -1", -1 * -1, nil},
			{"-1 / -1", -1 / -1, nil},
		}, nil},
	}},
}

type interpretFnc func(string) (float64, []error)

func interpret(s string) (float64, []error) {
	return interpreter.Interpret(s)
}

func interpreterOptimizerDisabled(s string) (float64, []error) {
	i := interpreter.NewInterpreter(s)
	return i.GetResult()
}

func interpreterOptimizerEnabled(s string) (float64, []error) {
	i := interpreter.NewInterpreter(s)
	i.EnableOptimizer()
	return i.GetResult()
}

func handleTestCases(cases []TestCaseWrapper, fnc interpretFnc) {
	for _, wrapper := range cases {
		Convey(wrapper.description, func() {
			if wrapper.cases != nil {
				for _, c := range wrapper.cases {
					result, errors := fnc(c.given)
					So(result, ShouldEqual, c.expectedValue)
					So(errors, ShouldEqualErrors, c.expectedErrors)
				}
			}

			if wrapper.wrappers != nil {
				handleTestCases(wrapper.wrappers, fnc)
			}
		})
	}
}

func TestInterpreter(t *testing.T) {
	Convey("Interpret Spec", t, func() {
		handleTestCases(testCases, interpret)
	})

	Convey("Interpreter Spec (optimizer disabled)", t, func() {
		handleTestCases(testCases, interpreterOptimizerDisabled)
	})

	Convey("Interpreter Spec (optimizer enabled)", t, func() {
		handleTestCases(testCases, interpreterOptimizerEnabled)
	})
}

func SkipTestInterpreter(t *testing.T) {
	Convey("interpreter works with", t, func() {
		Convey("nothing", func() {
			result, errors := interpreter.Interpret("")
			So(result, ShouldEqual, 0)
			So(errors, ShouldBeNil)
		})
		Convey("positive integers", func() {
			result, errors := interpreter.Interpret("1")
			So(result, ShouldEqual, 1)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("12345")
			So(result, ShouldEqual, 12345)
			So(errors, ShouldBeNil)
		})

		Convey("positive decimals", func() {
			result, errors := interpreter.Interpret("1.0")
			So(result, ShouldEqual, 1.0)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("1234.5678")
			So(result, ShouldEqual, 1234.5678)
			So(errors, ShouldBeNil)
		})

		Convey("negativ numbers", func() {
			result, errors := interpreter.Interpret("-1")
			So(result, ShouldEqual, -1)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("-1.1")
			So(result, ShouldEqual, -1.1)
			So(errors, ShouldBeNil)
		})

		Convey("simple additions with integers", func() {
			result, errors := interpreter.Interpret("1 + 1")
			So(result, ShouldEqual, 2)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("3 + 5")
			So(result, ShouldEqual, 3+5)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("1 + 2 + 3 + 4 + 5 + 6")
			So(result, ShouldEqual, 1+2+3+4+5+6)
			So(errors, ShouldBeNil)
		})

		Convey("simple additions with decimals", func() {
			result, errors := interpreter.Interpret("1.2 + 2.4")
			SkipSo(result, ShouldEqual, 1.2+2.4) // @todo: fix rounding error
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("0.7 + 2.4")
			SkipSo(result, ShouldEqual, 0.7+2.4) // @todo: fix rounding error
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("3.5 + 5.1")
			So(result, ShouldEqual, 3.5+5.1)
			So(errors, ShouldBeNil)
		})

		Convey("simple subtractions", func() {
			result, errors := interpreter.Interpret("1 - 1")
			So(result, ShouldEqual, 1-1)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("3 - 5")
			So(result, ShouldEqual, 3-5)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("1 - 2 - 3 - 4 - 5 - 6")
			So(result, ShouldEqual, 1-2-3-4-5-6)
			So(errors, ShouldBeNil)
		})

		Convey("simple multiplications", func() {
			result, errors := interpreter.Interpret("1 * 1")
			So(result, ShouldEqual, 1*1)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("3 * 5")
			So(result, ShouldEqual, 3*5)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("1 * 2 * 3 * 4 * 5 * 6")
			So(result, ShouldEqual, 1*2*3*4*5*6)
			So(errors, ShouldBeNil)
		})

		Convey("simple divisions", func() {
			result, errors := interpreter.Interpret("1 / 1")
			So(result, ShouldEqual, 1/1)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("3 / 5")
			So(result, ShouldEqual, 3.0/5.0)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("1 / 2 / 3 / 4 / 5 / 6")
			So(result, ShouldEqual, 1.0/2.0/3.0/4.0/5.0/6.0)
			So(errors, ShouldBeNil)
		})

		Convey("basic operations with negative numbers", func() {
			result, errors := interpreter.Interpret("-1 + 2")
			So(result, ShouldEqual, -1+2)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("-1 - 2")
			So(result, ShouldEqual, -1-2)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("-1 * 2")
			So(result, ShouldEqual, -1*2)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("-1 / 2")
			So(result, ShouldEqual, -1.0/2.0)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("-1 + -2")
			So(result, ShouldEqual, -1+-2)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("-1 - -2")
			So(result, ShouldEqual, -1 - -2)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("-1 * -2")
			So(result, ShouldEqual, -1*-2)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("-1 / -2")
			So(result, ShouldEqual, -1.0/-2.0)
			So(errors, ShouldBeNil)
		})

		Convey("'multiplication and division before addition and subtraction' rule", func() {
			result, errors := interpreter.Interpret("1 + 2 / 3")
			SkipSo(result, ShouldEqual, 1.0+2.0/3.0) // @todo: fix rounding error
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("1 - 2 / 3")
			SkipSo(result, ShouldEqual, 1.0-2.0/3.0) // @todo: fix rounding error
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("1 + 2 * 3")
			So(result, ShouldEqual, 1.0+2.0*3.0)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("1 - 2 * 3")
			So(result, ShouldEqual, 1.0-2.0*3.0)
			So(errors, ShouldBeNil)
		})

		Convey("brackets", func() {
			result, errors := interpreter.Interpret("(1 + 2) / 3")
			So(result, ShouldEqual, (1.0+2.0)/3.0)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("(1 - 2) / 3")
			So(result, ShouldEqual, (1.0-2.0)/3.0)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("(1 + 2) * 3")
			So(result, ShouldEqual, (1.0+2.0)*3.0)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("(1 - 2) * 3")
			So(result, ShouldEqual, (1.0-2.0)*3.0)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("2 + (1 - 2) / 3")
			So(result, ShouldEqual, 2.0+(1.0-2.0)/3.0)
			So(errors, ShouldBeNil)
		})

		Convey("nested brackets", func() {
			result, errors := interpreter.Interpret("((1 + 2) / 3) + 1")
			So(result, ShouldEqual, ((1.0+2.0)/3.0)+1)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("((2 + 3) / (1 + 2)) * 3")
			So(result, ShouldEqual, ((2.0+3.0)/(1.0+2.0))*3.0)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("(1 - 2) * (3 - 2) / (1 + 4)")
			So(result, ShouldEqual, (1.0-2.0)*(3.0-2.0)/(1.0+4.0))
			So(errors, ShouldBeNil)
		})

		Convey("brackets and 'multiplication and division before addition and subtraction' rule", func() {
			result, errors := interpreter.Interpret("1 + (1 + 2) * 3")
			So(result, ShouldEqual, 1.0+(1.0+2.0)*3.0)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("1 + (1 + 2) / 3")
			So(result, ShouldEqual, 1.0+(1.0+2.0)/3.0)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("1 - (1 + 2) * 3")
			So(result, ShouldEqual, 1.0-(1.0+2.0)*3.0)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("1 - (1 + 2) / 3")
			So(result, ShouldEqual, 1.0-(1.0+2.0)/3.0)
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("(1 + 2) * 3 + (4 - 6 / (5 + 2))")
			So(result, ShouldEqual, (1.0+2.0)*3.0+(4.0-6.0/(5.0+2.0)))
			So(errors, ShouldBeNil)
		})
	})

	Convey("Interpreter handles functions", t, func() {
		Convey("sqrt", func() {
			result, errors := interpreter.Interpret("sqrt(9)")
			So(result, ShouldEqual, math.Sqrt(9))
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("sqrt(3 * 3)")
			So(result, ShouldEqual, math.Sqrt(3*3))
			So(errors, ShouldBeNil)

			result, errors = interpreter.Interpret("1 + sqrt(3 * 3)")
			So(result, ShouldEqual, 1+math.Sqrt(3*3))
			So(errors, ShouldBeNil)
		})
	})

	Convey("interpreter handles variables", t, func() {
		Convey("works with simple variables", func() {
			i := interpreter.NewInterpreter("a")
			i.SetVar("a", 1.0)
			result, errors := i.GetResult()
			So(result, ShouldEqual, 1)
			So(errors, ShouldBeNil)

			i = interpreter.NewInterpreter("1 + a")
			i.SetVar("a", 1.0)
			result, errors = i.GetResult()
			So(result, ShouldEqual, 2)
			So(errors, ShouldBeNil)
		})

		Convey("works with multiple variables", func() {
			i := interpreter.NewInterpreter("a + b")
			i.SetVar("a", 1)
			i.SetVar("b", 2)
			result, errors := i.GetResult()
			So(result, ShouldEqual, 3)
			So(errors, ShouldBeNil)
		})

		Convey("works with reassining variables", func() {
			i := interpreter.NewInterpreter("1 + a")

			i.SetVar("a", 1.0)
			result, errors := i.GetResult()
			So(result, ShouldEqual, 2)
			So(errors, ShouldBeNil)

			i.SetVar("a", 3.0)
			result, errors = i.GetResult()
			So(result, ShouldEqual, 4)
			So(errors, ShouldBeNil)
		})

		Convey("returns error, when not providing variable", func() {
			i := interpreter.NewInterpreter("1 + a")
			result, errors := i.GetResult()
			So(result, ShouldEqual, 0)
			So(errors, ShouldEqualErrors, []error{
				interpreter.ErrorVariableNotDefined,
			})
		})
	})

	Convey("Interpret() works the same as InterpretAST()", t, func() {
		result1, errors1 := interpreter.Interpret("1 + 2")
		result2, err := interpreter.InterpretAST(parser.AST{
			Node: &parser.Node{
				Type:  parser.NAddition,
				Value: "",
				LeftChild: &parser.Node{
					Type:       parser.NInteger,
					Value:      "1",
					LeftChild:  nil,
					RightChild: nil,
				},
				RightChild: &parser.Node{
					Type:       parser.NInteger,
					Value:      "2",
					LeftChild:  nil,
					RightChild: nil,
				},
			},
		})
		errors2 := []error{err}
		if err == nil {
			errors2 = nil
		}
		So(result1, ShouldEqual, result2)
		So(errors1, ShouldEqualErrors, errors2)
	})

	Convey("Interpreter works with optimizer enabled", t, func() {
		Convey("simple number", func() {
			i := interpreter.NewInterpreter("1")
			i.EnableOptimizer()
			result, errors := i.GetResult()
			So(result, ShouldEqual, 1)
			So(errors, ShouldBeNil)
		})

		Convey("simple variable", func() {
			i := interpreter.NewInterpreter("a")
			i.EnableOptimizer()
			i.SetVar("a", 1.0)
			result, errors := i.GetResult()
			So(result, ShouldEqual, 1)
			So(errors, ShouldBeNil)
		})

		Convey("operations", func() {
			Convey("without variables", func() {
				i := interpreter.NewInterpreter("1 + 1")
				i.EnableOptimizer()
				result, errors := i.GetResult()
				So(result, ShouldEqual, 2)
				So(errors, ShouldBeNil)
			})

			Convey("with variables", func() {
				i := interpreter.NewInterpreter("1 + a")
				i.EnableOptimizer()
				i.SetVar("a", 1.0)
				result, errors := i.GetResult()
				So(result, ShouldEqual, 2)
				So(errors, ShouldBeNil)
			})
		})

		Convey("sqrt", func() {
			Convey("without variables", func() {
				i := interpreter.NewInterpreter("sqrt(9)")
				i.EnableOptimizer()
				result, errors := i.GetResult()
				So(result, ShouldEqual, 3)
				So(errors, ShouldBeNil)
			})

			Convey("with variables", func() {
				i := interpreter.NewInterpreter("sqrt(a)")
				i.EnableOptimizer()
				i.SetVar("a", 9.0)
				result, errors := i.GetResult()
				So(result, ShouldEqual, 3)
				So(errors, ShouldBeNil)
			})
		})

		Convey("handles errors correctly", func() {
			Convey("division by 0", func() {
				i := interpreter.NewInterpreter("1 / 0")
				i.EnableOptimizer()
				i.SetVar("a", 1.0)
				result, errors := i.GetResult()
				So(result, ShouldEqual, 0)
				So(errors, ShouldEqualErrors, []error{interpreter.ErrorDivisionByZero})
			})

			Convey("undefined variable", func() {
				i := interpreter.NewInterpreter("a")
				i.EnableOptimizer()
				result, errors := i.GetResult()
				So(result, ShouldEqual, 0)
				So(errors, ShouldEqualErrors, []error{interpreter.ErrorVariableNotDefined})
			})

			Convey("invalid node type", func() {
				i := interpreter.NewInterpreterFromAST(&parser.AST{
					Node: &parser.Node{
						Type:       30000,
						Value:      "a",
						LeftChild:  nil,
						RightChild: nil,
					},
				})
				i.EnableOptimizer()
				result, errors := i.GetResult()
				So(result, ShouldEqual, 0)
				So(errors, ShouldEqualErrors, []error{interpreter.ErrorInvalidNodeType})
			})

			Convey("error occurs not not on first node", func() {
				i := interpreter.NewInterpreterFromAST(&parser.AST{
					Node: &parser.Node{
						Type:  parser.NAddition,
						Value: "",
						LeftChild: &parser.Node{
							Type:       parser.NVariable,
							Value:      "a",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &parser.Node{
							Type:       parser.NInteger,
							Value:      "1",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
				})
				i.EnableOptimizer()
				result, errors := i.GetResult()
				So(result, ShouldEqual, 0)
				So(errors, ShouldEqualErrors, []error{interpreter.ErrorVariableNotDefined})

				i = interpreter.NewInterpreterFromAST(&parser.AST{
					Node: &parser.Node{
						Type:  parser.NAddition,
						Value: "",
						LeftChild: &parser.Node{
							Type:       parser.NInteger,
							Value:      "1",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &parser.Node{
							Type:       parser.NVariable,
							Value:      "a",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
				})
				i.EnableOptimizer()
				result, errors = i.GetResult()
				So(result, ShouldEqual, 0)
				So(errors, ShouldEqualErrors, []error{interpreter.ErrorVariableNotDefined})
			})
		})
	})

	Convey("interpreter returns errors, when parser returned errors", t, func() {
		result, errors := interpreter.Interpret("$")
		So(result, ShouldEqual, 0)
		So(errors, ShouldEqualErrors, []error{
			parser.ErrorExpectedNumberOrVariable,
		})

		result, errors = interpreter.Interpret("1 + #)")
		So(result, ShouldEqual, 0)
		So(errors, ShouldEqualErrors, []error{
			parser.ErrorExpectedNumberOrVariable,
			parser.ErrorUnexpectedClosingBracket,
		})
	})

	Convey("interpreter returns an error when", t, func() {
		Convey("dividing by 0", func() {
			result, errors := interpreter.Interpret("1 / 0")
			So(result, ShouldEqual, 0)
			So(errors, ShouldEqualErrors, []error{
				interpreter.ErrorDivisionByZero,
			})
		})
		Convey("a node child is missing", func() {
			result, errors := interpreter.InterpretAST(parser.AST{
				Node: &parser.Node{
					Type:      parser.NAddition,
					Value:     "",
					LeftChild: nil,
					RightChild: &parser.Node{
						Type:       parser.NInteger,
						Value:      "2",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(result, ShouldEqual, 0)
			So(errors, ShouldEqual, interpreter.ErrorMissingLeftChild)

			result, errors = interpreter.InterpretAST(parser.AST{
				Node: &parser.Node{
					Type:  parser.NAddition,
					Value: "",
					LeftChild: &parser.Node{
						Type:       parser.NInteger,
						Value:      "2",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: nil,
				},
			})
			So(result, ShouldEqual, 0)
			So(errors, ShouldEqual, interpreter.ErrorMissingRightChild)
		})

		Convey("a wrong number is given", func() {
			result, errors := interpreter.InterpretAST(parser.AST{
				Node: &parser.Node{
					Type:       parser.NInteger,
					Value:      "a",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(result, ShouldEqual, 0)
			So(errors, ShouldEqual, interpreter.ErrorInvalidInteger)

			result, errors = interpreter.InterpretAST(parser.AST{
				Node: &parser.Node{
					Type:       parser.NDecimal,
					Value:      "a",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(result, ShouldEqual, 0)
			So(errors, ShouldEqual, interpreter.ErrorInvalidDecimal)
		})

		Convey("an invalid node type is given", func() {
			result, errors := interpreter.InterpretAST(parser.AST{
				Node: &parser.Node{
					Type:  parser.NAddition,
					Value: "",
					LeftChild: &parser.Node{
						Type:       30000,
						Value:      "a",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &parser.Node{
						Type:       parser.NInteger,
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(result, ShouldEqual, 0)
			So(errors, ShouldEqual, interpreter.ErrorInvalidNodeType)

			result, errors = interpreter.InterpretAST(parser.AST{
				Node: &parser.Node{
					Type:  parser.NSubtraction,
					Value: "",
					LeftChild: &parser.Node{
						Type:       30000,
						Value:      "a",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &parser.Node{
						Type:       parser.NInteger,
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(result, ShouldEqual, 0)
			So(errors, ShouldEqual, interpreter.ErrorInvalidNodeType)

			result, errors = interpreter.InterpretAST(parser.AST{
				Node: &parser.Node{
					Type:  parser.NMultiplication,
					Value: "",
					LeftChild: &parser.Node{
						Type:       30000,
						Value:      "a",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &parser.Node{
						Type:       parser.NInteger,
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(result, ShouldEqual, 0)
			So(errors, ShouldEqual, interpreter.ErrorInvalidNodeType)

			result, errors = interpreter.InterpretAST(parser.AST{
				Node: &parser.Node{
					Type:  parser.NDivision,
					Value: "",
					LeftChild: &parser.Node{
						Type:       30000,
						Value:      "a",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &parser.Node{
						Type:       parser.NInteger,
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(result, ShouldEqual, 0)
			So(errors, ShouldEqual, interpreter.ErrorInvalidNodeType)
		})

		Convey("the error doesn't happen on the first node", func() {
			result, errors := interpreter.InterpretAST(parser.AST{
				Node: &parser.Node{
					Type:  parser.NAddition,
					Value: "",
					LeftChild: &parser.Node{
						Type:       parser.NInteger,
						Value:      "a",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &parser.Node{
						Type:       parser.NInteger,
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(result, ShouldEqual, 0)
			So(errors, ShouldEqual, interpreter.ErrorInvalidInteger)

			result, errors = interpreter.InterpretAST(parser.AST{
				Node: &parser.Node{
					Type:  parser.NAddition,
					Value: "",
					LeftChild: &parser.Node{
						Type:       parser.NInteger,
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &parser.Node{
						Type:       parser.NInteger,
						Value:      "a",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(result, ShouldEqual, 0)
			So(errors, ShouldEqual, interpreter.ErrorInvalidInteger)
		})
	})
}
