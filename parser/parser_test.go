package parser_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/relnod/calcgo/calcgotest"
	"github.com/relnod/calcgo/lexer"
	"github.com/relnod/calcgo/parser"
	. "github.com/smartystreets/goconvey/convey"
)

func astToString(ast parser.AST) string {
	str, err := json.MarshalIndent(ast, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	return string(str) + "\n"
}

func parseError(actual, expected parser.AST) string {
	return "Expected: \n" +
		astToString(expected) +
		"Actual: \n" +
		astToString(actual)
}

func eqNodes(n1, n2 *parser.Node) bool {
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
	actualAST := actual.(parser.AST)
	expectedAST := expected[0].(parser.AST)

	if eqNodes(actualAST.Node, expectedAST.Node) {
		return ""
	}

	return parseError(actualAST, expectedAST) + "(Should be Equal)"
}

func TestParser(t *testing.T) {
	Convey("Parser works with", t, func() {
		Convey("nothing", func() {
			ast, errors := parser.Parse("")
			So(ast, shouldEqualAST, parser.AST{})
			So(errors, ShouldBeNil)
		})

		Convey("simple numbers", func() {
			ast, errors := parser.Parse("20")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:       parser.NInteger,
					Value:      "20",
					LeftChild:  nil,
					RightChild: nil,
				},
			})

			So(errors, ShouldBeNil)

			ast, errors = parser.Parse("20.23")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:       parser.NDecimal,
					Value:      "20.23",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(errors, ShouldBeNil)

		})

		Convey("negativ numbers", func() {
			ast, errors := parser.Parse("-1")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:       parser.NInteger,
					Value:      "-1",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = parser.Parse("-123.45")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:       parser.NDecimal,
					Value:      "-123.45",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(errors, ShouldBeNil)
		})

		Convey("variables", func() {
			ast, errors := parser.Parse("a")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:       parser.NVariable,
					Value:      "a",
					LeftChild:  nil,
					RightChild: nil,
				},
			})

			So(errors, ShouldBeNil)

			ast, errors = parser.Parse("abcdef")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:       parser.NVariable,
					Value:      "abcdef",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(errors, ShouldBeNil)

		})

		Convey("additions", func() {
			ast, errors := parser.Parse("1 + 2")
			So(ast, shouldEqualAST, parser.AST{
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
			So(errors, ShouldBeNil)

			ast, errors = parser.Parse("1 + 2 + 3")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:  parser.NAddition,
					Value: "",
					LeftChild: &parser.Node{
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
					RightChild: &parser.Node{
						Type:       parser.NInteger,
						Value:      "3",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)
		})

		Convey("addition with decimals", func() {
			ast, errors := parser.Parse("1.2 + 2.4")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:  parser.NAddition,
					Value: "",
					LeftChild: &parser.Node{
						Type:       parser.NDecimal,
						Value:      "1.2",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &parser.Node{
						Type:       parser.NDecimal,
						Value:      "2.4",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)
		})

		Convey("subtractions", func() {
			ast, errors := parser.Parse("1 - 2")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:  parser.NSubtraction,
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
			So(errors, ShouldBeNil)

			ast, errors = parser.Parse("1 - 2 - 3")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:  parser.NSubtraction,
					Value: "",
					LeftChild: &parser.Node{
						Type:  parser.NSubtraction,
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
					RightChild: &parser.Node{
						Type:       parser.NInteger,
						Value:      "3",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)
		})

		Convey("Multiplications", func() {
			ast, errors := parser.Parse("1 * 2")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:  parser.NMultiplication,
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
			So(errors, ShouldBeNil)

			ast, errors = parser.Parse("1 * 2 * 3")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:  parser.NMultiplication,
					Value: "",
					LeftChild: &parser.Node{
						Type:  parser.NMultiplication,
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
					RightChild: &parser.Node{
						Type:       parser.NInteger,
						Value:      "3",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)
		})

		Convey("Divisions", func() {
			ast, errors := parser.Parse("1 / 2")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:  parser.NDivision,
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
			So(errors, ShouldBeNil)

			ast, errors = parser.Parse("1 / 2 / 3")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:  parser.NDivision,
					Value: "",
					LeftChild: &parser.Node{
						Type:  parser.NDivision,
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
					RightChild: &parser.Node{
						Type:       parser.NInteger,
						Value:      "3",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)
		})

		Convey("variables and operators", func() {
			ast, errors := parser.Parse("a + 2")
			So(ast, shouldEqualAST, parser.AST{
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
						Value:      "2",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = parser.Parse("a - b")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:  parser.NSubtraction,
					Value: "",
					LeftChild: &parser.Node{
						Type:       parser.NVariable,
						Value:      "a",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &parser.Node{
						Type:       parser.NVariable,
						Value:      "b",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)
		})

		Convey("'multiplication and division before addition and subtraction' rule", func() {
			ast, errors := parser.Parse("1 + 2 * 3")
			So(ast, shouldEqualAST, parser.AST{
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
						Type:  parser.NMultiplication,
						Value: "",
						LeftChild: &parser.Node{
							Type:       parser.NInteger,
							Value:      "2",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &parser.Node{
							Type:       parser.NInteger,
							Value:      "3",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = parser.Parse("1 - 2 / 3")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:  parser.NSubtraction,
					Value: "",
					LeftChild: &parser.Node{
						Type:       parser.NInteger,
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &parser.Node{
						Type:  parser.NDivision,
						Value: "",
						LeftChild: &parser.Node{
							Type:       parser.NInteger,
							Value:      "2",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &parser.Node{
							Type:       parser.NInteger,
							Value:      "3",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = parser.Parse("2 * 3 + 1")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:  parser.NAddition,
					Value: "",
					LeftChild: &parser.Node{
						Type:  parser.NMultiplication,
						Value: "",
						LeftChild: &parser.Node{
							Type:       parser.NInteger,
							Value:      "2",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &parser.Node{
							Type:       parser.NInteger,
							Value:      "3",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
					RightChild: &parser.Node{
						Type:       parser.NInteger,
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)
		})

		Convey("brackets", func() {
			ast, errors := parser.Parse("(1)")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:       parser.NInteger,
					Value:      "1",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = parser.Parse("(1 - 2)")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:  parser.NSubtraction,
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
			So(errors, ShouldBeNil)

			ast, errors = parser.Parse("(1 - 2) * 3")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:  parser.NMultiplication,
					Value: "",
					LeftChild: &parser.Node{
						Type:  parser.NSubtraction,
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
					RightChild: &parser.Node{
						Type:       parser.NInteger,
						Value:      "3",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = parser.Parse("3 * (1 - 2)")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:  parser.NMultiplication,
					Value: "",
					LeftChild: &parser.Node{
						Type:       parser.NInteger,
						Value:      "3",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &parser.Node{
						Type:  parser.NSubtraction,
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
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = parser.Parse("3 * (1 - 2) / 4")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:  parser.NDivision,
					Value: "",
					LeftChild: &parser.Node{
						Type:  parser.NMultiplication,
						Value: "",
						LeftChild: &parser.Node{
							Type:       parser.NInteger,
							Value:      "3",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &parser.Node{
							Type:  parser.NSubtraction,
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
					},
					RightChild: &parser.Node{
						Type:       parser.NInteger,
						Value:      "4",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)
		})

		Convey("nested brackets", func() {
			ast, errors := parser.Parse("((1))")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:       parser.NInteger,
					Value:      "1",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = parser.Parse("((1 - 2))")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:  parser.NSubtraction,
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
			So(errors, ShouldBeNil)

			ast, errors = parser.Parse("(3 * (1 - 2))")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:  parser.NMultiplication,
					Value: "",
					LeftChild: &parser.Node{
						Type:       parser.NInteger,
						Value:      "3",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &parser.Node{
						Type:  parser.NSubtraction,
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
				},
			})
			So(errors, ShouldBeNil)
		})

		Convey("Brackets with 'multiplication and division before addition and subtraction' rule", func() {
			ast, errors := parser.Parse("3 + (1 - 2) / 4")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:  parser.NAddition,
					Value: "",
					LeftChild: &parser.Node{
						Type:       parser.NInteger,
						Value:      "3",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &parser.Node{
						Type:  parser.NDivision,
						Value: "",
						LeftChild: &parser.Node{
							Type:  parser.NSubtraction,
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
						RightChild: &parser.Node{
							Type:       parser.NInteger,
							Value:      "4",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = parser.Parse("3 + (1 + 2) * 4")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:  parser.NAddition,
					Value: "",
					LeftChild: &parser.Node{
						Type:       parser.NInteger,
						Value:      "3",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &parser.Node{
						Type:  parser.NMultiplication,
						Value: "",
						LeftChild: &parser.Node{
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
						RightChild: &parser.Node{
							Type:       parser.NInteger,
							Value:      "4",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = parser.Parse("(1 + 2) * 4 + 1")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:  parser.NAddition,
					Value: "",
					LeftChild: &parser.Node{
						Type:  parser.NMultiplication,
						Value: "",
						LeftChild: &parser.Node{
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
						RightChild: &parser.Node{
							Type:       parser.NInteger,
							Value:      "4",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
					RightChild: &parser.Node{
						Type:       parser.NInteger,
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = parser.Parse("4 - 6 / (5 + 2)")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:  parser.NSubtraction,
					Value: "",
					LeftChild: &parser.Node{
						Type:       parser.NInteger,
						Value:      "4",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &parser.Node{
						Type:  parser.NDivision,
						Value: "",
						LeftChild: &parser.Node{
							Type:       parser.NInteger,
							Value:      "6",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &parser.Node{
							Type:  parser.NAddition,
							Value: "",
							LeftChild: &parser.Node{
								Type:       parser.NInteger,
								Value:      "5",
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
					},
				},
			})
			So(errors, ShouldBeNil)

			ast, errors = parser.Parse("(1 + 2) * 3 + (4 - 6 / (5 + 2))")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:  parser.NAddition,
					Value: "",
					LeftChild: &parser.Node{
						Type:  parser.NMultiplication,
						Value: "",
						LeftChild: &parser.Node{
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
						RightChild: &parser.Node{
							Type:       parser.NInteger,
							Value:      "3",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
					RightChild: &parser.Node{
						Type:  parser.NSubtraction,
						Value: "",
						LeftChild: &parser.Node{
							Type:       parser.NInteger,
							Value:      "4",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &parser.Node{
							Type:  parser.NDivision,
							Value: "",
							LeftChild: &parser.Node{
								Type:       parser.NInteger,
								Value:      "6",
								LeftChild:  nil,
								RightChild: nil,
							},
							RightChild: &parser.Node{
								Type:  parser.NAddition,
								Value: "",
								LeftChild: &parser.Node{
									Type:       parser.NInteger,
									Value:      "5",
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
						},
					},
				},
			})
			So(errors, ShouldBeNil)
		})

		Convey("functions", func() {
			Convey("sqrt", func() {
				ast, errors := parser.Parse("sqrt(1)")
				So(ast, shouldEqualAST, parser.AST{
					Node: &parser.Node{
						Type:  parser.NFuncSqrt,
						Value: "",
						LeftChild: &parser.Node{
							Type:       parser.NInteger,
							Value:      "1",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: nil,
					},
				})
				So(errors, ShouldBeNil)

				ast, errors = parser.Parse("sqrt(1 + 1)")
				So(ast, shouldEqualAST, parser.AST{
					Node: &parser.Node{
						Type:  parser.NFuncSqrt,
						Value: "",
						LeftChild: &parser.Node{
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
								Value:      "1",
								LeftChild:  nil,
								RightChild: nil,
							},
						},
						RightChild: nil,
					},
				})
				So(errors, ShouldBeNil)

				ast, errors = parser.Parse("1 + sqrt(1)")
				So(ast, shouldEqualAST, parser.AST{
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
							Type:  parser.NFuncSqrt,
							Value: "",
							LeftChild: &parser.Node{
								Type:       parser.NInteger,
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
			ast, errors := parser.Parse("1#")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:       parser.NInvalidNumber,
					Value:      "#",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(errors, ShouldEqualErrors, []error{
				parser.ErrorExpectedNumberOrVariable,
			})

			ast, errors = parser.Parse("1 + 3#")
			So(ast, shouldEqualAST, parser.AST{
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
						Type:       parser.NInvalidNumber,
						Value:      "#",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldEqualErrors, []error{
				parser.ErrorExpectedNumberOrVariable,
			})

			Convey("handles multiple invalid number errors", func() {
				ast, errors := parser.Parse("2# + 3'")
				So(ast, shouldEqualAST, parser.AST{
					Node: &parser.Node{
						Type:  parser.NAddition,
						Value: "",
						LeftChild: &parser.Node{
							Type:       parser.NInvalidNumber,
							Value:      "#",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &parser.Node{
							Type:       parser.NInvalidNumber,
							Value:      "'",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
				})
				So(errors, ShouldEqualErrors, []error{
					parser.ErrorExpectedNumberOrVariable,
					parser.ErrorExpectedNumberOrVariable,
				})
			})
		})

		Convey("handles invalid variable", func() {
			ast, errors := parser.Parse("a#")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:       parser.NInvalidVariable,
					Value:      "#",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(errors, ShouldEqualErrors, []error{
				parser.ErrorExpectedNumberOrVariable,
			})
		})

		Convey("handles invalid number or variable", func() {
			ast, errors := parser.Parse("#")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:       parser.NError,
					Value:      "#",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(errors, ShouldEqualErrors, []error{
				parser.ErrorExpectedNumberOrVariable,
			})
		})

		Convey("handles invalid operator", func() {
			ast, errors := parser.Parse("1 $ 1")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:  parser.NInvalidOperator,
					Value: "$",
					LeftChild: &parser.Node{
						Type:       parser.NInteger,
						Value:      "1",
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
			So(errors, ShouldEqualErrors, []error{
				parser.ErrorExpectedOperator,
			})

			Convey("handles multiple invalid operator errors", func() {
				ast, errors := parser.Parse("1 $ 1 $ 1")
				So(ast, shouldEqualAST, parser.AST{
					Node: &parser.Node{
						Type:  parser.NInvalidOperator,
						Value: "$",
						LeftChild: &parser.Node{
							Type:  parser.NInvalidOperator,
							Value: "$",
							LeftChild: &parser.Node{
								Type:       parser.NInteger,
								Value:      "1",
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
						RightChild: &parser.Node{
							Type:       parser.NInteger,
							Value:      "1",
							LeftChild:  nil,
							RightChild: nil,
						},
					},
				})
				So(errors, ShouldEqualErrors, []error{
					parser.ErrorExpectedOperator,
					parser.ErrorExpectedOperator,
				})
			})
		})

		Convey("handles multiple mixed errors", func() {
			ast, errors := parser.Parse("1# $ 1#")
			So(ast, shouldEqualAST, parser.AST{
				Node: &parser.Node{
					Type:  parser.NInvalidOperator,
					Value: "$",
					LeftChild: &parser.Node{
						Type:       parser.NInvalidNumber,
						Value:      "#",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &parser.Node{
						Type:       parser.NInvalidNumber,
						Value:      "#",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldEqualErrors, []error{
				parser.ErrorExpectedNumberOrVariable,
				parser.ErrorExpectedOperator,
				parser.ErrorExpectedNumberOrVariable,
			})
		})

		Convey("handles missing closing bracket", func() {
			ast, errors := parser.Parse("(1 + 1")
			So(ast, shouldEqualAST, parser.AST{
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
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldEqualErrors, []error{
				parser.ErrorMissingClosingBracket,
			})
		})

		Convey("handles unexpected closing bracket", func() {
			ast, errors := parser.Parse("1 + 1)")
			So(ast, shouldEqualAST, parser.AST{
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
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldEqualErrors, []error{
				parser.ErrorUnexpectedClosingBracket,
			})

			ast, errors = parser.Parse("(1 + 1))")
			So(ast, shouldEqualAST, parser.AST{
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
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(errors, ShouldEqualErrors, []error{
				parser.ErrorUnexpectedClosingBracket,
			})
		})
	})

	Convey("ParseTokens works the same as parse", t, func() {
		ast1, e1 := parser.ParseTokens(nil)
		ast2, e2 := parser.Parse("")
		So(ast1, shouldEqualAST, ast2)
		So(e1, ShouldEqualErrors, e2)

		ast1, e1 = parser.ParseTokens([]lexer.Token{
			{Type: lexer.TInteger, Value: "1"},
			{Type: lexer.TOperatorPlus, Value: ""},
			{Type: lexer.TInteger, Value: "1"},
		})
		ast2, e2 = parser.Parse("1 + 1")
		So(ast1, shouldEqualAST, ast2)
		So(e1, ShouldEqualErrors, e2)
	})
}
