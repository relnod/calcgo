package calcgo

// AST stores the data of the abstract syntax tree.
// The ast is in the form of a binary tree.
type AST struct {
	Node *Node
}

// Node represents a node
type Node struct {
	Type       NodeType `json:"type"`
	Value      string   `json:"value"`
	LeftChild  *Node    `json:"left"`
	RightChild *Node    `json:"right"`
}

// NodeType defines the type of a node
type NodeType uint

// Node types
const (
	NError NodeType = iota
	NInteger
	NDecimal
	NAddition
	NSubtraction
	NMultiplication
	NDivision
)

// Parse parses a string to an ast
//
// Example:
//  calcgo.Parse("(1 + 2) * 3")
//
// Result:
//  calcgo.AST{
//    Node: &calcgo.Node{
//    Type:  calcgo.NMultiplication,
//    Value: "",
//    LeftChild: &calcgo.Node{
//      Type:  calcgo.NSubtraction,
//      Value: "",
//      LeftChild: &calcgo.Node{
//        Type: calcgo.NInteger,
//        Value:      "1",
//        LeftChild:  nil,
//        RightChild: nil,
//      },
//      RightChild: &calcgo.Node{
//        Type: calcgo.NInteger,
//        Value:      "2",
//        LeftChild:  nil,
//        RightChild: nil,
//      },
//    },
//    RightChild: &calcgo.Node{
//      Type:       calcgo.NInteger,
//      Value:      "3",
//      LeftChild:  nil,
//      RightChild: nil,
//    },
//  },
func Parse(str string) AST {
	tokens := Lex(str)
	return ParseTokens(tokens)
}

// ParseTokens parses a list of tokens to an ast
func ParseTokens(tokens []Token) AST {
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

// IsOperator returns true if the given nodeType is an operator.
func IsOperator(op NodeType) bool {
	return op > NDecimal && op <= NDivision
}

func isHigherOperator(op1 NodeType, op2 NodeType) bool {
	if op1 <= NSubtraction && op2 <= NSubtraction {
		return false
	}

	if op1 > NSubtraction && op2 > NSubtraction {
		return false
	}

	return op1 < op2
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

func numberOrLeftBracketEntry(topNode *Node, tokens []Token, i int, current *Node) (*Node, int) {
	i++
	if i == len(tokens) {
		return topNode, i
	}

	if tokens[i].Type == TLeftBracket {
		topNode, i = numberOrLeftBracketEntry(topNode, tokens, i, current)

		return operatorAfterRightBracket(topNode, tokens, i, current)
	}

	nodeType := getNumberNodeType(tokens[i].Type)

	//@todo: handle wrong node type

	node := &Node{nodeType, tokens[i].Value, nil, nil}
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
		current.RightChild = rightNode

		return operator(topNode, tokens, i, current)
	}

	nodeType := getNumberNodeType(tokens[i].Type)

	//@todo: handle wrong node type

	node := &Node{nodeType, tokens[i].Value, nil, nil}
	current.RightChild = node

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

	node := &Node{nodeType, tokens[i].Value, nil, nil}
	if IsOperator(topNode.Type) && isHigherOperator(topNode.Type, nodeType) {
		node.LeftChild = topNode.RightChild
		topNode.RightChild = node
	} else {
		node.LeftChild = topNode
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

	// @todo: handle wrong node type

	node := &Node{nodeType, tokens[i].Value, nil, nil}
	node.LeftChild = topNode
	topNode = node

	return numberOrLeftBracket(topNode, tokens, i, node)
}
