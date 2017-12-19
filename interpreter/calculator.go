package interpreter

import (
	"math"

	"github.com/relnod/calcgo/parser"
)

func calculateOperator(left, right float64, nodeType parser.NodeType) (float64, error) {
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

func calculateFunction(arg float64, nodeType parser.NodeType) (float64, error) {
	var result float64

	switch nodeType {
	case parser.NFuncSqrt:
		result = math.Sqrt(arg)
	}

	return result, nil
}
