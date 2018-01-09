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
		{parser.NDecimal, false},
		{parser.NInteger, false},
		{parser.NVariable, false},
		{parser.NError, false},
		{parser.NAddition, true},
		{parser.NSubtraction, true},
		{parser.NMultiplication, true},
		{parser.NDivision, true},
	}

	Convey("IsOperator works", t, func() {
		for _, c := range cases {
			n := parser.Node{Type: c.nt}
			So(n.IsOperator(), ShouldEqual, c.isOperator)
		}
	})
}
