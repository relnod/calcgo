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
	Node parser.INode
}

// Root retruns root node.
func (a *OptimizedAST) Root() parser.INode {
	return a.Node
}

// Optimized returns true.
func (a *OptimizedAST) Optimized() bool {
	return true
}

// OptimizedNode holds an optimized node
type OptimizedNode struct {
	Type  parser.NodeType
	Value float64
}

// GetType return the type of the node.
func (n *OptimizedNode) GetType() parser.NodeType { return n.Type }

// GetValue return the value of the node.
func (n *OptimizedNode) GetValue() string { return "" }

// Left returns the left child.
func (n *OptimizedNode) Left() parser.INode { return nil }

// Right returns the right child.
func (n *OptimizedNode) Right() parser.INode { return nil }

// SetLeft sets the left child.
func (n *OptimizedNode) SetLeft(l parser.INode) {
	panic("") // @todo
}

// SetRight sets the right child.
func (n *OptimizedNode) SetRight(r parser.INode) {
	panic("") // @todo
}

// Calculate returns the calculated value if it is pre calculated.
// Otherwise returns the result of the calculation visitor.
func (n *OptimizedNode) Calculate(fn parser.CalcVisitor) (float64, error) {
	if n.GetType() == parser.NDec {
		return n.Value, nil
	}

	return fn(n)
}

// newOptimizedNode returns a new optimized node.
func newOptimizedNode(value float64) *OptimizedNode {
	return &OptimizedNode{
		Type:  parser.NDec,
		Value: value,
	}
}

// Optimize optimizes an ast.
// Interprets all integer and decimal nodes.
// Interprets all operations, if their child nodes can already be interpreted
func Optimize(ast parser.IAST) (*OptimizedAST, error) {
	if ast == nil {
		return nil, nil
	}

	optimizedNode, err := optimizeNode(ast.Root())
	if err != nil {
		return nil, err
	}

	return &OptimizedAST{Node: optimizedNode}, nil
}

// optimizeNode recursively optimizes all nodes, that can be optimized.
func optimizeNode(n parser.INode) (parser.INode, error) {
	if parser.IsLiteral(n) {
		return optimizeLiteral(n)
	}
	if parser.IsOperator(n) {
		return optimizeOperator(n)
	}

	if parser.IsFunction(n) {
		return optimizeFunction(n)
	}

	return nil, ErrorInvalidNodeType
}

func optimizeLiteral(n parser.INode) (parser.INode, error) {
	if n.GetType() == parser.NVar {
		return n, nil
	}

	var result float64
	var err error

	result, err = calculator.ConvertLiteral(n.GetValue(), n.GetType())
	if err != nil {
		return nil, err
	}

	return newOptimizedNode(result), nil
}

// optimizeOperator recursively optimizes an operator node and its child nodes.
func optimizeOperator(n parser.INode) (parser.INode, error) {
	left, right, err := getOptimizedNodeChilds(n)
	if err != nil {
		return nil, err
	}

	if left.GetType() != parser.NDec || right.GetType() != parser.NDec {
		n.SetLeft(left)
		n.SetRight(right)
		return n, nil
	}

	var result float64
	leftVal, _ := left.Calculate(nil)
	rightVal, _ := right.Calculate(nil)
	result, err = calculator.CalculateOperator(leftVal, rightVal, n.GetType())
	if err != nil {
		return nil, err
	}

	return newOptimizedNode(result), nil
}

// getOptimizedNodeChilds returns all optimized child nodes of a node.
func getOptimizedNodeChilds(n parser.INode) (parser.INode, parser.INode, error) {
	if n.Left() == nil {
		return nil, nil, ErrorMissingLeftChild
	}
	if n.Right() == nil {
		return nil, nil, ErrorMissingRightChild
	}

	left, err := optimizeNode(n.Left())
	if err != nil {
		return nil, nil, err
	}
	right, err := optimizeNode(n.Right())
	if err != nil {
		return nil, nil, err
	}

	return left, right, nil
}

// optimizeFunction recursively optimizes a function node and its arguments.
func optimizeFunction(n parser.INode) (parser.INode, error) {
	left, err := optimizeNode(n.Left())
	if err != nil {
		return nil, err
	}

	if left.GetType() != parser.NDec {
		n.SetLeft(left)
		return n, nil
	}

	var result float64
	val, _ := left.Calculate(nil)
	result, err = calculator.CalculateFunction(val, n.GetType())
	if err != nil {
		return nil, err
	}

	return newOptimizedNode(result), nil
}
