package parser_test

import (
	"testing"

	"github.com/relnod/calcgo/parser"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAST(t *testing.T) {
	var cases = []struct {
		nt         parser.NodeType
		isOperator bool
	}{
		{parser.NDec, false},
		{parser.NInt, false},
		{parser.NVar, false},
		{parser.NError, false},
		{parser.NAdd, true},
		{parser.NSub, true},
		{parser.NMult, true},
		{parser.NDiv, true},
	}

	Convey("IsOperator works", t, func() {
		for _, c := range cases {
			n := parser.Node{Type: c.nt}
			So(n.IsOperator(), ShouldEqual, c.isOperator)
		}
	})
}
