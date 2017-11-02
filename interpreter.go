package calcgo

import "strconv"

func Interpret(ast AST) float64 {
	return calculateNode(ast.Node)
}

func calculateNode(node *Node) float64 {
	if (node.Type == NNumber) {
		integer, err := strconv.Atoi(node.Value)
		if (err != nil) {
			// @todo
		}

		return float64(integer)
	}

	c1 := calculateNode(node.Childs[0])
	c2 := calculateNode(node.Childs[1])

	switch(node.Type) {
	case NAddition:
		return c1 + c2
	case NSubtraction:
		return c1 - c2
	case NMultiplication:
		return c1 * c2
	case NDivision:
		return c1 / c2
	}
	
	return 0
}