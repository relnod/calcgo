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

	literalBeg
	// Numbers
	NInt
	NDec
	NBin
	NHex
	NExp

	// Variable
	NVar
	literalEnd

	operatorBeg
	// Operators
	NAdd
	NSub
	NMult
	NDiv
	operatorEnd

	functionBeg
	// Functions
	NFnSqrt
	NFnSin
	NFnCos
	NFnTan
	functionEnd
)

// AST stores the data of the abstract syntax tree.
// The ast is in the form of a binary tree.
type AST struct {
	Node *Node
}

// NodeType defines the type of a node
type NodeType uint

// IsLiteral returns true if t is a literal.
func (t NodeType) IsLiteral() bool {
	return literalBeg < t && t < literalEnd
}

// IsOperator returns true if t is an operator.
func (t NodeType) IsOperator() bool {
	return operatorBeg < t && t < operatorEnd
}

// IsFunction returns true if t is a function.
func (t NodeType) IsFunction() bool {
	return functionBeg < t && t < functionEnd
}

// Node represents a node
type Node struct {
	Type       NodeType `json:"type"`
	Value      string   `json:"value"`
	LeftChild  *Node    `json:"left"`
	RightChild *Node    `json:"right"`
}

// IsLiteral returns true if n is literal node.
func (n *Node) IsLiteral() bool {
	return n.Type.IsLiteral()
}

// IsOperator returns true if n is an operator node.
func (n *Node) IsOperator() bool {
	return n.Type.IsOperator()
}

// IsFunction returns true if n is a function node.
func (n *Node) IsFunction() bool {
	return n.Type.IsFunction()
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
