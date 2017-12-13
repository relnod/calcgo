package calcgo

import (
	"errors"
	"strconv"
)

// Errors, that can be returned by the interpreter
var (
	ErrorMissingLeftChild   = errors.New("Error: Missing left child of node")
	ErrorMissingRightChild  = errors.New("Error: Missing right child of node")
	ErrorInvalidNodeType    = errors.New("Error: Invalid node type")
	ErrorInvalidInteger     = errors.New("Error: Invalid Integer")
	ErrorInvalidDecimal     = errors.New("Error: Invalid Decimal")
	ErrorInvalidVariable    = errors.New("Error: Invalid Variable")
	ErrorParserError        = errors.New("Error: Parser error")
	ErrorDivisionByZero     = errors.New("Error: Division by zero")
	ErrorVariableNotDefined = errors.New("Error: A variable was not defined")
)

// Interpreter holds state of interpreter
type Interpreter struct {
	str  string
	ast  *AST
	vars map[string]float64
}

// NewInterpreter returns a new interpreter from a string
func NewInterpreter(str string) *Interpreter {
	return &Interpreter{
		str:  str,
		ast:  nil,
		vars: make(map[string]float64),
	}
}

// NewInterpreterFromAST returns a new interpreter from an ast
func NewInterpreterFromAST(ast *AST) *Interpreter {
	return &Interpreter{
		str:  "",
		ast:  ast,
		vars: make(map[string]float64),
	}
}

// SetVar sets the value of a variable
func (i *Interpreter) SetVar(name string, value float64) {
	i.vars[name] = value
}

// GetResult interprets the ast.
// All variables have to be set up to this point
//
// If the interpreter was initialized with a string, the ast gets generated first
func (i *Interpreter) GetResult() (float64, []error) {
	if i.ast == nil {
		ast, errors := Parse(i.str)
		if errors != nil {
			return 0, errors
		}

		i.ast = &ast
	}

	result, err := i.interpretNode(i.ast.Node)
	if err != nil {
		return 0, []error{err}
	}

	return result, nil
}

// Interpret interprets a given string.
// Can return an error if parsing failed
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
// Can return an error if the ast is invalid.
func InterpretAST(ast AST) (float64, error) {
	interpreter := NewInterpreterFromAST(&ast)

	result, errors := interpreter.GetResult()
	if errors != nil {
		return 0, errors[0]
	}

	return result, nil
}

func (i *Interpreter) interpretNode(node *Node) (float64, error) {
	switch node.Type {
	case NInteger:
		return i.interpretInteger(node)
	case NDecimal:
		return i.interpretDecimal(node)
	case NVariable:
		return i.interpretVariable(node)
	case NAddition:
		return i.interpretAddition(node)
	case NSubtraction:
		return i.interpretSubtraction(node)
	case NMultiplication:
		return i.interpretMultiplication(node)
	case NDivision:
		return i.interpretDivision(node)
	}

	return 0, ErrorInvalidNodeType
}

func (i *Interpreter) interpretInteger(node *Node) (float64, error) {
	integer, err := strconv.Atoi(node.Value)
	if err != nil {
		return 0, ErrorInvalidInteger
	}
	return float64(integer), nil
}

func (i *Interpreter) interpretDecimal(node *Node) (float64, error) {
	decimal, err := strconv.ParseFloat(node.Value, 64)
	if err != nil {
		return 0, ErrorInvalidDecimal
	}
	return decimal, nil
}

func (i *Interpreter) interpretVariable(node *Node) (float64, error) {
	number, ok := i.vars[node.Value]
	if ok {
		return number, nil
	}

	return 0, ErrorVariableNotDefined
}

func (i *Interpreter) interpretAddition(node *Node) (float64, error) {
	left, right, err := i.getInterpretedNodeChilds(node)
	if err != nil {
		return 0, err
	}

	return left + right, nil
}

func (i *Interpreter) interpretSubtraction(node *Node) (float64, error) {
	left, right, err := i.getInterpretedNodeChilds(node)
	if err != nil {
		return 0, err
	}

	return left - right, nil
}

func (i *Interpreter) interpretMultiplication(node *Node) (float64, error) {
	left, right, err := i.getInterpretedNodeChilds(node)
	if err != nil {
		return 0, err
	}

	return left * right, nil
}

func (i *Interpreter) interpretDivision(node *Node) (float64, error) {
	left, right, err := i.getInterpretedNodeChilds(node)
	if err != nil {
		return 0, err
	}

	if right == 0 {
		return 0, ErrorDivisionByZero
	}

	return left / right, nil
}

func (i *Interpreter) getInterpretedNodeChilds(node *Node) (float64, float64, error) {
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
