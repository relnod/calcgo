package parser

import "github.com/relnod/calcgo/token"

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
	NMod
	NOr
	NXor
	NAnd
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

type CalcVisitor func(INode) (float64, error)

// INode defines an interface for a node.
type INode interface {
	Calculate(CalcVisitor) (float64, error)
	GetType() NodeType
	GetValue() string
	Left() INode
	Right() INode

	// IsLiteral() bool
	IsOperator() bool
	IsFunction() bool
}

// Calculatable defines an interface for a calculatable node.
type Calculatable interface {
	// Value returns the calculated result of a node
	// Calculate(CalcVisitor) float64
}

// Node represents a node
type Node struct {
	Type       NodeType `json:"type"`
	Value      string   `json:"value"`
	LeftChild  *Node    `json:"left"`
	RightChild *Node    `json:"right"`
}

func (n *Node) GetType() NodeType { return n.Type }
func (n *Node) GetValue() string  { return n.Value }
func (n *Node) Left() INode {
	if n.LeftChild == nil {
		return nil
	}

	return n.LeftChild
}
func (n *Node) Right() INode {
	if n.RightChild == nil {
		return nil
	}

	return n.RightChild
}

func (n *Node) Calculate(fn CalcVisitor) (float64, error) { return fn(n) }

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
func getOperatorNodeType(t token.Token) (NodeType, bool) {
	switch t.Type {
	case token.Plus:
		return NAdd, true
	case token.Minus:
		return NSub, true
	case token.Mult:
		return NMult, true
	case token.Div:
		return NDiv, true
	case token.Mod:
		return NMod, true
	case token.Or:
		return NOr, true
	case token.Xor:
		return NXor, true
	case token.And:
		return NAnd, true
	}

	return NInvalidOperator, false
}

// getOperatorNodeType converts a token type to a node type.
// The given token should be a number or variable. Returns an invalid number or
// variable node otherwise.
func getNumberOrVariableNodeType(t token.Token) (NodeType, bool) {
	switch t.Type {
	case token.Int:
		return NInt, true
	case token.Dec:
		return NDec, true
	case token.Bin:
		return NBin, true
	case token.Hex:
		return NHex, true
	case token.Exp:
		return NExp, true
	case token.Var:
		return NVar, true
	}

	switch t.Type {
	case token.InvalidCharacterInNumber:
		return NInvalidNumber, false
	case token.InvalidCharacterInVariable:
		return NInvalidVariable, false
	}

	return NError, false
}

// getOperatorNodeType converts a token type to a node type.
// The given token should be a function.
func getFunctionNodeType(t token.Token) (NodeType, bool) {
	switch t.Type {
	case token.Sqrt:
		return NFnSqrt, true
	case token.Sin:
		return NFnSin, true
	case token.Cos:
		return NFnCos, true
	case token.Tan:
		return NFnTan, true
	}

	return NInvalidFunction, false
}
