package calcgo_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/relnod/calcgo"
)

func errorsToString(errors []error) string {
	var str string

	str += "(\n"
	for _, err := range errors {
		str += err.Error() + "\n"
	}
	str += ")\n"

	return str
}

func errorsError(actual []error, expected []error) string {
	return "Expected: \n" +
		errorsToString(expected) +
		"Actual: \n" +
		errorsToString(actual)
}

func eqErrors(e1 []error, e2 []error) bool {
	if len(e1) != len(e2) {
		return false
	}

	for i := 0; i < len(e1); i++ {
		if e1[i] != e2[i] {
			return false
		}
	}

	return true
}

func ShouldEqualErrors(actual interface{}, expected ...interface{}) string {
	e1 := actual.([]error)
	e2 := expected[0].([]error)

	if eqErrors(e1, e2) {
		return ""
	}

	return errorsError(e1, e2) + "(Should be Equal)"
}

func TestInterpreter(t *testing.T) {
	Convey("interpreter works with", t, func() {
		Convey("nothing", func() {
			result, errors := calcgo.Interpret("")
			So(result, ShouldEqual, 0)
			So(errors, ShouldBeNil)
		})
		Convey("positive integers", func() {
			result, errors := calcgo.Interpret("1")
			So(result, ShouldEqual, 1)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("12345")
			So(result, ShouldEqual, 12345)
			So(errors, ShouldBeNil)
		})

		Convey("positive decimals", func() {
			result, errors := calcgo.Interpret("1.0")
			So(result, ShouldEqual, 1.0)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("1234.5678")
			So(result, ShouldEqual, 1234.5678)
			So(errors, ShouldBeNil)
		})

		Convey("negativ numbers", func() {
			result, errors := calcgo.Interpret("-1")
			So(result, ShouldEqual, -1)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("-1.1")
			So(result, ShouldEqual, -1.1)
			So(errors, ShouldBeNil)
		})

		Convey("simple additions with integers", func() {
			result, errors := calcgo.Interpret("1 + 1")
			So(result, ShouldEqual, 2)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("3 + 5")
			So(result, ShouldEqual, 3+5)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("1 + 2 + 3 + 4 + 5 + 6")
			So(result, ShouldEqual, 1+2+3+4+5+6)
			So(errors, ShouldBeNil)
		})

		Convey("simple additions with decimals", func() {
			result, errors := calcgo.Interpret("1.2 + 2.4")
			SkipSo(result, ShouldEqual, 1.2+2.4) // @todo: fix rounding error
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("0.7 + 2.4")
			SkipSo(result, ShouldEqual, 0.7+2.4) // @todo: fix rounding error
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("3.5 + 5.1")
			So(result, ShouldEqual, 3.5+5.1)
			So(errors, ShouldBeNil)
		})

		Convey("simple subtractions", func() {
			result, errors := calcgo.Interpret("1 - 1")
			So(result, ShouldEqual, 1-1)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("3 - 5")
			So(result, ShouldEqual, 3-5)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("1 - 2 - 3 - 4 - 5 - 6")
			So(result, ShouldEqual, 1-2-3-4-5-6)
			So(errors, ShouldBeNil)
		})

		Convey("simple multiplications", func() {
			result, errors := calcgo.Interpret("1 * 1")
			So(result, ShouldEqual, 1*1)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("3 * 5")
			So(result, ShouldEqual, 3*5)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("1 * 2 * 3 * 4 * 5 * 6")
			So(result, ShouldEqual, 1*2*3*4*5*6)
			So(errors, ShouldBeNil)
		})

		Convey("simple divisions", func() {
			result, errors := calcgo.Interpret("1 / 1")
			So(result, ShouldEqual, 1/1)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("3 / 5")
			So(result, ShouldEqual, 3.0/5.0)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("1 / 2 / 3 / 4 / 5 / 6")
			So(result, ShouldEqual, 1.0/2.0/3.0/4.0/5.0/6.0)
			So(errors, ShouldBeNil)
		})

		Convey("basic operations with negative numbers", func() {
			result, errors := calcgo.Interpret("-1 + 2")
			So(result, ShouldEqual, -1+2)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("-1 - 2")
			So(result, ShouldEqual, -1-2)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("-1 * 2")
			So(result, ShouldEqual, -1*2)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("-1 / 2")
			So(result, ShouldEqual, -1.0/2.0)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("-1 + -2")
			So(result, ShouldEqual, -1+-2)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("-1 - -2")
			So(result, ShouldEqual, -1 - -2)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("-1 * -2")
			So(result, ShouldEqual, -1*-2)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("-1 / -2")
			So(result, ShouldEqual, -1.0/-2.0)
			So(errors, ShouldBeNil)
		})

		Convey("dot before line rule", func() {
			result, errors := calcgo.Interpret("1 + 2 / 3")
			SkipSo(result, ShouldEqual, 1.0+2.0/3.0) // @todo: fix rounding error
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("1 - 2 / 3")
			SkipSo(result, ShouldEqual, 1.0-2.0/3.0) // @todo: fix rounding error
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("1 + 2 * 3")
			So(result, ShouldEqual, 1.0+2.0*3.0)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("1 - 2 * 3")
			So(result, ShouldEqual, 1.0-2.0*3.0)
			So(errors, ShouldBeNil)
		})

		Convey("brackets", func() {
			result, errors := calcgo.Interpret("(1 + 2) / 3")
			So(result, ShouldEqual, (1.0+2.0)/3.0)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("(1 - 2) / 3")
			So(result, ShouldEqual, (1.0-2.0)/3.0)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("(1 + 2) * 3")
			So(result, ShouldEqual, (1.0+2.0)*3.0)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("(1 - 2) * 3")
			So(result, ShouldEqual, (1.0-2.0)*3.0)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("2 + (1 - 2) / 3")
			So(result, ShouldEqual, 2.0+(1.0-2.0)/3.0)
			So(errors, ShouldBeNil)
		})

		Convey("nested brackets", func() {
			result, errors := calcgo.Interpret("((1 + 2) / 3) + 1")
			So(result, ShouldEqual, ((1.0+2.0)/3.0)+1)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("((2 + 3) / (1 + 2)) * 3")
			So(result, ShouldEqual, ((2.0+3.0)/(1.0+2.0))*3.0)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("(1 - 2) * (3 - 2) / (1 + 4)")
			So(result, ShouldEqual, (1.0-2.0)*(3.0-2.0)/(1.0+4.0))
			So(errors, ShouldBeNil)
		})

		Convey("brackets and dot before line rule", func() {
			result, errors := calcgo.Interpret("1 + (1 + 2) * 3")
			So(result, ShouldEqual, 1.0+(1.0+2.0)*3.0)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("1 + (1 + 2) / 3")
			So(result, ShouldEqual, 1.0+(1.0+2.0)/3.0)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("1 - (1 + 2) * 3")
			So(result, ShouldEqual, 1.0-(1.0+2.0)*3.0)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("1 - (1 + 2) / 3")
			So(result, ShouldEqual, 1.0-(1.0+2.0)/3.0)
			So(errors, ShouldBeNil)

			result, errors = calcgo.Interpret("(1 + 2) * 3 + (4 - 6 / (5 + 2))")
			So(result, ShouldEqual, (1.0+2.0)*3.0+(4.0-6.0/(5.0+2.0)))
			So(errors, ShouldBeNil)
		})
	})

	Convey("variables", t, func() {
		Convey("works with simple variables", func() {
			i := calcgo.NewInterpreter("a")
			i.SetVar("a", 1.0)
			result, errors := i.GetResult()
			So(result, ShouldEqual, 1)
			So(errors, ShouldBeNil)

			i = calcgo.NewInterpreter("1 + a")
			i.SetVar("a", 1.0)
			result, errors = i.GetResult()
			So(result, ShouldEqual, 2)
			So(errors, ShouldBeNil)
		})

		Convey("works with multiple variables", func() {
			i := calcgo.NewInterpreter("a + b")
			i.SetVar("a", 1)
			i.SetVar("b", 2)
			result, errors := i.GetResult()
			So(result, ShouldEqual, 3)
			So(errors, ShouldBeNil)
		})

		Convey("works with reassining variables", func() {
			i := calcgo.NewInterpreter("1 + a")

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
			i := calcgo.NewInterpreter("1 + a")
			result, errors := i.GetResult()
			So(result, ShouldEqual, 0)
			So(errors, ShouldEqualErrors, []error{
				calcgo.ErrorVariableNotDefined,
			})
		})
	})

	Convey("interpreter returns errors, when parser returned errors", t, func() {
		result, errors := calcgo.Interpret("$")
		So(result, ShouldEqual, 0)
		So(errors, ShouldEqualErrors, []error{
			calcgo.ErrorExpectedNumberOrVariable,
		})

		result, errors = calcgo.Interpret("1 + #)")
		So(result, ShouldEqual, 0)
		So(errors, ShouldEqualErrors, []error{
			calcgo.ErrorExpectedNumberOrVariable,
			calcgo.ErrorUnexpectedClosingBracket,
		})
	})

	Convey("interpreter returns error when", t, func() {
		Convey("dividing by 0", func() {
			result, errors := calcgo.Interpret("1 / 0")
			So(result, ShouldEqual, 0)
			So(errors, ShouldEqualErrors, []error{
				calcgo.ErrorDivisionByZero,
			})
		})
		Convey("a node child is missing", func() {
			result, errors := calcgo.InterpretAST(calcgo.AST{
				Node: &calcgo.Node{
					Type:      calcgo.NAddition,
					Value:     "",
					LeftChild: nil,
					RightChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "2",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(result, ShouldEqual, 0)
			So(errors, ShouldEqual, calcgo.ErrorMissingLeftChild)

			result, errors = calcgo.InterpretAST(calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NAddition,
					Value: "",
					LeftChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "2",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: nil,
				},
			})
			So(result, ShouldEqual, 0)
			So(errors, ShouldEqual, calcgo.ErrorMissingRightChild)
		})

		Convey("a wrong number is given", func() {
			result, errors := calcgo.InterpretAST(calcgo.AST{
				Node: &calcgo.Node{
					Type:       calcgo.NInteger,
					Value:      "a",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(result, ShouldEqual, 0)
			So(errors, ShouldEqual, calcgo.ErrorInvalidInteger)

			result, errors = calcgo.InterpretAST(calcgo.AST{
				Node: &calcgo.Node{
					Type:       calcgo.NDecimal,
					Value:      "a",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(result, ShouldEqual, 0)
			So(errors, ShouldEqual, calcgo.ErrorInvalidDecimal)
		})

		Convey("an invalid node type if given", func() {
			result, errors := calcgo.InterpretAST(calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NAddition,
					Value: "",
					LeftChild: &calcgo.Node{
						Type:       30000,
						Value:      "a",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(result, ShouldEqual, 0)
			So(errors, ShouldEqual, calcgo.ErrorInvalidNodeType)
		})

		Convey("the error doesn't happen on the first node", func() {
			result, errors := calcgo.InterpretAST(calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NAddition,
					Value: "",
					LeftChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "a",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(result, ShouldEqual, 0)
			So(errors, ShouldEqual, calcgo.ErrorInvalidInteger)

			result, errors = calcgo.InterpretAST(calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NAddition,
					Value: "",
					LeftChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "a",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(result, ShouldEqual, 0)
			So(errors, ShouldEqual, calcgo.ErrorInvalidInteger)
		})
	})
}
