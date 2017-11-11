package calcgo_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/relnod/calcgo"
)

func TestInterpreter(t *testing.T) {
	Convey("calcgo.Interpreter works with", t, func() {
		Convey("simple integers", func() {
			So(calcgo.Interpret("1"), ShouldEqual, 1)
			So(calcgo.Interpret("12345"), ShouldEqual, 12345)
		})

		Convey("simple decimals", func() {
			So(calcgo.Interpret("1.0"), ShouldEqual, 1.0)
			So(calcgo.Interpret("1234.5678"), ShouldEqual, 1234.5678)
		})

		Convey("simple additions with integers", func() {
			So(calcgo.Interpret("1 + 1"), ShouldEqual, 2)
			So(calcgo.Interpret("3 + 5"), ShouldEqual, 3+5)
			So(calcgo.Interpret("1 + 2 + 3 + 4 + 5 + 6"), ShouldEqual, 1+2+3+4+5+6)
		})

		Convey("simple additions with decimals", func() {
			SkipSo(calcgo.Interpret("1.2 + 2.4"), ShouldEqual, 1.2+2.4) // @todo: fix rounding error
			SkipSo(calcgo.Interpret("0.7 + 2.4"), ShouldEqual, 0.7+2.4) // @todo: fix rounding error
			So(calcgo.Interpret("3.5 + 5.1"), ShouldEqual, 3.5+5.1)
		})

		Convey("simple subtractions", func() {
			So(calcgo.Interpret("1 - 1"), ShouldEqual, 1-1)
			So(calcgo.Interpret("3 - 5"), ShouldEqual, 3-5)
			So(calcgo.Interpret("1 - 2 - 3 - 4 - 5 - 6"), ShouldEqual, 1-2-3-4-5-6)
		})

		Convey("simple multiplications", func() {
			So(calcgo.Interpret("1 * 1"), ShouldEqual, 1*1)
			So(calcgo.Interpret("3 * 5"), ShouldEqual, 3*5)
			So(calcgo.Interpret("1 * 2 * 3 * 4 * 5 * 6"), ShouldEqual, 1*2*3*4*5*6)
		})

		Convey("simple divisions", func() {
			So(calcgo.Interpret("1 / 1"), ShouldEqual, 1/1)
			So(calcgo.Interpret("3 / 5"), ShouldEqual, 3.0/5.0)
			So(calcgo.Interpret("1 / 2 / 3 / 4 / 5 / 6"), ShouldEqual, 1.0/2.0/3.0/4.0/5.0/6.0)
		})

		Convey("dot before line rule", func() {
			SkipSo(calcgo.Interpret("1 + 2 / 3"), ShouldEqual, 1.0+2.0/3.0) // @todo: fix rounding error
			SkipSo(calcgo.Interpret("1 - 2 / 3"), ShouldEqual, 1.0-2.0/3.0) // @todo: fix rounding error
			So(calcgo.Interpret("1 + 2 * 3"), ShouldEqual, 1.0+2.0*3.0)
			So(calcgo.Interpret("1 - 2 * 3"), ShouldEqual, 1.0-2.0*3.0)
		})

		Convey("brackets", func() {
			So(calcgo.Interpret("(1 + 2) / 3"), ShouldEqual, (1.0+2.0)/3.0)
			So(calcgo.Interpret("(1 - 2) / 3"), ShouldEqual, (1.0-2.0)/3.0)
			So(calcgo.Interpret("(1 + 2) * 3"), ShouldEqual, (1.0+2.0)*3.0)
			So(calcgo.Interpret("(1 - 2) * 3"), ShouldEqual, (1.0-2.0)*3.0)
			So(calcgo.Interpret("2 + (1 - 2) / 3"), ShouldEqual, 2.0+(1.0-2.0)/3.0)
		})

		SkipConvey("nested brackets", func() {
			So(calcgo.Interpret("((1 + 2) / 3) + 1"), ShouldEqual, ((1.0+2.0)/3.0)+1)
			So(calcgo.Interpret("((2 + 3) / (1 + 2)) * 3"), ShouldEqual, ((2.0+3.0)/(1.0+2.0))*3.0)
			So(calcgo.Interpret("(1 - 2) * (3 - 2) / (1 + 4)"), ShouldEqual, (1.0-2.0)*(3.0-2.0)/(1.0+4.0))
			SkipSo(calcgo.Interpret("(1 + 2) * 3 + (4 - 6 / (5 + 2))"), ShouldEqual, (1+2)*3+(4-6/(5+2)))
		})

		Convey("brackets and dot before line rule", func() {
			So(calcgo.Interpret("1 + (1 + 2) * 3"), ShouldEqual, 1.0+(1.0+2.0)*3.0)
			So(calcgo.Interpret("1 + (1 + 2) / 3"), ShouldEqual, 1.0+(1.0+2.0)/3.0)
			So(calcgo.Interpret("1 - (1 + 2) * 3"), ShouldEqual, 1.0-(1.0+2.0)*3.0)
			So(calcgo.Interpret("1 - (1 + 2) / 3"), ShouldEqual, 1.0-(1.0+2.0)/3.0)
		})

	})
}
