package parser_test

import (
	"testing"

	"github.com/relnod/calcgo/parser"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAST(t *testing.T) {
	Convey("AST spec", t, func() {
		Convey("IsLiteral()", func() {
			var cases = []struct {
				nt         parser.NodeType
				isOperator bool
			}{
				{parser.NError, false},
				{parser.NInt, true},
				{parser.NVar, true},
				{parser.NAdd, false},
			}

			for _, c := range cases {
				n := parser.Node{Type: c.nt}
				So(n.IsLiteral(), ShouldEqual, c.isOperator)
			}
		})
		Convey("IsOperator()", func() {
			var cases = []struct {
				nt         parser.NodeType
				isOperator bool
			}{
				{parser.NVar, false},
				{parser.NAdd, true},
				{parser.NSub, true},
				{parser.NMult, true},
				{parser.NDiv, true},
				{parser.NMod, true},
				{parser.NOr, true},
				{parser.NXor, true},
				{parser.NAnd, true},
				{parser.NFnSqrt, false},
			}

			for _, c := range cases {
				n := parser.Node{Type: c.nt}
				So(n.IsOperator(), ShouldEqual, c.isOperator)
			}
		})
		Convey("IsFunction()", func() {
			var cases = []struct {
				nt         parser.NodeType
				isOperator bool
			}{
				{parser.NDiv, false},
				{parser.NFnCos, true},
				{parser.NFnSin, true},
			}

			for _, c := range cases {
				n := parser.Node{Type: c.nt}
				So(n.IsFunction(), ShouldEqual, c.isOperator)
			}
		})
	})
}
