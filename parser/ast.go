package parser

import "github.com/relnod/calcgo/lexer"

// Node types
const (
	NError NodeType = iota
	NInvalidNumber
	NInvalidVariable
	NInvalidOperator
	NInvalidFunction
	NInteger
	NDecimal
	NVariable
	NAddition
	NFuncSqrt
	NFuncSin
	NFuncCos
	NFuncTan
	NSubtraction
	NMultiplication
	NDivision
)

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

// IsFunction returns true if the type of n is a function.
func (n *Node) IsFunction() bool {
	return n.Type == NFuncSqrt ||
		n.Type == NFuncSin ||
		n.Type == NFuncCos ||
		n.Type == NFuncTan
}

// IsOperator returns true if the type of n is an operator.
func (n *Node) IsOperator() bool {
	return n.Type == NAddition ||
		n.Type == NSubtraction ||
		n.Type == NMultiplication ||
		n.Type == NDivision
}

// isHigherOperator returns true if operator n is of higher than n2.
// Order is defined as the following:
// NAddition, NSubtraction = 0
// NMultiplication, NDivision = 1
func (n *Node) isHigherOperator(n2 *Node) bool {
	if n.Type <= NSubtraction && n2.Type <= NSubtraction {
		return false
	}

	if n.Type > NSubtraction && n2.Type > NSubtraction {
		return false
	}

	return n.Type < n2.Type
}

// getOperatorNodeType converts a token type to a node type.
// The given token should be an operator. Returns an invalid operator node
// otherwise.
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

// getOperatorNodeType converts a token type to a node type.
// The given token should be a number or variable. Returns an invalid number or
// variable node otherwise.
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

// getOperatorNodeType converts a token type to a node type.
// The given token should be a function.
func getFunctionNodeType(t lexer.Token) (NodeType, bool) {
	switch t.Type {
	case lexer.TFuncSqrt:
		return NFuncSqrt, true
	case lexer.TFuncSin:
		return NFuncSin, true
	case lexer.TFuncCos:
		return NFuncCos, true
	case lexer.TFuncTan:
		return NFuncTan, true
	}

	return NInvalidFunction, false
}
