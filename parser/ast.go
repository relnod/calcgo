package parser

// AST stores the data of the abstract syntax tree.
// The ast is in the form of a binary tree.
type AST struct {
	Node *Node
}

// NodeType defines the type of a node
type NodeType uint

// Node represents a node
type Node struct {
	Type       NodeType `json:"type"`
	Value      string   `json:"value"`
	LeftChild  *Node    `json:"left"`
	RightChild *Node    `json:"right"`
}

// IsOperator returns true if the given nodeType is an operator.
func (n *Node) IsOperator() bool {
	return n.Type > NDecimal && n.Type <= NDivision
}

func (n *Node) isHigherOperator(n2 *Node) bool {
	if n.Type <= NSubtraction && n2.Type <= NSubtraction {
		return false
	}

	if n.Type > NSubtraction && n2.Type > NSubtraction {
		return false
	}

	return n.Type < n2.Type
}
