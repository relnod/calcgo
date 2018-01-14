package calculator_test

import (
	"testing"

	"github.com/relnod/calcgo/interpreter/calculator"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCalculator(t *testing.T) {
	Convey("Calculator Spec", t, func() {
		Convey("ConvertInteger()", func() {
			Convey("correct integer", func() {
				result, err := calculator.ConvertInteger("123")
				So(result, ShouldEqual, 123)
				So(err, ShouldBeNil)
			})

			Convey("invalid character", func() {
				result, err := calculator.ConvertInteger("a")
				So(result, ShouldEqual, 0)
				So(err, ShouldEqual, calculator.ErrorInvalidInteger)
			})

			Convey("overflow", func() {
				result, err := calculator.ConvertInteger("11111111111111111111111111111")
				So(result, ShouldEqual, 0)
				So(err, ShouldEqual, calculator.ErrorInvalidInteger)
			})
		})

		Convey("ConvertDecimal()", func() {
			Convey("correct decimal", func() {
				result, err := calculator.ConvertDecimal("1.23")
				So(result, ShouldEqual, 1.23)
				So(err, ShouldBeNil)
			})

			Convey("invalid character", func() {
				result, err := calculator.ConvertDecimal("a")
				So(result, ShouldEqual, 0)
				So(err, ShouldEqual, calculator.ErrorInvalidDecimal)
			})

			Convey("overflow", func() {
				result, err := calculator.ConvertDecimal("1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111.1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111")
				So(result, ShouldEqual, 0)
				So(err, ShouldEqual, calculator.ErrorInvalidDecimal)
			})
		})

		Convey("ConvertHex()", func() {
			Convey("correct hex", func() {
				result, err := calculator.ConvertHex("0x1A")
				So(result, ShouldEqual, 26)
				So(err, ShouldBeNil)
			})

			Convey("invalid character", func() {
				result, err := calculator.ConvertHex("0xH")
				So(result, ShouldEqual, 0)
				So(err, ShouldEqual, calculator.ErrorInvalidHexadecimal)
			})

			Convey("overflow", func() {
				result, err := calculator.ConvertHex("0x111111111111111111111111")
				So(result, ShouldEqual, 0)
				So(err, ShouldEqual, calculator.ErrorInvalidHexadecimal)
			})
		})

		Convey("ConvertBin()", func() {
			Convey("correct bin", func() {
				result, err := calculator.ConvertBin("0b101")
				So(result, ShouldEqual, 5)
				So(err, ShouldBeNil)
			})

			Convey("invalid character", func() {
				result, err := calculator.ConvertBin("0b2")
				So(result, ShouldEqual, 0)
				So(err, ShouldEqual, calculator.ErrorInvalidBinary)
			})

			Convey("overflow", func() {
				result, err := calculator.ConvertBin("0b1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111")
				So(result, ShouldEqual, 0)
				So(err, ShouldEqual, calculator.ErrorInvalidBinary)
			})
		})

		Convey("ConvertExponential()", func() {
			Convey("correct exponential", func() {
				result, err := calculator.ConvertExponential("2^2")
				So(result, ShouldEqual, 4)
				So(err, ShouldBeNil)
			})

			Convey("invalid character", func() {
				result, err := calculator.ConvertExponential("2^$")
				So(result, ShouldEqual, 0)
				So(err, ShouldEqual, calculator.ErrorInvalidExponential)
			})

			SkipConvey("overflow", func() {
				result, err := calculator.ConvertExponential("999999^999")
				So(result, ShouldEqual, 0)
				So(err, ShouldEqual, calculator.ErrorInvalidExponential)
			})
		})
	})
}
