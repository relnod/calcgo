package calcgo_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/relnod/calcgo"
	"github.com/relnod/calcgo/lexer"
)

func astToString(ast calcgo.AST) string {
	str, err := json.MarshalIndent(ast, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	return string(str) + "\n"
}

func parseError(actual, expected calcgo.AST) string {
	return "Expected: \n" +
		astToString(expected) +
		"Actual: \n" +
		astToString(actual)
}

func eqNodes(n1, n2 *calcgo.Node) bool {
	if n1 == nil && n2 == nil {
		return true
	}

	if n1 == nil && n2 != nil {
		return false
	}

	if n1 != nil && n2 == nil {
		return false
	}

	if n1.Type != n2.Type {
		return false
	}

	if n1.Value != n2.Value {
		return false
	}

	if !eqNodes(n1.LeftChild, n2.LeftChild) {
		return false
	}

	if !eqNodes(n1.RightChild, n2.RightChild) {
		return false
	}

	return true
}

func shouldEqualAST(actual interface{}, expected ...interface{}) string {
	actualAST := actual.(calcgo.AST)
	expectedAST := expected[0].(calcgo.AST)

	if eqNodes(actualAST.Node, expectedAST.Node) {
		return ""
	}

	return parseError(actualAST, expectedAST) + "(Should be Equal)"
}

func TestParser(t *testing.T) {
	Convey("IsOperator works", t, func() {
		So(calcgo.IsOperator(calcgo.NDecimal), ShouldEqual, false)
		So(calcgo.IsOperator(calcgo.NInteger), ShouldEqual, false)
		So(calcgo.IsOperator(calcgo.NError), ShouldEqual, false)
		So(calcgo.IsOperator(calcgo.NAddition), ShouldEqual, true)
		So(calcgo.IsOperator(calcgo.NSubtraction), ShouldEqual, true)
		So(calcgo.IsOperator(calcgo.NMultiplication), ShouldEqual, true)
		So(calcgo.IsOperator(calcgo.NDivision), ShouldEqual, true)
	})

	Convey("Parser works with", t, func() {
		Convey("nothing", func() {
			ast, errors := calcgo.Parse("")
			So(ast, shouldEqualAST, calcgo.AST{})
			So(errors, ShouldBeNil)
		})

		Convey("simple numbers", func() {
			ast, errors := calcgo.Parse("20")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:       calcgo.NInteger,
					Value:      "20",
					LeftChild:  nil,
					RightChild: nil,
				},
			})

			So(errors, ShouldBeNil)

			ast, errors = calcgo.Parse("20.23")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:       calcgo.NDecimal,
					Value:      "20.23",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(errors, ShouldBeNil)

		})

		Convey("negativ numbers", func() {
			ast, errors := calcgo.Parse("-1")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:       calcgo.NInteger,
					Value:      "-1",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = calcgo.Parse("-123.45")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:       calcgo.NDecimal,
					Value:      "-123.45",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(errors, ShouldBeNil)
		})

		Convey("variables", func() {
			ast, errors := calcgo.Parse("a")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:       calcgo.NVariable,
					Value:      "a",
					LeftChild:  nil,
					RightChild: nil,
				},
			})

			So(errors, ShouldBeNil)

			ast, errors = calcgo.Parse("abcdef")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:       calcgo.NVariable,
					Value:      "abcdef",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(errors, ShouldBeNil)

		})

		Convey("additions", func() {
			ast, errors := calcgo.Parse("1 + 2")
			So(ast, shouldEqualAST, calcgo.AST{
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
						Value:      "2",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = calcgo.Parse("1 + 2 + 3")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NAddition,
					Value: "",
					LeftChild: &calcgo.Node{
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
							Value:      "2",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
					RightChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "3",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)
		})

		Convey("addition with decimals", func() {
			ast, errors := calcgo.Parse("1.2 + 2.4")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NAddition,
					Value: "",
					LeftChild: &calcgo.Node{
						Type:       calcgo.NDecimal,
						Value:      "1.2",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &calcgo.Node{
						Type:       calcgo.NDecimal,
						Value:      "2.4",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)
		})

		Convey("subtractions", func() {
			ast, errors := calcgo.Parse("1 - 2")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NSubtraction,
					Value: "",
					LeftChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "2",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = calcgo.Parse("1 - 2 - 3")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NSubtraction,
					Value: "",
					LeftChild: &calcgo.Node{
						Type:  calcgo.NSubtraction,
						Value: "",
						LeftChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "1",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "2",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
					RightChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "3",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)
		})

		Convey("Multiplications", func() {
			ast, errors := calcgo.Parse("1 * 2")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NMultiplication,
					Value: "",
					LeftChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "2",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = calcgo.Parse("1 * 2 * 3")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NMultiplication,
					Value: "",
					LeftChild: &calcgo.Node{
						Type:  calcgo.NMultiplication,
						Value: "",
						LeftChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "1",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "2",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
					RightChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "3",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)
		})

		Convey("Divisions", func() {
			ast, errors := calcgo.Parse("1 / 2")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NDivision,
					Value: "",
					LeftChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "2",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = calcgo.Parse("1 / 2 / 3")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NDivision,
					Value: "",
					LeftChild: &calcgo.Node{
						Type:  calcgo.NDivision,
						Value: "",
						LeftChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "1",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "2",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
					RightChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "3",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)
		})

		Convey("variables and operators", func() {
			ast, errors := calcgo.Parse("a + 2")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NAddition,
					Value: "",
					LeftChild: &calcgo.Node{
						Type:       calcgo.NVariable,
						Value:      "a",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "2",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = calcgo.Parse("a - b")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NSubtraction,
					Value: "",
					LeftChild: &calcgo.Node{
						Type:       calcgo.NVariable,
						Value:      "a",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &calcgo.Node{
						Type:       calcgo.NVariable,
						Value:      "b",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)
		})

		Convey("'multiplication and division before addition and subtraction' rule", func() {
			ast, errors := calcgo.Parse("1 + 2 * 3")
			So(ast, shouldEqualAST, calcgo.AST{
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
						Type:  calcgo.NMultiplication,
						Value: "",
						LeftChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "2",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "3",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = calcgo.Parse("1 - 2 / 3")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NSubtraction,
					Value: "",
					LeftChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &calcgo.Node{
						Type:  calcgo.NDivision,
						Value: "",
						LeftChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "2",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "3",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = calcgo.Parse("2 * 3 + 1")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NAddition,
					Value: "",
					LeftChild: &calcgo.Node{
						Type:  calcgo.NMultiplication,
						Value: "",
						LeftChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "2",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "3",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
					RightChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)
		})

		Convey("brackets", func() {
			ast, errors := calcgo.Parse("(1)")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:       calcgo.NInteger,
					Value:      "1",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = calcgo.Parse("(1 - 2)")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NSubtraction,
					Value: "",
					LeftChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "2",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = calcgo.Parse("(1 - 2) * 3")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NMultiplication,
					Value: "",
					LeftChild: &calcgo.Node{
						Type:  calcgo.NSubtraction,
						Value: "",
						LeftChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "1",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "2",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
					RightChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "3",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = calcgo.Parse("3 * (1 - 2)")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NMultiplication,
					Value: "",
					LeftChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "3",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &calcgo.Node{
						Type:  calcgo.NSubtraction,
						Value: "",
						LeftChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "1",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "2",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = calcgo.Parse("3 * (1 - 2) / 4")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NDivision,
					Value: "",
					LeftChild: &calcgo.Node{
						Type:  calcgo.NMultiplication,
						Value: "",
						LeftChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "3",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &calcgo.Node{
							Type:  calcgo.NSubtraction,
							Value: "",
							LeftChild: &calcgo.Node{
								Type:       calcgo.NInteger,
								Value:      "1",
								LeftChild:  nil,
								RightChild: nil,
							},
							RightChild: &calcgo.Node{
								Type:       calcgo.NInteger,
								Value:      "2",
								LeftChild:  nil,
								RightChild: nil,
							},
						},
					},
					RightChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "4",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)
		})

		Convey("nested brackets", func() {
			ast, errors := calcgo.Parse("((1))")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:       calcgo.NInteger,
					Value:      "1",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = calcgo.Parse("((1 - 2))")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NSubtraction,
					Value: "",
					LeftChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "2",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = calcgo.Parse("(3 * (1 - 2))")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NMultiplication,
					Value: "",
					LeftChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "3",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &calcgo.Node{
						Type:  calcgo.NSubtraction,
						Value: "",
						LeftChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "1",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "2",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
				},
			})
			So(errors, ShouldBeNil)
		})

		Convey("Brackets with 'multiplication and division before addition and subtraction' rule", func() {
			ast, errors := calcgo.Parse("3 + (1 - 2) / 4")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NAddition,
					Value: "",
					LeftChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "3",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &calcgo.Node{
						Type:  calcgo.NDivision,
						Value: "",
						LeftChild: &calcgo.Node{
							Type:  calcgo.NSubtraction,
							Value: "",
							LeftChild: &calcgo.Node{
								Type:       calcgo.NInteger,
								Value:      "1",
								LeftChild:  nil,
								RightChild: nil,
							},
							RightChild: &calcgo.Node{
								Type:       calcgo.NInteger,
								Value:      "2",
								LeftChild:  nil,
								RightChild: nil,
							},
						},
						RightChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "4",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = calcgo.Parse("3 + (1 + 2) * 4")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NAddition,
					Value: "",
					LeftChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "3",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &calcgo.Node{
						Type:  calcgo.NMultiplication,
						Value: "",
						LeftChild: &calcgo.Node{
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
								Value:      "2",
								LeftChild:  nil,
								RightChild: nil,
							},
						},
						RightChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "4",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = calcgo.Parse("(1 + 2) * 4 + 1")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NAddition,
					Value: "",
					LeftChild: &calcgo.Node{
						Type:  calcgo.NMultiplication,
						Value: "",
						LeftChild: &calcgo.Node{
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
								Value:      "2",
								LeftChild:  nil,
								RightChild: nil,
							},
						},
						RightChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "4",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
					RightChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = calcgo.Parse("4 - 6 / (5 + 2)")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NSubtraction,
					Value: "",
					LeftChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "4",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &calcgo.Node{
						Type:  calcgo.NDivision,
						Value: "",
						LeftChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "6",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &calcgo.Node{
							Type:  calcgo.NAddition,
							Value: "",
							LeftChild: &calcgo.Node{
								Type:       calcgo.NInteger,
								Value:      "5",
								LeftChild:  nil,
								RightChild: nil,
							},
							RightChild: &calcgo.Node{
								Type:       calcgo.NInteger,
								Value:      "2",
								LeftChild:  nil,
								RightChild: nil,
							},
						},
					},
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = calcgo.Parse("(1 + 2) * 3 + (4 - 6 / (5 + 2))")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NAddition,
					Value: "",
					LeftChild: &calcgo.Node{
						Type:  calcgo.NMultiplication,
						Value: "",
						LeftChild: &calcgo.Node{
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
								Value:      "2",
								LeftChild:  nil,
								RightChild: nil,
							},
						},
						RightChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "3",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
					RightChild: &calcgo.Node{
						Type:  calcgo.NSubtraction,
						Value: "",
						LeftChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "4",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &calcgo.Node{
							Type:  calcgo.NDivision,
							Value: "",
							LeftChild: &calcgo.Node{
								Type:       calcgo.NInteger,
								Value:      "6",
								LeftChild:  nil,
								RightChild: nil,
							},
							RightChild: &calcgo.Node{
								Type:  calcgo.NAddition,
								Value: "",
								LeftChild: &calcgo.Node{
									Type:       calcgo.NInteger,
									Value:      "5",
									LeftChild:  nil,
									RightChild: nil,
								},
								RightChild: &calcgo.Node{
									Type:       calcgo.NInteger,
									Value:      "2",
									LeftChild:  nil,
									RightChild: nil,
								},
							},
						},
					},
				},
			})
			So(errors, ShouldBeNil)
		})

		Convey("functions", func() {
			Convey("sqrt", func() {
				ast, errors := calcgo.Parse("sqrt(1)")
				So(ast, shouldEqualAST, calcgo.AST{
					Node: &calcgo.Node{
						Type:  calcgo.NFuncSqrt,
						Value: "",
						LeftChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "1",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: nil,
					},
				})
				So(errors, ShouldBeNil)

				ast, errors = calcgo.Parse("sqrt(1 + 1)")
				So(ast, shouldEqualAST, calcgo.AST{
					Node: &calcgo.Node{
						Type:  calcgo.NFuncSqrt,
						Value: "",
						LeftChild: &calcgo.Node{
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
								Value:      "1",
								LeftChild:  nil,
								RightChild: nil,
							},
						},
						RightChild: nil,
					},
				})
				So(errors, ShouldBeNil)

				ast, errors = calcgo.Parse("1 + sqrt(1)")
				So(ast, shouldEqualAST, calcgo.AST{
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
							Type:  calcgo.NFuncSqrt,
							Value: "",
							LeftChild: &calcgo.Node{
								Type:       calcgo.NInteger,
								Value:      "1",
								LeftChild:  nil,
								RightChild: nil,
							},
							RightChild: nil,
						},
					},
				})
				So(errors, ShouldBeNil)
			})
		})
	})

	Convey("Parser works with errors", t, func() {
		Convey("handles invalid number", func() {
			ast, errors := calcgo.Parse("1#")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:       calcgo.NInvalidNumber,
					Value:      "#",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(errors, ShouldEqualErrors, []error{
				calcgo.ErrorExpectedNumberOrVariable,
			})

			ast, errors = calcgo.Parse("1 + 3#")
			So(ast, shouldEqualAST, calcgo.AST{
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
						Type:       calcgo.NInvalidNumber,
						Value:      "#",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldEqualErrors, []error{
				calcgo.ErrorExpectedNumberOrVariable,
			})

			Convey("handles multiple invalid number errors", func() {
				ast, errors := calcgo.Parse("2# + 3'")
				So(ast, shouldEqualAST, calcgo.AST{
					Node: &calcgo.Node{
						Type:  calcgo.NAddition,
						Value: "",
						LeftChild: &calcgo.Node{
							Type:       calcgo.NInvalidNumber,
							Value:      "#",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &calcgo.Node{
							Type:       calcgo.NInvalidNumber,
							Value:      "'",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
				})
				So(errors, ShouldEqualErrors, []error{
					calcgo.ErrorExpectedNumberOrVariable,
					calcgo.ErrorExpectedNumberOrVariable,
				})
			})
		})

		Convey("handles invalid variable", func() {
			ast, errors := calcgo.Parse("a#")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:       calcgo.NInvalidVariable,
					Value:      "#",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(errors, ShouldEqualErrors, []error{
				calcgo.ErrorExpectedNumberOrVariable,
			})
		})

		Convey("handles invalid number or variable", func() {
			ast, errors := calcgo.Parse("#")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:       calcgo.NError,
					Value:      "#",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(errors, ShouldEqualErrors, []error{
				calcgo.ErrorExpectedNumberOrVariable,
			})
		})

		Convey("handles invalid operator", func() {
			ast, errors := calcgo.Parse("1 $ 1")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NInvalidOperator,
					Value: "$",
					LeftChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "1",
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
			So(errors, ShouldEqualErrors, []error{
				calcgo.ErrorExpectedOperator,
			})

			Convey("handles multiple invalid operator errors", func() {
				ast, errors := calcgo.Parse("1 $ 1 $ 1")
				So(ast, shouldEqualAST, calcgo.AST{
					Node: &calcgo.Node{
						Type:  calcgo.NInvalidOperator,
						Value: "$",
						LeftChild: &calcgo.Node{
							Type:  calcgo.NInvalidOperator,
							Value: "$",
							LeftChild: &calcgo.Node{
								Type:       calcgo.NInteger,
								Value:      "1",
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
						RightChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "1",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
				})
				So(errors, ShouldEqualErrors, []error{
					calcgo.ErrorExpectedOperator,
					calcgo.ErrorExpectedOperator,
				})
			})
		})

		Convey("handles multiple mixed errors", func() {
			ast, errors := calcgo.Parse("1# $ 1#")
			So(ast, shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NInvalidOperator,
					Value: "$",
					LeftChild: &calcgo.Node{
						Type:       calcgo.NInvalidNumber,
						Value:      "#",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &calcgo.Node{
						Type:       calcgo.NInvalidNumber,
						Value:      "#",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldEqualErrors, []error{
				calcgo.ErrorExpectedNumberOrVariable,
				calcgo.ErrorExpectedOperator,
				calcgo.ErrorExpectedNumberOrVariable,
			})
		})

		Convey("handles missing closing bracket", func() {
			ast, errors := calcgo.Parse("(1 + 1")
			So(ast, shouldEqualAST, calcgo.AST{
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
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldEqualErrors, []error{
				calcgo.ErrorMissingClosingBracket,
			})
		})

		Convey("handles unexpected closing bracket", func() {
			ast, errors := calcgo.Parse("1 + 1)")
			So(ast, shouldEqualAST, calcgo.AST{
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
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldEqualErrors, []error{
				calcgo.ErrorUnexpectedClosingBracket,
			})

			ast, errors = calcgo.Parse("(1 + 1))")
			So(ast, shouldEqualAST, calcgo.AST{
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
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldEqualErrors, []error{
				calcgo.ErrorUnexpectedClosingBracket,
			})
		})
	})

	Convey("ParseTokens works the same as parse", t, func() {
		ast1, e1 := calcgo.ParseTokens(nil)
		ast2, e2 := calcgo.Parse("")
		So(ast1, shouldEqualAST, ast2)
		So(e1, ShouldEqualErrors, e2)

		ast1, e1 = calcgo.ParseTokens([]lexer.Token{
			{Type: lexer.TInteger, Value: "1"},
			{Type: lexer.TOperatorPlus, Value: ""},
			{Type: lexer.TInteger, Value: "1"},
		})
		ast2, e2 = calcgo.Parse("1 + 1")
		So(ast1, shouldEqualAST, ast2)
		So(e1, ShouldEqualErrors, e2)
	})
}
