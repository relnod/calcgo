package parser

import "github.com/relnod/calcgo/lexer"

// Node types
const (
	// Errors
	NError NodeType = iota
	NInvalidNumber
	NInvalidVariable
	NInvalidOperator
	NInvalidFunction

	// Numbers
	NInt
	NDec
	NBin
	NHex
	NExp

	// Variable
	NVar

	// Operators
	NAdd
	NSub
	NMult
	NDiv

	// Functions
	NFnSqrt
	NFnSin
	NFnCos
	NFnTan
)

// AST stores the data of the abstract syntax tree.
// The ast is in the form of a binary tree.
type AST struct {
	Node *Node
}

// NodeType defines the type of a node
type NodeType uint

// IsFunction returns true if the type of n is a function.
func (n NodeType) IsFunction() bool {
	return n == NFnSqrt ||
		n == NFnSin ||
		n == NFnCos ||
		n == NFnTan
}

// IsOperator returns true if the type of n is an operator.
func (n NodeType) IsOperator() bool {
	return n == NAdd ||
		n == NSub ||
		n == NMult ||
		n == NDiv
}

// Node represents a node
type Node struct {
	Type       NodeType `json:"type"`
	Value      string   `json:"value"`
	LeftChild  *Node    `json:"left"`
	RightChild *Node    `json:"right"`
}

// IsFunction returns true if the type of n is a function.
func (n *Node) IsFunction() bool {
	return n.Type.IsFunction()
}

// IsOperator returns true if the type of n is an operator.
func (n *Node) IsOperator() bool {
	return n.Type.IsOperator()
}

// isHigherOperator returns true if operator n is of higher than n2.
// Order is defined as the following:
// NAddition, NSubtraction = 0
// NMultiplication, NDivision = 1
func (n *Node) isHigherOperator(n2 *Node) bool {
	if n.Type <= NSub && n2.Type <= NSub {
		return false
	}

	if n.Type > NSub && n2.Type > NSub {
		return false
	}

	return n.Type < n2.Type
}

// getOperatorNodeType converts a token type to a node type.
// The given token should be an operator. Returns an invalid operator node
// otherwise.
func getOperatorNodeType(t lexer.Token) (NodeType, bool) {
	switch t.Type {
	case lexer.TOpPlus:
		return NAdd, true
	case lexer.TOpMinus:
		return NSub, true
	case lexer.TOpMult:
		return NMult, true
	case lexer.TOpDiv:
		return NDiv, true
	}

	return NInvalidOperator, false
}

// getOperatorNodeType converts a token type to a node type.
// The given token should be a number or variable. Returns an invalid number or
// variable node otherwise.
func getNumberOrVariableNodeType(t lexer.Token) (NodeType, bool) {
	switch t.Type {
	case lexer.TInt:
		return NInt, true
	case lexer.TDec:
		return NDec, true
	case lexer.TBin:
		return NBin, true
	case lexer.THex:
		return NHex, true
	case lexer.TExp:
		return NExp, true
	case lexer.TVar:
		return NVar, true
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
	case lexer.TFnSqrt:
		return NFnSqrt, true
	case lexer.TFnSin:
		return NFnSin, true
	case lexer.TFnCos:
		return NFnCos, true
	case lexer.TFnTan:
		return NFnTan, true
	}

	return NInvalidFunction, false
}
