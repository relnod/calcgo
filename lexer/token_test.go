package lexer_test

import (
	"testing"

	"github.com/relnod/calcgo/lexer"
	. "github.com/smartystreets/goconvey/convey"
)

func TestToken(t *testing.T) {
	Convey("Token Spec", t, func() {
		Convey("IsLiteral()", func() {
			var cases = []struct {
				t lexer.TokenType
				r bool
			}{
				{lexer.TEOF, false},
				{lexer.TInt, true},
				{lexer.TDec, true},
				{lexer.THex, true},
				{lexer.TBin, true},
				{lexer.TExp, true},
				{lexer.TVar, true},
				{lexer.TOpPlus, false},
			}

			for _, c := range cases {
				n := lexer.Token{Type: c.t}
				So(n.IsLiteral(), ShouldEqual, c.r)
			}
		})

		Convey("IsOperator()", func() {
			var cases = []struct {
				t lexer.TokenType
				r bool
			}{
				{lexer.TVar, false},
				{lexer.TOpPlus, true},
				{lexer.TOpMinus, true},
				{lexer.TOpMult, true},
				{lexer.TOpDiv, true},
				{lexer.TOpOr, true},
				{lexer.TOpXor, true},
				{lexer.TOpAnd, true},
				{lexer.TLParen, false},
			}

			for _, c := range cases {
				n := lexer.Token{Type: c.t}
				So(n.IsOperator(), ShouldEqual, c.r)
			}
		})

		Convey("IsFunction()", func() {
			var cases = []struct {
				t lexer.TokenType
				r bool
			}{
				{lexer.TOpDiv, false},
				{lexer.TFnSqrt, true},
				{lexer.TFnSin, true},
				{lexer.TFnCos, true},
				{lexer.TFnTan, true},
				{lexer.TLParen, false},
			}

			for _, c := range cases {
				n := lexer.Token{Type: c.t}
				So(n.IsFunction(), ShouldEqual, c.r)
			}
		})

		Convey("String()", func() {
			var cases = []struct {
				t lexer.TokenType
				v string

				r string
			}{
				{lexer.TOpDiv, "/", `{"/", /}`},
				{255, "bla", `{"bla", Unknown token}`},
			}

			for _, c := range cases {
				n := lexer.Token{Type: c.t, Value: c.v}
				So(n.String(), ShouldEqual, c.r)
			}
		})
	})
}
