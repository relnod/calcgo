package calcgo

type AST struct {
	Node *Node
}

type Node struct {
	Type   NodeType `json:"type"`
	Value  string   `json:"value"`
	Childs []*Node  `json:"childs"`
}

type NodeType uint

const (
	NInteger        NodeType = iota
	NDecimal        NodeType = iota
	NAddition       NodeType = iota
	NSubtraction    NodeType = iota
	NMultiplication NodeType = iota
	NDivision       NodeType = iota
	NError          NodeType = iota
)

func isOperator(op NodeType) bool {
	return op > NDecimal
}

func isHigherOperator(op1 NodeType, op2 NodeType) bool {
	if op1 <= NSubtraction && op2 <= NSubtraction {
		return false
	}

	if op1 > NSubtraction && op2 > NSubtraction {
		return false
	}

	return true
}

func getNumberNodeType(token TokenType) NodeType {
	switch token {
	case TInteger:
		return NInteger
	case TDecimal:
		return NDecimal
	}

	return NError
}

func getOperatorNodeType(token TokenType) NodeType {
	switch token {
	case TOperatorPlus:
		return NAddition
	case TOperatorMinus:
		return NSubtraction
	case TOperatorMult:
		return NMultiplication
	case TOperatorDiv:
		return NDivision
	}

	return NError
}

func Parse(tokens []Token) AST {
	var topNode *Node
	var current *Node
	var i int

	topNode, i = numberOrLeftBracketEntry(topNode, tokens, -1, current)
	if i > 0 {
		// @todo
	}

	ast := AST{topNode}

	return ast
}

func numberOrLeftBracketEntry(topNode *Node, tokens []Token, i int, current *Node) (*Node, int) {
	i++
	if i == len(tokens) {
		return topNode, i
	}
	
	if tokens[i].Type == TLeftBracket {
		topNode, i := numberOrLeftBracketEntry(topNode, tokens, i, current)

		return operatorAfterRightBracket(topNode, tokens, i, current)
	}

	nodeType := getNumberNodeType(tokens[i].Type)

	//@todo: handle wrong node type		

	node := &Node{nodeType, tokens[i].Value, nil}
	topNode = node

	return operator(topNode, tokens, i, current)
}

func numberOrLeftBracket(topNode *Node, tokens []Token, i int, current *Node) (*Node, int) {
	i++
	if i == len(tokens) {
		return topNode, i
	}

	if tokens[i].Type == TLeftBracket {
		topNodeNested := topNode
		rightNode, i := numberOrLeftBracketEntry(topNodeNested, tokens, i, current)
		topNode.Childs = append(topNode.Childs, rightNode)

		return operator(topNode, tokens, i, current)
	}
	
	nodeType := getNumberNodeType(tokens[i].Type)

	//@todo: handle wrong node type		

	node := &Node{nodeType, tokens[i].Value, nil}
	current.Childs = append(current.Childs, node)

	return operator(topNode, tokens, i, current)
}

func operator(topNode *Node, tokens []Token, i int, current *Node) (*Node, int) {
	i++
	if i == len(tokens) {
		return topNode, i
	}

	if tokens[i].Type == TRightBracket {
		return topNode, i
	}

	nodeType := getOperatorNodeType(tokens[i].Type)

	// @todo: handle wrong node type

	node := &Node{nodeType, tokens[i].Value, nil}
	if isOperator(topNode.Type) && isHigherOperator(topNode.Type, nodeType) {
		node.Childs = []*Node{topNode.Childs[1]}
		topNode.Childs[1] = node
	} else {
		node.Childs = []*Node{topNode}
		topNode = node
	}

	return numberOrLeftBracket(topNode, tokens, i, node)
}

func operatorAfterRightBracket(topNode *Node, tokens []Token, i int, current *Node) (*Node, int) {
	i++
	if i == len(tokens) {
		return topNode, i
	}

	if tokens[i].Type == TRightBracket {
		return topNode, i
	}

	nodeType := getOperatorNodeType(tokens[i].Type)

	// @todo: hanndle wrong node type

	node := &Node{nodeType, tokens[i].Value, nil}
	node.Childs = []*Node{topNode}
	topNode = node

	return numberOrLeftBracket(topNode, tokens, i, node)
}