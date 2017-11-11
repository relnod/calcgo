package calcgo_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/relnod/calcgo"
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

	if n1 == nil && n1 != nil {
		return false
	}

	if n1 != nil && n1 == nil {
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
	Convey("Parser works with", t, func() {
		Convey("nothing", func() {
			So(calcgo.Parse(""), shouldEqualAST, calcgo.AST{})
		})

		Convey("simple numbers", func() {
			So(calcgo.Parse("20"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:       calcgo.NInteger,
					Value:      "20",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(calcgo.Parse("20.23"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:       calcgo.NDecimal,
					Value:      "20.23",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
		})

		Convey("additions", func() {
			So(calcgo.Parse("1 + 2"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NAddition,
					Value: "+",
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
			So(calcgo.Parse("1 + 2 + 3"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NAddition,
					Value: "+",
					LeftChild: &calcgo.Node{
						Type:  calcgo.NAddition,
						Value: "+",
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
		})

		Convey("addition with decimals", func() {
			So(calcgo.Parse("1.2 + 2.4"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NAddition,
					Value: "+",
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
		})

		Convey("subtractions", func() {
			So(calcgo.Parse("1 - 2"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NSubtraction,
					Value: "-",
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
			So(calcgo.Parse("1 - 2 - 3"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NSubtraction,
					Value: "-",
					LeftChild: &calcgo.Node{
						Type:  calcgo.NSubtraction,
						Value: "-",
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
		})

		Convey("Multiplications", func() {
			So(calcgo.Parse("1 * 2"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NMultiplication,
					Value: "*",
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
			So(calcgo.Parse("1 * 2 * 3"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NMultiplication,
					Value: "*",
					LeftChild: &calcgo.Node{
						Type:  calcgo.NMultiplication,
						Value: "*",
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
		})

		Convey("Divisions", func() {
			So(calcgo.Parse("1 / 2"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NDivision,
					Value: "/",
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
			So(calcgo.Parse("1 / 2 / 3"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NDivision,
					Value: "/",
					LeftChild: &calcgo.Node{
						Type:  calcgo.NDivision,
						Value: "/",
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
		})

		Convey("dot before line", func() {
			So(calcgo.Parse("1 + 2 * 3"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NAddition,
					Value: "+",
					LeftChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &calcgo.Node{
						Type:  calcgo.NMultiplication,
						Value: "*",
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
			So(calcgo.Parse("1 - 2 / 3"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NSubtraction,
					Value: "-",
					LeftChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "1",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &calcgo.Node{
						Type:  calcgo.NDivision,
						Value: "/",
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

			So(calcgo.Parse("2 * 3 + 1"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NAddition,
					Value: "+",
					LeftChild: &calcgo.Node{
						Type:  calcgo.NMultiplication,
						Value: "*",
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
		})

		Convey("brackets", func() {
			So(calcgo.Parse("(1)"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:       calcgo.NInteger,
					Value:      "1",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(calcgo.Parse("(1 - 2)"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NSubtraction,
					Value: "-",
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
			So(calcgo.Parse("(1 - 2) * 3"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NMultiplication,
					Value: "*",
					LeftChild: &calcgo.Node{
						Type:  calcgo.NSubtraction,
						Value: "-",
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
			So(calcgo.Parse("3 * (1 - 2)"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NMultiplication,
					Value: "*",
					LeftChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "3",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &calcgo.Node{
						Type:  calcgo.NSubtraction,
						Value: "-",
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
			So(calcgo.Parse("3 * (1 - 2) / 4"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NDivision,
					Value: "/",
					LeftChild: &calcgo.Node{
						Type:  calcgo.NMultiplication,
						Value: "*",
						LeftChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "3",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &calcgo.Node{
							Type:  calcgo.NSubtraction,
							Value: "-",
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
		})

		Convey("nested brackets", func() {
			So(calcgo.Parse("((1))"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:       calcgo.NInteger,
					Value:      "1",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(calcgo.Parse("((1 - 2))"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NSubtraction,
					Value: "-",
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
			So(calcgo.Parse("(3 * (1 - 2))"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NMultiplication,
					Value: "*",
					LeftChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "3",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &calcgo.Node{
						Type:  calcgo.NSubtraction,
						Value: "-",
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
		})

		Convey("Brackets with dot before line rule", func() {
			So(calcgo.Parse("3 + (1 - 2) / 4"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NAddition,
					Value: "+",
					LeftChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "3",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &calcgo.Node{
						Type:  calcgo.NDivision,
						Value: "/",
						LeftChild: &calcgo.Node{
							Type:  calcgo.NSubtraction,
							Value: "-",
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

			So(calcgo.Parse("3 + (1 + 2) * 4"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NAddition,
					Value: "+",
					LeftChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "3",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &calcgo.Node{
						Type:  calcgo.NMultiplication,
						Value: "*",
						LeftChild: &calcgo.Node{
							Type:  calcgo.NAddition,
							Value: "+",
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

			So(calcgo.Parse("(1 + 2) * 4 + 1"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NAddition,
					Value: "+",
					LeftChild: &calcgo.Node{
						Type:  calcgo.NMultiplication,
						Value: "*",
						LeftChild: &calcgo.Node{
							Type:  calcgo.NAddition,
							Value: "+",
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

			So(calcgo.Parse("4 - 6 / (5 + 2)"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NSubtraction,
					Value: "-",
					LeftChild: &calcgo.Node{
						Type:       calcgo.NInteger,
						Value:      "4",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &calcgo.Node{
						Type:  calcgo.NDivision,
						Value: "/",
						LeftChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "6",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &calcgo.Node{
							Type:  calcgo.NAddition,
							Value: "+",
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

			So(calcgo.Parse("(1 + 2) * 3 + (4 - 6 / (5 + 2))"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node{
					Type:  calcgo.NAddition,
					Value: "+",
					LeftChild: &calcgo.Node{
						Type:  calcgo.NMultiplication,
						Value: "*",
						LeftChild: &calcgo.Node{
							Type:  calcgo.NAddition,
							Value: "+",
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
						Value: "-",
						LeftChild: &calcgo.Node{
							Type:       calcgo.NInteger,
							Value:      "4",
							LeftChild:  nil,
							RightChild: nil,
						},
						RightChild: &calcgo.Node{
							Type:  calcgo.NDivision,
							Value: "/",
							LeftChild: &calcgo.Node{
								Type:       calcgo.NInteger,
								Value:      "6",
								LeftChild:  nil,
								RightChild: nil,
							},
							RightChild: &calcgo.Node{
								Type:  calcgo.NAddition,
								Value: "+",
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
		})
	})
}
