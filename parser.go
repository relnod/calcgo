package calcgo

type AST struct {
	Node Node
}

type Node struct {
	Type NodeType `json:"type"`
	Value string `json:"value"`
	Childs []Node `json:"childs"`
}

type NodeType uint

const (
	NNumber NodeType = iota
	NAddition NodeType = iota
	NSubtraction NodeType = iota
	NMultiplication NodeType = iota
	NDivision NodeType = iota
)

func Parse(tokens []Token) AST {
	var ast AST
	var stack *Node

	ast = numberOrLeftBracketEntry(ast, tokens, 0, stack)

	return ast
}

func numberOrLeftBracketEntry(ast AST, tokens []Token, i int, stack *Node) AST {
	if i == len(tokens) {
		return ast
	}

	switch tokens[i].Type {
	case TNumber:
		node := Node{NNumber, tokens[i].Value, nil}
		ast.Node = node
		stack = &node

		i++
		return operator(ast, tokens, i, stack)
	}

	return ast
}

func numberOrLeftBracket(ast AST, tokens []Token, i int, stack *Node) AST {
	if i == len(tokens) {
		return ast
	}

	switch tokens[i].Type {
	case TNumber:
		node := Node{NNumber, tokens[i].Value, nil}
		ast.Node.Childs = append(ast.Node.Childs, node)
		stack = nil
	}

	i++
	return operator(ast, tokens, i, stack)
}

func operator(ast AST, tokens []Token, i int, stack *Node) AST {
	if i == len(tokens) {
		return ast
	}

	var node Node

	switch tokens[i].Type {
	case TOperatorPlus:
		node = Node{NAddition, "+", []Node{ast.Node}}
	case TOperatorMinus:
		node = Node{NSubtraction, "-", []Node{ast.Node}}
	case TOperatorMult:
		node = Node{NMultiplication, "*", []Node{ast.Node}}
	case TOperatorDiv:
		node = Node{NDivision, "/", []Node{ast.Node}}
	}

	ast.Node = node
	i++
	return numberOrLeftBracket(ast, tokens, i, stack)
}