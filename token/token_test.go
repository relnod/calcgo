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
				t token.Type
				r bool
			}{
				{token.EOF, false},
				{token.Int, true},
				{token.Dec, true},
				{token.Hex, true},
				{token.Bin, true},
				{token.Exp, true},
				{token.Var, true},
				{token.Plus, false},
			}

			for _, c := range cases {
				n := token.Token{Type: c.t}
				So(n.IsLiteral(), ShouldEqual, c.r)
			}
		})

		Convey("IsOperator()", func() {
			var cases = []struct {
				t token.Type
				r bool
			}{
				{token.Var, false},
				{token.Plus, true},
				{token.Minus, true},
				{token.Mult, true},
				{token.Div, true},
				{token.Or, true},
				{token.Xor, true},
				{token.And, true},
				{token.ParenL, false},
			}

			for _, c := range cases {
				n := token.Token{Type: c.t}
				So(n.IsOperator(), ShouldEqual, c.r)
			}
		})

		Convey("IsFunction()", func() {
			var cases = []struct {
				t token.Type
				r bool
			}{
				{token.Div, false},
				{token.Sqrt, true},
				{token.Sin, true},
				{token.Cos, true},
				{token.Tan, true},
				{token.ParenL, false},
			}

			for _, c := range cases {
				n := token.Token{Type: c.t}
				So(n.IsFunction(), ShouldEqual, c.r)
			}
		})

		Convey("String()", func() {
			var cases = []struct {
				t token.Type
				v string

				r string
			}{
				{token.Div, "/", `{Value: '/', Type: '/', Start: '0', End: '0', }`},
				{255, "bla", `{Value: 'bla', Type: 'Unknown token', Start: '0', End: '0', }`},
			}

			for _, c := range cases {
				n := token.Token{Type: c.t, Value: c.v}
				So(n.String(), ShouldEqual, c.r)
			}
		})
	})
}
