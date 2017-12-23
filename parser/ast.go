package parser

import "github.com/relnod/calcgo/lexer"

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

func getOperatorNodeType(t lexer.Token) (NodeType, bool) {
	switch t.Type {
	case lexer.TOperatorPlus:
		return NAddition, true
	case lexer.TOperatorMinus:
		return NSubtraction, true
	case lexer.TOperatorMult:
		return NMultiplication, true
	case lexer.TOperatorDiv:
		return NDivision, true
	}

	return NInvalidOperator, false
}

func getNumberOrVariableNodeType(t lexer.Token) (NodeType, bool) {
	switch t.Type {
	case lexer.TInteger:
		return NInteger, true
	case lexer.TDecimal:
		return NDecimal, true
	case lexer.TVariable:
		return NVariable, true
	}

	switch t.Type {
	case lexer.TInvalidCharacterInNumber:
		return NInvalidNumber, false
	case lexer.TInvalidCharacterInVariable:
		return NInvalidVariable, false
	}

	return NError, false
}

func getFunctionNodeType(t lexer.Token) (NodeType, bool) {
	return NFuncSqrt, true
}
