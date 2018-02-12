package token_test

import (
	"testing"

	"github.com/relnod/calcgo/token"
	. "github.com/smartystreets/goconvey/convey"
)

func TestToken(t *testing.T) {
	Convey("Token Spec", t, func() {
		Convey("IsLiteral()", func() {
			var cases = []struct {
				t token.TokenType
				r bool
			}{
				{token.TEOF, false},
				{token.TInt, true},
				{token.TDec, true},
				{token.THex, true},
				{token.TBin, true},
				{token.TExp, true},
				{token.TVar, true},
				{token.TOpPlus, false},
			}

			for _, c := range cases {
				n := token.Token{Type: c.t}
				So(n.IsLiteral(), ShouldEqual, c.r)
			}
		})

		Convey("IsOperator()", func() {
			var cases = []struct {
				t token.TokenType
				r bool
			}{
				{token.TVar, false},
				{token.TOpPlus, true},
				{token.TOpMinus, true},
				{token.TOpMult, true},
				{token.TOpDiv, true},
				{token.TOpOr, true},
				{token.TOpXor, true},
				{token.TOpAnd, true},
				{token.TLParen, false},
			}

			for _, c := range cases {
				n := token.Token{Type: c.t}
				So(n.IsOperator(), ShouldEqual, c.r)
			}
		})

		Convey("IsFunction()", func() {
			var cases = []struct {
				t token.TokenType
				r bool
			}{
				{token.TOpDiv, false},
				{token.TFnSqrt, true},
				{token.TFnSin, true},
				{token.TFnCos, true},
				{token.TFnTan, true},
				{token.TLParen, false},
			}

			for _, c := range cases {
				n := token.Token{Type: c.t}
				So(n.IsFunction(), ShouldEqual, c.r)
			}
		})

		Convey("String()", func() {
			var cases = []struct {
				t token.TokenType
				v string

				r string
			}{
				{token.TOpDiv, "/", `{Value: '/', Type: '/', Start: '0', End: '0', }`},
				{255, "bla", `{Value: 'bla', Type: 'Unknown token', Start: '0', End: '0', }`},
			}

			for _, c := range cases {
				n := token.Token{Type: c.t, Value: c.v}
				So(n.String(), ShouldEqual, c.r)
			}
		})
	})
}
