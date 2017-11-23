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

func shouldEqualErrors(actual interface{}, expected ...interface{}) string {
	e1 := actual.([]error)
	e2 := expected[0].([]error)

	if eqErrors(e1, e2) {
		return ""
	}

	return errorsError(e1, e2) + "(Should be Equal)"
}

func getInterpretNumber(str string) float64 {
	number, _ := calcgo.Interpret(str)
	return number
}

func getInterpretError(str string) []error {
	_, errors := calcgo.Interpret(str)
	return errors
}

func getInterpretASTError(ast calcgo.AST) error {
	_, err := calcgo.InterpretAST(ast)
	return err
}

func TestInterpreter(t *testing.T) {
	Convey("interpreter works with", t, func() {
		Convey("nothing", func() {
			So(getInterpretNumber(""), ShouldEqual, 0)
		})
		Convey("positive integers", func() {
			So(getInterpretNumber("1"), ShouldEqual, 1)
			So(getInterpretNumber("12345"), ShouldEqual, 12345)
		})

		Convey("positive decimals", func() {
			So(getInterpretNumber("1.0"), ShouldEqual, 1.0)
			So(getInterpretNumber("1234.5678"), ShouldEqual, 1234.5678)
		})

		Convey("negativ numbers", func() {
			So(getInterpretNumber("-1"), ShouldEqual, -1)
			So(getInterpretNumber("-1.1"), ShouldEqual, -1.1)
		})

		Convey("simple additions with integers", func() {
			So(getInterpretNumber("1 + 1"), ShouldEqual, 2)
			So(getInterpretNumber("3 + 5"), ShouldEqual, 3+5)
			So(getInterpretNumber("1 + 2 + 3 + 4 + 5 + 6"), ShouldEqual, 1+2+3+4+5+6)
		})

		Convey("simple additions with decimals", func() {
			SkipSo(getInterpretNumber("1.2 + 2.4"), ShouldEqual, 1.2+2.4) // @todo: fix rounding error
			SkipSo(getInterpretNumber("0.7 + 2.4"), ShouldEqual, 0.7+2.4) // @todo: fix rounding error
			So(getInterpretNumber("3.5 + 5.1"), ShouldEqual, 3.5+5.1)
		})

		Convey("simple subtractions", func() {
			So(getInterpretNumber("1 - 1"), ShouldEqual, 1-1)
			So(getInterpretNumber("3 - 5"), ShouldEqual, 3-5)
			So(getInterpretNumber("1 - 2 - 3 - 4 - 5 - 6"), ShouldEqual, 1-2-3-4-5-6)
		})

		Convey("simple multiplications", func() {
			So(getInterpretNumber("1 * 1"), ShouldEqual, 1*1)
			So(getInterpretNumber("3 * 5"), ShouldEqual, 3*5)
			So(getInterpretNumber("1 * 2 * 3 * 4 * 5 * 6"), ShouldEqual, 1*2*3*4*5*6)
		})

		Convey("simple divisions", func() {
			So(getInterpretNumber("1 / 1"), ShouldEqual, 1/1)
			So(getInterpretNumber("3 / 5"), ShouldEqual, 3.0/5.0)
			So(getInterpretNumber("1 / 2 / 3 / 4 / 5 / 6"), ShouldEqual, 1.0/2.0/3.0/4.0/5.0/6.0)
		})

		Convey("basic operations with negative numbers", func() {
			So(getInterpretNumber("-1 + 2"), ShouldEqual, -1+2)
			So(getInterpretNumber("-1 - 2"), ShouldEqual, -1-2)
			So(getInterpretNumber("-1 * 2"), ShouldEqual, -1*2)
			So(getInterpretNumber("-1 / 2"), ShouldEqual, -1.0/2.0)

			So(getInterpretNumber("-1 + -2"), ShouldEqual, -1+-2)
			So(getInterpretNumber("-1 - -2"), ShouldEqual, -1 - -2)
			So(getInterpretNumber("-1 * -2"), ShouldEqual, -1*-2)
			So(getInterpretNumber("-1 / -2"), ShouldEqual, -1.0/-2.0)
		})

		Convey("dot before line rule", func() {
			SkipSo(getInterpretNumber("1 + 2 / 3"), ShouldEqual, 1.0+2.0/3.0) // @todo: fix rounding error
			SkipSo(getInterpretNumber("1 - 2 / 3"), ShouldEqual, 1.0-2.0/3.0) // @todo: fix rounding error
			So(getInterpretNumber("1 + 2 * 3"), ShouldEqual, 1.0+2.0*3.0)
			So(getInterpretNumber("1 - 2 * 3"), ShouldEqual, 1.0-2.0*3.0)
		})

		Convey("brackets", func() {
			So(getInterpretNumber("(1 + 2) / 3"), ShouldEqual, (1.0+2.0)/3.0)
			So(getInterpretNumber("(1 - 2) / 3"), ShouldEqual, (1.0-2.0)/3.0)
			So(getInterpretNumber("(1 + 2) * 3"), ShouldEqual, (1.0+2.0)*3.0)
			So(getInterpretNumber("(1 - 2) * 3"), ShouldEqual, (1.0-2.0)*3.0)
			So(getInterpretNumber("2 + (1 - 2) / 3"), ShouldEqual, 2.0+(1.0-2.0)/3.0)
		})

		Convey("nested brackets", func() {
			So(getInterpretNumber("((1 + 2) / 3) + 1"), ShouldEqual, ((1.0+2.0)/3.0)+1)
			So(getInterpretNumber("((2 + 3) / (1 + 2)) * 3"), ShouldEqual, ((2.0+3.0)/(1.0+2.0))*3.0)
			So(getInterpretNumber("(1 - 2) * (3 - 2) / (1 + 4)"), ShouldEqual, (1.0-2.0)*(3.0-2.0)/(1.0+4.0))
		})

		Convey("brackets and dot before line rule", func() {
			So(getInterpretNumber("1 + (1 + 2) * 3"), ShouldEqual, 1.0+(1.0+2.0)*3.0)
			So(getInterpretNumber("1 + (1 + 2) / 3"), ShouldEqual, 1.0+(1.0+2.0)/3.0)
			So(getInterpretNumber("1 - (1 + 2) * 3"), ShouldEqual, 1.0-(1.0+2.0)*3.0)
			So(getInterpretNumber("1 - (1 + 2) / 3"), ShouldEqual, 1.0-(1.0+2.0)/3.0)
			So(getInterpretNumber("(1 + 2) * 3 + (4 - 6 / (5 + 2))"), ShouldEqual, (1.0+2.0)*3.0+(4.0-6.0/(5.0+2.0)))
		})
	})

	Convey("interpreter returns error when", t, func() {
		Convey("dividing by 0", func() {
			result, errors := calcgo.Interpret("1 / 0")
			So(errors, shouldEqualErrors, []error{calcgo.ErrorDivisionByZero})
			So(result, ShouldEqual, 0)
		})
		Convey("a node child is missing", func() {
			So(getInterpretASTError(calcgo.AST{
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
			}), ShouldEqual, calcgo.ErrorMissingLeftChild)

			So(getInterpretASTError(calcgo.AST{
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
			}), ShouldEqual, calcgo.ErrorMissingRightChild)
		})

		Convey("a wrong number is given", func() {
			So(getInterpretASTError(calcgo.AST{
				Node: &calcgo.Node{
					Type:       calcgo.NInteger,
					Value:      "a",
					LeftChild:  nil,
					RightChild: nil,
				},
			}), ShouldEqual, calcgo.ErrorInvalidInteger)

			So(getInterpretASTError(calcgo.AST{
				Node: &calcgo.Node{
					Type:       calcgo.NDecimal,
					Value:      "a",
					LeftChild:  nil,
					RightChild: nil,
				},
			}), ShouldEqual, calcgo.ErrorInvalidDecimal)
		})

		Convey("an invalid node type if given", func() {
			So(getInterpretASTError(calcgo.AST{
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
			}), ShouldEqual, calcgo.ErrorInvalidNodeType)
		})

		Convey("the error doesn't happen on the first node", func() {
			So(getInterpretASTError(calcgo.AST{
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
			}), ShouldEqual, calcgo.ErrorInvalidInteger)

			So(getInterpretASTError(calcgo.AST{
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
			}), ShouldEqual, calcgo.ErrorInvalidInteger)
		})
	})
}
