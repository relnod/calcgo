package calcgo

type AST struct {
	Node *Node
}

type Node struct {
	Type NodeType `json:"type"`
	Value string `json:"value"`
	Childs []*Node `json:"childs"`
}

type NodeType uint

const (
	NNumber NodeType = iota
	NAddition NodeType = iota
	NSubtraction NodeType = iota
	NMultiplication NodeType = iota
	NDivision NodeType = iota
)

func isOperator(op NodeType) bool {
	return op > NNumber
}
func isHigherOperator(op1 NodeType, op2 NodeType) bool {
	return op1 < op2
}

func Parse(tokens []Token) AST {
	var ast AST
	var current *Node

	ast = numberOrLeftBracketEntry(ast, tokens, 0, current)

	return ast
}

func numberOrLeftBracketEntry(ast AST, tokens []Token, i int, current *Node) AST {
	if i == len(tokens) {
		return ast
	}

	switch tokens[i].Type {
	case TNumber:
		node := &Node{NNumber, tokens[i].Value, nil}
		ast.Node = node	
	}

	i++
	return operator(ast, tokens, i, current)
}

func numberOrLeftBracket(ast AST, tokens []Token, i int, current *Node) AST {
	if i == len(tokens) {
		return ast
	}

	switch tokens[i].Type {
	case TNumber:
		node := &Node{NNumber, tokens[i].Value, nil}
		current.Childs = append(current.Childs, node)
	}

	i++
	return operator(ast, tokens, i, current)
}

func operator(ast AST, tokens []Token, i int, current *Node) AST {
	if i == len(tokens) {
		return ast
	}

	var nodeType NodeType

	switch tokens[i].Type {
	case TOperatorPlus:
		nodeType = NAddition
	case TOperatorMinus:
		nodeType = NSubtraction
	case TOperatorMult:
		nodeType = NMultiplication
	case TOperatorDiv:
		nodeType = NDivision
	}

	node := &Node{nodeType, tokens[i].Value, nil}
	if isOperator(ast.Node.Type) && isHigherOperator(ast.Node.Type, nodeType) {
		node.Childs = []*Node{ast.Node.Childs[1]}
		ast.Node.Childs[1] = node
	} else {
		node.Childs = []*Node{ast.Node}
		ast.Node = node
	}

	current = node
	i++
	return numberOrLeftBracket(ast, tokens, i, current)
}