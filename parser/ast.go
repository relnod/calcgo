package parser

import "github.com/relnod/calcgo/token"

// NodeType defines the type of a node
type NodeType uint

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

// CalcVisitor defines the visitor function called when calculation a node.
type CalcVisitor func(INode) (float64, error)

// INode defines an interface for a node.
type INode interface {
	GetType() NodeType
	GetValue() string
	Left() INode
	Right() INode

	SetLeft(INode)
	SetRight(INode)

	Calculate(CalcVisitor) (float64, error)
}

// IsLiteral returns true if t is a literal.
func IsLiteral(n INode) bool {
	return literalBeg < n.GetType() && n.GetType() < literalEnd
}

// IsOperator returns true if t is an operator.
func IsOperator(n INode) bool {
	return operatorBeg < n.GetType() && n.GetType() < operatorEnd
}

// IsFunction returns true if t is a function.
func IsFunction(n INode) bool {
	return functionBeg < n.GetType() && n.GetType() < functionEnd
}

// IAST defines an interface for an ast.
type IAST interface {
	// Root returns the root node.
	Root() INode

	// Optimized returns true if the ast is optimized.
	Optimized() bool
}

// Node represents a node
type Node struct {
	Type       NodeType
	Value      string
	LeftChild  INode
	RightChild INode
}

// GetType returns the type of the node.
func (n *Node) GetType() NodeType { return n.Type }

// GetValue returns the value of the node.
func (n *Node) GetValue() string { return n.Value }

// Left returns the left child.
func (n *Node) Left() INode {
	if n.LeftChild == nil {
		return nil
	}

	return n.LeftChild
}

// Right returns the left child.
func (n *Node) Right() INode {
	if n.RightChild == nil {
		return nil
	}

	return n.RightChild
}

func (n *Node) SetLeft(l INode) {
	n.LeftChild = l
}

func (n *Node) SetRight(r INode) {
	n.RightChild = r
}

// Calculate returns the result of the calculation visitor.
func (n *Node) Calculate(fn CalcVisitor) (float64, error) { return fn(n) }

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

// AST stores the data of the abstract syntax tree.
// The ast is in the form of a binary tree.
type AST struct {
	Node *Node
}

// Root returns the root node.
func (a *AST) Root() INode {
	return a.Node
}

// Optimized returns false.
func (a *AST) Optimized() bool {
	return false
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
