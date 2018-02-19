package optimizer

import (
	"errors"

	"github.com/relnod/calcgo/interpreter/calculator"
	"github.com/relnod/calcgo/parser"
)

// Errors, that can occur during interpreting
var (
	ErrorMissingLeftChild       = errors.New("Error: Missing left child of node")
	ErrorMissingRightChild      = errors.New("Error: Missing right child of node")
	ErrorMissingFunctionArguent = errors.New("Error: Missing function argument")
	ErrorInvalidNodeType        = errors.New("Error: Invalid node type")
	ErrorInvalidVariable        = errors.New("Error: Invalid Variable")
	ErrorParserError            = errors.New("Error: Parser error")
	ErrorVariableNotDefined     = errors.New("Error: A variable was not defined")
)

// OptimizedAST holds an optimized ast.
// For all integer and decimal the value was already interpreted.
// All operations are already interpreted, if both child nodes could already be
// interpreted.
type OptimizedAST struct {
	Node *OptimizedNode
}

// OptimizedNode holds an optimized node
type OptimizedNode struct {
	Type        parser.NodeType
	Value       float64
	OldValue    string
	IsOptimized bool
	LeftChild   *OptimizedNode
	RightChild  *OptimizedNode
}

func (n *OptimizedNode) GetType() parser.NodeType { return n.Type }
func (n *OptimizedNode) GetValue() string         { return n.OldValue }
func (n *OptimizedNode) Left() parser.INode {
	if n.LeftChild == nil {
		return nil
	}

	return n.LeftChild
}
func (n *OptimizedNode) Right() parser.INode {
	if n.RightChild == nil {
		return nil
	}

	return n.RightChild
}

func (n *OptimizedNode) Calculate(fn parser.CalcVisitor) (float64, error) {
	if n.IsOptimized {
		return n.Value, nil
	}

	return fn(n)
}

// IsFunction returns true if the type of n is a function.
func (n *OptimizedNode) IsFunction() bool {
	return n.Type.IsFunction()
}

// IsOperator returns true if the type of n is an operator.
func (n *OptimizedNode) IsOperator() bool {
	return n.Type.IsOperator()
}

// newOptimizedNode returns a new optimized node.
func newOptimizedNode(value float64) *OptimizedNode {
	return &OptimizedNode{
		Type:        parser.NDec,
		Value:       value,
		OldValue:    "",
		IsOptimized: true,
		LeftChild:   nil,
		RightChild:  nil,
	}
}

// Optimize optimizes an ast.
// Interprets all integer and decimal nodes.
// Interprets all operations, if their child nodes can already be interpreted
func Optimize(ast *parser.AST) (*OptimizedAST, error) {
	if ast == nil {
		return nil, nil
	}

	optimizedNode, err := optimizeNode(ast.Node)
	if err != nil {
		return nil, err
	}

	return &OptimizedAST{Node: optimizedNode}, nil
}

// optimizeNode recursively optimizes all nodes, that can be optimized.
func optimizeNode(node *parser.Node) (*OptimizedNode, error) {
	if node.IsLiteral() {
		return optimizeLiteral(node)
	}
	if node.IsOperator() {
		return optimizeOperator(node)
	}

	if node.IsFunction() {
		return optimizeFunction(node)
	}

	return nil, ErrorInvalidNodeType
}

func optimizeLiteral(node *parser.Node) (*OptimizedNode, error) {
	if node.Type == parser.NVar {
		return &OptimizedNode{
			Type:        parser.NVar,
			Value:       0,
			OldValue:    node.Value,
			IsOptimized: false,
			LeftChild:   nil,
			RightChild:  nil,
		}, nil
	}

	var result float64
	var err error

	result, err = calculator.ConvertLiteral(node)
	if err != nil {
		return nil, err
	}

	return newOptimizedNode(result), nil
}

// optimizeOperator recursively optimizes an operator node and its child nodes.
func optimizeOperator(node *parser.Node) (*OptimizedNode, error) {
	left, right, err := getOptimizedNodeChilds(node)
	if err != nil {
		return nil, err
	}

	if !left.IsOptimized || !right.IsOptimized {
		return &OptimizedNode{
			Type:        node.Type,
			Value:       0,
			OldValue:    "",
			IsOptimized: false,
			LeftChild:   left,
			RightChild:  right,
		}, nil
	}

	var result float64
	result, err = calculator.CalculateOperator(left.Value, right.Value, node.Type)
	if err != nil {
		return nil, err
	}

	return newOptimizedNode(result), nil
}

// getOptimizedNodeChilds returns all optimized child nodes of a node.
func getOptimizedNodeChilds(node *parser.Node) (*OptimizedNode, *OptimizedNode, error) {
	if node.LeftChild == nil {
		return nil, nil, ErrorMissingLeftChild
	}
	if node.RightChild == nil {
		return nil, nil, ErrorMissingRightChild
	}

	left, err := optimizeNode(node.LeftChild)
	if err != nil {
		return nil, nil, err
	}
	right, err := optimizeNode(node.RightChild)
	if err != nil {
		return nil, nil, err
	}

	return left, right, nil
}

// optimizeFunction recursively optimizes a function node and its arguments.
func optimizeFunction(node *parser.Node) (*OptimizedNode, error) {
	left, err := optimizeNode(node.LeftChild)
	if err != nil {
		return nil, err
	}

	if !left.IsOptimized {
		return &OptimizedNode{
			Type:        node.Type,
			Value:       0,
			OldValue:    "",
			IsOptimized: false,
			LeftChild:   left,
			RightChild:  nil,
		}, nil
	}

	var result float64
	result, err = calculator.CalculateFunction(left.Value, node.Type)
	if err != nil {
		return nil, err
	}

	return newOptimizedNode(result), nil
}
