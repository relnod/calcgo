package calcgo_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"gitlab.com/relnod/calcgo"
)

func interpret(str string) float64 {
	tokens := calcgo.Lex(str)
	ast := calcgo.Parse(tokens)
	number := calcgo.Interpret(ast)

	return number
}

func TestInterpreter(t *testing.T) {
	Convey("interpreter works with", t, func() {
		Convey("simple integers", func() {
			So(interpret("1"), ShouldEqual, 1)
			So(interpret("12345"), ShouldEqual, 12345)
		})

		Convey("simple decimals", func() {
			So(interpret("1.0"), ShouldEqual, 1.0)
			So(interpret("1234.5678"), ShouldEqual, 1234.5678)
		})

		Convey("simple additions with integers", func() {
			So(interpret("1 + 1"), ShouldEqual, 2)
			So(interpret("3 + 5"), ShouldEqual, 3+5)
			So(interpret("1 + 2 + 3 + 4 + 5 + 6"), ShouldEqual, 1+2+3+4+5+6)
		})

		Convey("simple additions with decimals", func() {
			SkipSo(interpret("1.2 + 2.4"), ShouldEqual, 1.2+2.4) // @todo: fix rounding error
			SkipSo(interpret("0.7 + 2.4"), ShouldEqual, 0.7+2.4) // @todo: fix rounding error			
			So(interpret("3.5 + 5.1"), ShouldEqual, 3.5+5.1)
		})

		Convey("simple subtractions", func() {
			So(interpret("1 - 1"), ShouldEqual, 1-1)
			So(interpret("3 - 5"), ShouldEqual, 3-5)
			So(interpret("1 - 2 - 3 - 4 - 5 - 6"), ShouldEqual, 1-2-3-4-5-6)
		})

		Convey("simple multiplications", func() {
			So(interpret("1 * 1"), ShouldEqual, 1*1)
			So(interpret("3 * 5"), ShouldEqual, 3*5)
			So(interpret("1 * 2 * 3 * 4 * 5 * 6"), ShouldEqual, 1*2*3*4*5*6)
		})

		Convey("simple divisions", func() {
			So(interpret("1 / 1"), ShouldEqual, 1/1)
			So(interpret("3 / 5"), ShouldEqual, 3.0/5.0)
			So(interpret("1 / 2 / 3 / 4 / 5 / 6"), ShouldEqual, 1.0/2.0/3.0/4.0/5.0/6.0)
		})

		Convey("dot before line rule", func() {
			SkipSo(interpret("1 + 2 / 3"), ShouldEqual, 1.0+2.0/3.0) // @todo: fix rounding error
			SkipSo(interpret("1 - 2 / 3"), ShouldEqual, 1.0-2.0/3.0) // @todo: fix rounding error
			So(interpret("1 + 2 * 3"), ShouldEqual, 1.0+2.0*3.0)
			So(interpret("1 - 2 * 3"), ShouldEqual, 1.0-2.0*3.0)
		})
	})
}
