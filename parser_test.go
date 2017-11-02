package calcgo_test

import (
	"encoding/json"	
	"fmt"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	
	"gitlab.com/relnod/calcgo"
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

	if (n1.Type != n2.Type) {
		return false
	}

	if (n1.Value != n2.Value) {
		return false
	}

	if (len(n1.Childs) != len(n2.Childs)) {
		return false
	}

	for i := 0; i < len(n1.Childs); i++ {
		if (!eqNodes(n1.Childs[i], n2.Childs[i])) {
			return false
		}
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

func parse(str string) calcgo.AST {
	tokens := calcgo.Lex(str)
	ast := calcgo.Parse(tokens)
	return ast
}

func TestParser(t *testing.T) {
	Convey("Parser works with", t, func() {
		Convey("nothing", func() {
			So(parse(""), shouldEqualAST, calcgo.AST{})
		})

		Convey("simple numbers", func() {
			So(parse("20"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node {
					Type: calcgo.NInteger,
					Value: "20",
					Childs: nil,
				},
			})
			So(parse("20.23"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node {
					Type: calcgo.NDecimal,
					Value: "20.23",
					Childs: nil,
				},
			})
		})

		Convey("additions", func() {
			So(parse("1 + 2"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node {
					Type: calcgo.NAddition,
					Value: "+",
					Childs: []*calcgo.Node {
						&calcgo.Node {
							Type: calcgo.NInteger,
							Value: "1",
							Childs: nil,
						},
						&calcgo.Node {
							Type: calcgo.NInteger,
							Value: "2",
							Childs: nil,
						},	
					},
				},
			})
			So(parse("1 + 2 + 3"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node {
					Type: calcgo.NAddition,
					Value: "+",
					Childs: []*calcgo.Node {
						&calcgo.Node {
							Type: calcgo.NAddition,
							Value: "+",
							Childs: []*calcgo.Node {
								&calcgo.Node {
									Type: calcgo.NInteger,
									Value: "1",
									Childs: nil,
								},
								&calcgo.Node {
									Type: calcgo.NInteger,
									Value: "2",
									Childs: nil,
								},	
							},
						},
						&calcgo.Node {
							Type: calcgo.NInteger,
							Value: "3",
							Childs: nil,
						},	
					},
				},
			})
		})

		Convey("subtractions", func() {
			So(parse("1 - 2"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node {
					Type: calcgo.NSubtraction,
					Value: "-",
					Childs: []*calcgo.Node {
						&calcgo.Node {
							Type: calcgo.NInteger,
							Value: "1",
							Childs: nil,
						},
						&calcgo.Node {
							Type: calcgo.NInteger,
							Value: "2",
							Childs: nil,
						},	
					},
				},
			})
			So(parse("1 - 2 - 3"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node {
					Type: calcgo.NSubtraction,
					Value: "-",
					Childs: []*calcgo.Node {
						&calcgo.Node {
							Type: calcgo.NSubtraction,
							Value: "-",
							Childs: []*calcgo.Node {
								&calcgo.Node {
									Type: calcgo.NInteger,
									Value: "1",
									Childs: nil,
								},
								&calcgo.Node {
									Type: calcgo.NInteger,
									Value: "2",
									Childs: nil,
								},	
							},
						},
						&calcgo.Node {
							Type: calcgo.NInteger,
							Value: "3",
							Childs: nil,
						},	
					},
				},
			})
		})

		Convey("Multiplications", func() {
			So(parse("1 * 2"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node {
					Type: calcgo.NMultiplication,
					Value: "*",
					Childs: []*calcgo.Node {
						&calcgo.Node {
							Type: calcgo.NInteger,
							Value: "1",
							Childs: nil,
						},
						&calcgo.Node {
							Type: calcgo.NInteger,
							Value: "2",
							Childs: nil,
						},	
					},
				},
			})
			So(parse("1 * 2 * 3"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node {
					Type: calcgo.NMultiplication,
					Value: "*",
					Childs: []*calcgo.Node {
						&calcgo.Node {
							Type: calcgo.NMultiplication,
							Value: "*",
							Childs: []*calcgo.Node {
								&calcgo.Node {
									Type: calcgo.NInteger,
									Value: "1",
									Childs: nil,
								},
								&calcgo.Node {
									Type: calcgo.NInteger,
									Value: "2",
									Childs: nil,
								},	
							},
						},
						&calcgo.Node {
							Type: calcgo.NInteger,
							Value: "3",
							Childs: nil,
						},	
					},
				},
			})
		})

		Convey("Divisions", func() {
			So(parse("1 / 2"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node {
					Type: calcgo.NDivision,
					Value: "/",
					Childs: []*calcgo.Node {
						&calcgo.Node {
							Type: calcgo.NInteger,
							Value: "1",
							Childs: nil,
						},
						&calcgo.Node {
							Type: calcgo.NInteger,
							Value: "2",
							Childs: nil,
						},	
					},
				},
			})
			So(parse("1 / 2 / 3"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node {
					Type: calcgo.NDivision,
					Value: "/",
					Childs: []*calcgo.Node {
						&calcgo.Node {
							Type: calcgo.NDivision,
							Value: "/",
							Childs: []*calcgo.Node {
								&calcgo.Node {
									Type: calcgo.NInteger,
									Value: "1",
									Childs: nil,
								},
								&calcgo.Node {
									Type: calcgo.NInteger,
									Value: "2",
									Childs: nil,
								},	
							},
						},
						&calcgo.Node {
							Type: calcgo.NInteger,
							Value: "3",
							Childs: nil,
						},	
					},
				},
			})
		})

		Convey("dot before line", func() {
			So(parse("1 + 2 * 3"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node {
					Type: calcgo.NAddition,
					Value: "+",
					Childs: []*calcgo.Node {
						&calcgo.Node {
							Type: calcgo.NInteger,
							Value: "1",
							Childs: nil,
						},
						&calcgo.Node {
							Type: calcgo.NMultiplication,
							Value: "*",
							Childs: []*calcgo.Node {
								&calcgo.Node {
									Type: calcgo.NInteger,
									Value: "2",
									Childs: nil,
								},
								&calcgo.Node {
									Type: calcgo.NInteger,
									Value: "3",
									Childs: nil,
								},	
							},
						},
					},
				},
			})
			So(parse("1 - 2 / 3"), shouldEqualAST, calcgo.AST{
				Node: &calcgo.Node {
					Type: calcgo.NSubtraction,
					Value: "-",
					Childs: []*calcgo.Node {
						&calcgo.Node {
							Type: calcgo.NInteger,
							Value: "1",
							Childs: nil,
						},
						&calcgo.Node {
							Type: calcgo.NDivision,
							Value: "/",
							Childs: []*calcgo.Node {
								&calcgo.Node {
									Type: calcgo.NInteger,
									Value: "2",
									Childs: nil,
								},
								&calcgo.Node {
									Type: calcgo.NInteger,
									Value: "3",
									Childs: nil,
								},	
							},
						},	
					},
				},
			})
		})
	})
}