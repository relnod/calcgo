package calculator_test

import (
	"math"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/relnod/calcgo/interpreter/calculator"
	"github.com/relnod/calcgo/parser"
)

var (
	intOverflow = "11111111111111111111111111111"
	decOverflow = "1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111.1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111"
	hexOverflow = "0x111111111111111111111111"
	binOverflow = "0b1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111"
	expOverflow = "99999^999"
)

func TestCalculator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Calculator Suite")
}

var _ = DescribeTable("ConvertInteger()",
	func(in string, expRes float64, expErr error) {
		result, err := calculator.ConvertInteger(in)
		Expect(result).To(BeNumerically("==", expRes))
		if expErr != nil {
			Expect(err).To(Equal(expErr))
		} else {
			Expect(err).To(BeNil())
		}
	},
	Entry("works with correct integer", "123", 123.0, nil),
	Entry("handles invalid integer", "a", 0.0, calculator.ErrorInvalidInteger),
	Entry("handles overflow", intOverflow, 0.0, calculator.ErrorInvalidInteger),
)

var _ = DescribeTable("ConvertDecimal()",
	func(in string, expRes float64, expErr error) {
		result, err := calculator.ConvertDecimal(in)
		Expect(result).To(BeNumerically("==", expRes))
		if expErr != nil {
			Expect(err).To(Equal(expErr))
		} else {
			Expect(err).To(BeNil())
		}
	},
	Entry("works with correct decimal", "123.456", 123.456, nil),
	Entry("handles invalid decimal", "a", 0.0, calculator.ErrorInvalidDecimal),
	Entry("handles overflow", decOverflow, 0.0, calculator.ErrorInvalidDecimal),
)

var _ = DescribeTable("ConvertHex()",
	func(in string, expRes float64, expErr error) {
		result, err := calculator.ConvertHex(in)
		Expect(result).To(BeNumerically("==", expRes))
		if expErr != nil {
			Expect(err).To(Equal(expErr))
		} else {
			Expect(err).To(BeNil())
		}
	},
	Entry("works with correct hex", "0x1A", 26.0, nil),
	Entry("handles invalid hex", "0xH", 0.0, calculator.ErrorInvalidHexadecimal),
	Entry("handles overflow", hexOverflow, 0.0, calculator.ErrorInvalidHexadecimal),
)

var _ = DescribeTable("ConvertBin()",
	func(in string, expRes float64, expErr error) {
		result, err := calculator.ConvertBin(in)
		Expect(result).To(BeNumerically("==", expRes))
		if expErr != nil {
			Expect(err).To(Equal(expErr))
		} else {
			Expect(err).To(BeNil())
		}
	},
	Entry("works with correct binary", "0b10", 2.0, nil),
	Entry("handles invalid binary", "0b2", 0.0, calculator.ErrorInvalidBinary),
	Entry("handles overflow", binOverflow, 0.0, calculator.ErrorInvalidBinary),
)

var _ = DescribeTable("ConvertExponential()",
	func(in string, expRes float64, expErr error) {
		result, err := calculator.ConvertExponential(in)
		Expect(result).To(BeNumerically("==", expRes))
		if expErr != nil {
			Expect(err).To(Equal(expErr))
		} else {
			Expect(err).To(BeNil())
		}
	},
	Entry("works with correct exponential", "2^2", 4.0, nil),
	Entry("handles invalid exponential", "2^$", 0.0, calculator.ErrorInvalidExponential),
	Entry("handles overflow", expOverflow, 0.0, calculator.ErrorInvalidExponential),
)

var _ = DescribeTable("CalculateOperator()",
	func(left, right float64, nodeType parser.NodeType, expRes float64, expErr error) {
		result, err := calculator.CalculateOperator(left, right, nodeType)
		Expect(result).To(BeNumerically("==", expRes))
		if expErr != nil {
			Expect(err).To(Equal(expErr))
		} else {
			Expect(err).To(BeNil())
		}
	},
	Entry("add", 1.0, 1.0, parser.NAdd, 2.0, nil),
	Entry("sub", 1.0, 1.0, parser.NSub, 0.0, nil),
	Entry("mult", 1.0, 2.0, parser.NMult, 2.0, nil),
	Entry("div", 1.0, 2.0, parser.NDiv, 0.5, nil),
	Entry("div division by zero", 1.0, 0.0, parser.NDiv, 0.0, calculator.ErrorDivisionByZero),
	Entry("mod 1", 5.0, 6.0, parser.NMod, 5.0, nil),
	Entry("mod 2", 6.0, 5.0, parser.NMod, 1.0, nil),
	Entry("mod 3", 7.0, 2.0, parser.NMod, 1.0, nil),
	Entry("mod 4", 7.0, 7.0, parser.NMod, 0.0, nil),
	Entry("mod 4", 4.0, 2.0, parser.NMod, 0.0, nil),
	Entry("or", 1.0, 1.0, parser.NOr, 1.0, nil),
	Entry("xor", 1.0, 1.0, parser.NXor, 0.0, nil),
	Entry("and", 1.0, 0.0, parser.NAnd, 0.0, nil),
)

var _ = DescribeTable("CalculateFunction()",
	func(arg float64, nodeType parser.NodeType, expRes float64) {
		result, err := calculator.CalculateFunction(arg, nodeType)
		Expect(result).To(BeNumerically("==", expRes))
		Expect(err).To(BeNil())
	},
	Entry("sqrt", 9.0, parser.NFnSqrt, math.Sqrt(9.0)),
	Entry("sin", 9.0, parser.NFnSin, math.Sin(9.0)),
	Entry("cos", 9.0, parser.NFnCos, math.Cos(9.0)),
	Entry("tan", 9.0, parser.NFnTan, math.Tan(9.0)),
)
