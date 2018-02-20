package interpreter

import (
	"errors"

	"github.com/relnod/calcgo/interpreter/calculator"
	"github.com/relnod/calcgo/interpreter/optimizer"
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

// Interpreter holds state of interpreter
type Interpreter struct {
	str              string
	ast              parser.IAST
	vars             map[string]float64
	optimizerEnabled bool
}

// NewInterpreter returns a new interpreter from a string
func NewInterpreter(str string) *Interpreter {
	return &Interpreter{
		str:              str,
		ast:              nil,
		vars:             make(map[string]float64),
		optimizerEnabled: false,
	}
}

// NewInterpreterFromAST returns a new interpreter from an ast
func NewInterpreterFromAST(ast parser.IAST) *Interpreter {
	return &Interpreter{
		str:              "",
		ast:              ast,
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
	if i.str == "" && i.ast == nil {
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
	if i.optimizerEnabled && !i.ast.Optimized() {
		oast, err := optimizer.Optimize(i.ast)
		if err != nil {
			return 0, []error{err}
		}

		i.ast = oast
	}

	result, err = i.ast.Root().Calculate(i.calcVisitor)

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
//
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
//
func InterpretAST(ast parser.IAST) (float64, error) {
	i := NewInterpreterFromAST(ast)

	result, errors := i.GetResult()
	if errors != nil {
		return 0, errors[0]
	}

	return result, nil
}

func (i *Interpreter) calcVisitor(n parser.INode) (float64, error) {
	switch n.GetType() {
	case parser.NVar:
		return i.interpretVariable(n)
	case parser.NInt:
		return calculator.ConvertInteger(n.GetValue())
	case parser.NDec:
		return calculator.ConvertDecimal(n.GetValue())
	case parser.NBin:
		return calculator.ConvertBin(n.GetValue())
	case parser.NHex:
		return calculator.ConvertHex(n.GetValue())
	case parser.NExp:
		return calculator.ConvertExponential(n.GetValue())
	}

	if parser.IsOperator(n) {
		return i.interpretOperator(n)
	}

	if parser.IsFunction(n) {
		return i.interpretFunction(n)
	}

	return 0, ErrorInvalidNodeType
}

// interpretVariable interprets a variable node.
// Returns an error if the variable is not defined.
func (i *Interpreter) interpretVariable(n parser.INode) (float64, error) {
	number, ok := i.vars[n.GetValue()]
	if ok {
		return number, nil
	}

	return 0, ErrorVariableNotDefined
}

// interpretOperator recursively interprets an operator node.
func (i *Interpreter) interpretOperator(n parser.INode) (float64, error) {
	left, right, err := i.getInterpretedNodeChilds(n)
	if err != nil {
		return 0, err
	}

	return calculator.CalculateOperator(left, right, n.GetType())
}

// interpretFunction interprets a function node
func (i *Interpreter) interpretFunction(n parser.INode) (float64, error) {
	left, err := n.Left().Calculate(i.calcVisitor)
	if err != nil {
		return 0, err
	}

	return calculator.CalculateFunction(left, n.GetType())
}

// getInterpretedNodeChilds returns the interpreted child nodes of a given node.
// Both child nodes have to be defined. Retruns an error otherwise.
func (i *Interpreter) getInterpretedNodeChilds(n parser.INode) (float64, float64, error) {
	if n.Left() == nil {
		return 0, 0, ErrorMissingLeftChild
	}
	if n.Right() == nil {
		return 0, 0, ErrorMissingRightChild
	}

	left, err := n.Left().Calculate(i.calcVisitor)
	if err != nil {
		return 0, 0, err
	}
	right, err := n.Right().Calculate(i.calcVisitor)
	if err != nil {
		return 0, 0, err
	}

	return left, right, nil
}
