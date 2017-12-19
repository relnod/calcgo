package interpreter

import (
	"errors"
	"math"
	"strconv"

	"github.com/relnod/calcgo/parser"
)

// Errors, that can be returned by the interpreter
var (
	ErrorMissingLeftChild       = errors.New("Error: Missing left child of node")
	ErrorMissingRightChild      = errors.New("Error: Missing right child of node")
	ErrorMissingFunctionArguent = errors.New("Error: Missing function argument")
	ErrorInvalidNodeType        = errors.New("Error: Invalid node type")
	ErrorInvalidInteger         = errors.New("Error: Invalid Integer")
	ErrorInvalidDecimal         = errors.New("Error: Invalid Decimal")
	ErrorInvalidVariable        = errors.New("Error: Invalid Variable")
	ErrorParserError            = errors.New("Error: Parser error")
	ErrorDivisionByZero         = errors.New("Error: Division by zero")
	ErrorVariableNotDefined     = errors.New("Error: A variable was not defined")
)

// Interpreter holds state of interpreter
type Interpreter struct {
	str              string
	ast              *parser.AST
	oast             *OptimizedAST
	vars             map[string]float64
	optimizerEnabled bool
}

// NewInterpreter returns a new interpreter from a string
func NewInterpreter(str string) *Interpreter {
	return &Interpreter{
		str:              str,
		ast:              nil,
		oast:             nil,
		vars:             make(map[string]float64),
		optimizerEnabled: false,
	}
}

// NewInterpreterFromAST returns a new interpreter from an ast
func NewInterpreterFromAST(ast *parser.AST) *Interpreter {
	return &Interpreter{
		str:              "",
		ast:              ast,
		oast:             nil,
		vars:             make(map[string]float64),
		optimizerEnabled: false,
	}
}

// SetVar sets the value of a variable
func (i *Interpreter) SetVar(name string, value float64) {
	i.vars[name] = value
}

// EnableOptimizer enables optimization of the ast.
// Optimization happens at the next GetResult() call
func (i *Interpreter) EnableOptimizer() {
	i.optimizerEnabled = true
}

// GetResult interprets the ast.
// All variables have to be set up to this point
//
// If the interpreter was initialized with a string,
// the ast gets generated first
func (i *Interpreter) GetResult() (float64, []error) {
	if i.str == "" {
		return 0, nil
	}

	if i.ast == nil {
		ast, errors := parser.Parse(i.str)
		if errors != nil {
			return 0, errors
		}

		i.ast = &ast
	}

	var result float64
	var err error
	if i.optimizerEnabled {
		if i.oast == nil {
			oast, err := Optimize(i.ast)
			if err != nil {
				return 0, []error{err}
			}

			i.oast = oast
		}

		result, err = i.interpretOptimizedNode(i.oast.Node)
	} else {
		result, err = i.interpretNode(i.ast.Node)
	}

	if err != nil {
		return 0, []error{err}
	}

	return result, nil
}

// Interpret interprets a given string.
// Returns errors if lexing, parsing or interpreting failed
//
// Examples:
//  caclgo.Interpret("(1 + 2) * 3") // Result: 9
//  caclgo.Interpret("1 + 2 * 3")   // Result: 7
func Interpret(str string) (float64, []error) {
	if len(str) == 0 {
		return 0, nil
	}

	interpreter := NewInterpreter(str)

	return interpreter.GetResult()
}

// InterpretAST interprets a given ast.
// Returns an error if the ast is invalid.
//
// Example:
//  result, err := calcgo.InterpretAST(calcgo.AST{
//		Node: &calcgo.Node{
//			Type:  calcgo.NAddition,
//			Value: "",
//			LeftChild: &calcgo.Node{
//				Type:       calcgo.NInteger,
//				Value:      "1",
//				LeftChild:  nil,
//				RightChild: nil,
//			},
//			RightChild: &calcgo.Node{
//				Type:       calcgo.NInteger,
//				Value:      "2",
//				LeftChild:  nil,
//				RightChild: nil,
//			},
//		},
//  })
func InterpretAST(ast parser.AST) (float64, error) {
	interpreter := NewInterpreterFromAST(&ast)

	result, errors := interpreter.GetResult()
	if errors != nil {
		return 0, errors[0]
	}

	return result, nil
}

func (i *Interpreter) interpretNode(node *parser.Node) (float64, error) {
	switch node.Type {
	case parser.NInteger:
		return interpretInteger(node)
	case parser.NDecimal:
		return interpretDecimal(node)
	case parser.NVariable:
		return i.interpretVariable(node)
	case parser.NAddition, parser.NSubtraction, parser.NMultiplication, parser.NDivision:
		return i.interpretOperator(node)
	case parser.NFuncSqrt:
		return i.interpretSqrt(node)
	}

	return 0, ErrorInvalidNodeType
}

func (i *Interpreter) interpretOptimizedNode(node *OptimizedNode) (float64, error) {
	if node.IsOptimized {
		return node.Value, nil
	}

	switch node.Type {
	case parser.NVariable:
		return i.interpretOptimizedVariable(node)
	case parser.NFuncSqrt:
		return i.interpretOptimizedSqrt(node)
	}

	return i.interpretOptimizedOperator(node)
}

func interpretInteger(node *parser.Node) (float64, error) {
	integer, err := strconv.Atoi(node.Value)
	if err != nil {
		return 0, ErrorInvalidInteger
	}
	return float64(integer), nil
}

func interpretDecimal(node *parser.Node) (float64, error) {
	decimal, err := strconv.ParseFloat(node.Value, 64)
	if err != nil {
		return 0, ErrorInvalidDecimal
	}
	return decimal, nil
}

func (i *Interpreter) interpretVariable(node *parser.Node) (float64, error) {
	number, ok := i.vars[node.Value]
	if ok {
		return number, nil
	}

	return 0, ErrorVariableNotDefined
}

func (i *Interpreter) interpretOptimizedVariable(node *OptimizedNode) (float64, error) {
	number, ok := i.vars[node.OldValue]
	if ok {
		return number, nil
	}

	return 0, ErrorVariableNotDefined
}

func (i *Interpreter) interpretOperator(node *parser.Node) (float64, error) {
	left, right, err := i.getInterpretedNodeChilds(node)
	if err != nil {
		return 0, err
	}

	return calculateOperator(left, right, node.Type)
}

func (i *Interpreter) interpretOptimizedOperator(node *OptimizedNode) (float64, error) {
	left, right, err := i.getInterpretedOptimizedNodeChilds(node)
	if err != nil {
		return 0, err
	}

	return calculateOperator(left, right, node.Type)
}

func (i *Interpreter) getInterpretedNodeChilds(node *parser.Node) (float64, float64, error) {
	if node.LeftChild == nil {
		return 0, 0, ErrorMissingLeftChild
	}
	if node.RightChild == nil {
		return 0, 0, ErrorMissingRightChild
	}

	left, err := i.interpretNode(node.LeftChild)
	if err != nil {
		return 0, 0, err
	}
	right, err := i.interpretNode(node.RightChild)
	if err != nil {
		return 0, 0, err
	}

	return left, right, nil
}

func (i *Interpreter) getInterpretedOptimizedNodeChilds(node *OptimizedNode) (float64, float64, error) {
	left, err := i.interpretOptimizedNode(node.LeftChild)
	if err != nil {
		return 0, 0, err
	}
	right, err := i.interpretOptimizedNode(node.RightChild)
	if err != nil {
		return 0, 0, err
	}

	return left, right, nil
}

func (i *Interpreter) interpretSqrt(node *parser.Node) (float64, error) {
	left, err := i.interpretNode(node.LeftChild)
	if err != nil {
		return 0, err
	}

	return math.Sqrt(left), nil
}

func (i *Interpreter) interpretOptimizedSqrt(node *OptimizedNode) (float64, error) {
	left, err := i.interpretOptimizedNode(node.LeftChild)
	if err != nil {
		return 0, err
	}

	return math.Sqrt(left), nil
}
