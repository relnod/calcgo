package interpreter_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/relnod/calcgo/interpreter"
	"github.com/relnod/calcgo/parser"
	. "github.com/smartystreets/goconvey/convey"
)

func oastToString(ast *interpreter.OptimizedAST) string {
	str, err := json.MarshalIndent(ast, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	return string(str) + "\n"
}

func optimizeError(actual, expected *interpreter.OptimizedAST) string {
	return "Expected: \n" +
		oastToString(expected) +
		"Actual: \n" +
		oastToString(actual)
}

func eqOptimizedNodes(n1, n2 *interpreter.OptimizedNode) bool {
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

	if n1.OldValue != n2.OldValue {
		return false
	}

	if n1.IsOptimized != n2.IsOptimized {
		return false
	}

	if !eqOptimizedNodes(n1.LeftChild, n2.LeftChild) {
		return false
	}

	if !eqOptimizedNodes(n1.RightChild, n2.RightChild) {
		return false
	}

	return true
}

func eqOptimizedAST(oast1, oast2 *interpreter.OptimizedAST) bool {
	if oast1 != nil && oast2 != nil {
		return true
	}

	if oast1 != nil && oast2 == nil {
		return false
	}

	if oast1 == nil && oast2 != nil {
		return false
	}

	return eqOptimizedNodes(oast1.Node, oast2.Node)
}

func ShouldEqualOptimizedAST(actual interface{}, expected ...interface{}) string {
	actualOAST := actual.(*interpreter.OptimizedAST)
	expectedOAST := expected[0].(*interpreter.OptimizedAST)

	if eqOptimizedAST(actualOAST, expectedOAST) {
		return ""
	}

	return optimizeError(actualOAST, expectedOAST) + "(Should be Equal)"
}

func optimize(str string) (*interpreter.OptimizedAST, error) {
	ast, errors := parser.Parse(str)
	if errors != nil {
		return nil, errors[0]
	}

	return interpreter.Optimize(&ast)
}

func TestOptimizer(t *testing.T) {
	Convey("Optimizer works with", t, func() {
		Convey("nothing", func() {
			oast, err := interpreter.Optimize(nil)
			So(oast, ShouldBeNil)
			So(err, ShouldBeNil)
		})

		Convey("ast without variable", func() {
			Convey("integer number", func() {
				oast, err := optimize("1")
				So(oast, ShouldEqualOptimizedAST, &interpreter.OptimizedAST{
					Node: &interpreter.OptimizedNode{
						Type:        parser.NDecimal,
						Value:       1.0,
						OldValue:    "",
						IsOptimized: true,
						LeftChild:   nil,
						RightChild:  nil,
					},
				})
				So(err, ShouldBeNil)
			})

			Convey("decimal number", func() {
				oast, err := optimize("1.3")
				So(oast, ShouldEqualOptimizedAST, &interpreter.OptimizedAST{
					Node: &interpreter.OptimizedNode{
						Type:        parser.NDecimal,
						Value:       1.3,
						OldValue:    "",
						IsOptimized: true,
						LeftChild:   nil,
						RightChild:  nil,
					},
				})
				So(err, ShouldBeNil)
			})

			Convey("addition", func() {
				oast, err := optimize("1 + 1")
				So(oast, ShouldEqualOptimizedAST, &interpreter.OptimizedAST{
					Node: &interpreter.OptimizedNode{
						Type:        parser.NDecimal,
						Value:       2.0,
						OldValue:    "",
						IsOptimized: true,
						LeftChild:   nil,
						RightChild:  nil,
					},
				})
				So(err, ShouldBeNil)
			})

			Convey("subtraction", func() {
				oast, err := optimize("1 - 1")
				So(oast, ShouldEqualOptimizedAST, &interpreter.OptimizedAST{
					Node: &interpreter.OptimizedNode{
						Type:        parser.NDecimal,
						Value:       0.0,
						OldValue:    "",
						IsOptimized: true,
						LeftChild:   nil,
						RightChild:  nil,
					},
				})
				So(err, ShouldBeNil)
			})

			Convey("multiplication", func() {
				oast, err := optimize("1 * 1")
				So(oast, ShouldEqualOptimizedAST, &interpreter.OptimizedAST{
					Node: &interpreter.OptimizedNode{
						Type:        parser.NDecimal,
						Value:       1.0,
						OldValue:    "",
						IsOptimized: true,
						LeftChild:   nil,
						RightChild:  nil,
					},
				})
				So(err, ShouldBeNil)
			})

			Convey("division", func() {
				oast, err := optimize("1 / 1")
				So(oast, ShouldEqualOptimizedAST, &interpreter.OptimizedAST{
					Node: &interpreter.OptimizedNode{
						Type:        parser.NDecimal,
						Value:       1.0,
						OldValue:    "",
						IsOptimized: true,
						LeftChild:   nil,
						RightChild:  nil,
					},
				})
				So(err, ShouldBeNil)
			})
		})

		Convey("ast with variable", func() {
			Convey("only vairable", func() {
				oast, err := optimize("a")
				So(oast, ShouldEqualOptimizedAST, &interpreter.OptimizedAST{
					Node: &interpreter.OptimizedNode{
						Type:        parser.NDecimal,
						Value:       0,
						OldValue:    "a",
						IsOptimized: false,
						LeftChild:   nil,
						RightChild:  nil,
					},
				})
				So(err, ShouldBeNil)
			})

			Convey("addition", func() {
				oast, err := optimize("1 + a")
				So(oast, ShouldEqualOptimizedAST, &interpreter.OptimizedAST{
					Node: &interpreter.OptimizedNode{
						Type:        parser.NAddition,
						Value:       0,
						OldValue:    "",
						IsOptimized: false,
						LeftChild: &interpreter.OptimizedNode{
							Type:        parser.NDecimal,
							Value:       1.0,
							OldValue:    "",
							IsOptimized: true,
							LeftChild:   nil,
							RightChild:  nil,
						},
						RightChild: &interpreter.OptimizedNode{
							Type:        parser.NVariable,
							Value:       0,
							OldValue:    "a",
							IsOptimized: false,
							LeftChild:   nil,
							RightChild:  nil,
						},
					},
				})
				So(err, ShouldBeNil)
			})

			Convey("subtraction", func() {
				oast, err := optimize("1 - a")
				So(oast, ShouldEqualOptimizedAST, &interpreter.OptimizedAST{
					Node: &interpreter.OptimizedNode{
						Type:        parser.NSubtraction,
						Value:       0,
						OldValue:    "",
						IsOptimized: false,
						LeftChild: &interpreter.OptimizedNode{
							Type:        parser.NDecimal,
							Value:       1.0,
							OldValue:    "",
							IsOptimized: true,
							LeftChild:   nil,
							RightChild:  nil,
						},
						RightChild: &interpreter.OptimizedNode{
							Type:        parser.NVariable,
							Value:       0,
							OldValue:    "a",
							IsOptimized: false,
							LeftChild:   nil,
							RightChild:  nil,
						},
					},
				})
				So(err, ShouldBeNil)
			})

			Convey("multiplication", func() {
				oast, err := optimize("1 - a")
				So(oast, ShouldEqualOptimizedAST, &interpreter.OptimizedAST{
					Node: &interpreter.OptimizedNode{
						Type:        parser.NMultiplication,
						Value:       0,
						OldValue:    "",
						IsOptimized: false,
						LeftChild: &interpreter.OptimizedNode{
							Type:        parser.NDecimal,
							Value:       1.0,
							OldValue:    "",
							IsOptimized: true,
							LeftChild:   nil,
							RightChild:  nil,
						},
						RightChild: &interpreter.OptimizedNode{
							Type:        parser.NVariable,
							Value:       0,
							OldValue:    "a",
							IsOptimized: false,
							LeftChild:   nil,
							RightChild:  nil,
						},
					},
				})
				So(err, ShouldBeNil)
			})

			Convey("division", func() {
				oast, err := optimize("1 / a")
				So(oast, ShouldEqualOptimizedAST, &interpreter.OptimizedAST{
					Node: &interpreter.OptimizedNode{
						Type:        parser.NDivision,
						Value:       0,
						OldValue:    "",
						IsOptimized: false,
						LeftChild: &interpreter.OptimizedNode{
							Type:        parser.NDecimal,
							Value:       1.0,
							OldValue:    "",
							IsOptimized: true,
							LeftChild:   nil,
							RightChild:  nil,
						},
						RightChild: &interpreter.OptimizedNode{
							Type:        parser.NVariable,
							Value:       0,
							OldValue:    "a",
							IsOptimized: false,
							LeftChild:   nil,
							RightChild:  nil,
						},
					},
				})
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Optimizer returns an error when", t, func() {
		Convey("dividing by 0", func() {
			oast, err := optimize("1 / 0")
			So(oast, ShouldBeNil)
			So(err, ShouldEqual, interpreter.ErrorDivisionByZero)
		})

		Convey("interpreting wrong number", func() {
			oast, err := interpreter.Optimize(&parser.AST{
				Node: &parser.Node{
					Type:       parser.NInteger,
					Value:      "a",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(oast, ShouldBeNil)
			So(err, ShouldEqual, interpreter.ErrorInvalidInteger)

			oast, err = interpreter.Optimize(&parser.AST{
				Node: &parser.Node{
					Type:       parser.NDecimal,
					Value:      "a",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(oast, ShouldBeNil)
			So(err, ShouldEqual, interpreter.ErrorInvalidDecimal)
		})

		Convey("invalid node type", func() {
			oast, err := interpreter.Optimize(&parser.AST{
				Node: &parser.Node{
					Type:       3000,
					Value:      "",
					LeftChild:  nil,
					RightChild: nil,
				},
			})
			So(oast, ShouldBeNil)
			So(err, ShouldEqual, interpreter.ErrorInvalidNodeType)
		})

		Convey("left child is missing", func() {
			oast, err := interpreter.Optimize(&parser.AST{
				Node: &parser.Node{
					Type:      parser.NAddition,
					Value:     "",
					LeftChild: nil,
					RightChild: &parser.Node{
						Type:       3000,
						Value:      "",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(oast, ShouldBeNil)
			So(err, ShouldEqual, interpreter.ErrorMissingLeftChild)
		})

		Convey("right child is missing", func() {
			oast, err := interpreter.Optimize(&parser.AST{
				Node: &parser.Node{
					Type:  parser.NAddition,
					Value: "",
					LeftChild: &parser.Node{
						Type:       3000,
						Value:      "",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: nil,
				},
			})
			So(oast, ShouldBeNil)
			So(err, ShouldEqual, interpreter.ErrorMissingRightChild)
		})

		Convey("error happens in nested node", func() {
			oast, err := interpreter.Optimize(&parser.AST{
				Node: &parser.Node{
					Type:  parser.NAddition,
					Value: "",
					LeftChild: &parser.Node{
						Type:       3000,
						Value:      "",
						LeftChild:  nil,
						RightChild: nil,
					},
					RightChild: &parser.Node{
						Type:       3000,
						Value:      "",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(oast, ShouldBeNil)
			So(err, ShouldEqual, interpreter.ErrorInvalidNodeType)

			oast, err = interpreter.Optimize(&parser.AST{
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
						Type:       3000,
						Value:      "",
						LeftChild:  nil,
						RightChild: nil,
					},
				},
			})
			So(oast, ShouldBeNil)
			So(err, ShouldEqual, interpreter.ErrorInvalidNodeType)
		})
	})
}
