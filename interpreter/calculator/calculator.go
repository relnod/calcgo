package calculator

import (
	"errors"
	"math"
	"strconv"

	"github.com/relnod/calcgo/parser"
)

// Errors that can occur during calculation or conversion.
var (
	ErrorInvalidInteger = errors.New("Invalid Integer")
	ErrorInvalidDecimal = errors.New("Invalid Decimal")
	ErrorDivisionByZero = errors.New("Division by zero")
)

// ConvertInteger converts a string to a float64.
// Returns an error if conversion failed.
func ConvertInteger(value string) (float64, error) {
	integer, err := strconv.Atoi(value)
	if err != nil {
		return 0, ErrorInvalidInteger
	}
	return float64(integer), nil
}

// ConvertDecimal converts a string to a float64.
// Returns an error if conversion failed.
func ConvertDecimal(value string) (float64, error) {
	decimal, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, ErrorInvalidDecimal
	}
	return decimal, nil
}

// CalculateOperator calculates the result of an operator.
func CalculateOperator(left, right float64, nodeType parser.NodeType) (float64, error) {
	var result float64

	switch nodeType {
	case parser.NAddition:
		result = left + right
	case parser.NSubtraction:
		result = left - right
	case parser.NMultiplication:
		result = left * right
	case parser.NDivision:
		if right == 0 {
			return 0, ErrorDivisionByZero
		}
		result = left / right
	}

	return result, nil
}

// CalculateFunction calculates the result of a function.
func CalculateFunction(arg float64, nodeType parser.NodeType) (float64, error) {
	var result float64

	switch nodeType {
	case parser.NFuncSqrt:
		result = math.Sqrt(arg)
	}

	return result, nil
}
