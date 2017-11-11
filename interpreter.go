package calcgo

import "strconv"

func Interpret(str string) float64 {
	ast := Parse(str)
	number := InterpretAST(ast)

	return number
}

func InterpretAST(ast AST) float64 {
	return calculateNode(ast.Node)
}

func calculateNode(node *Node) float64 {
	switch node.Type {
	case NInteger:
		integer, err := strconv.Atoi(node.Value)
		if err != nil {
			// @todo
		}

		return float64(integer)
	case NDecimal:
		decimal, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			// @todo
		}
		return decimal
	}

	if node.LeftChild == nil || node.RightChild == nil {
		return -1 // @todo correct error handling
	}

	c1 := calculateNode(node.LeftChild)
	c2 := calculateNode(node.RightChild)

	switch node.Type {
	case NAddition:
		return c1 + c2
	case NSubtraction:
		return c1 - c2
	case NMultiplication:
		return c1 * c2
	case NDivision:
		return c1 / c2
	}

	return 0 // @todo correct error handling
}
