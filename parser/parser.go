// Package parser contains the parser of calcgo.
package parser

import (
	"errors"

	"github.com/relnod/calcgo/lexer"
	"github.com/relnod/calcgo/token"
)

// parseState defines a state function of the parser machine.
type parseState func(*Parser) parseState

// Parser holds state of parser
type Parser struct {
	reader    token.Reader
	currToken token.Token
	topNode   *Node
	current   *Node
	errors    []error
	nested    bool
}

// Errors, that can occur during parsing
var (
	ErrorExpectedNumberOrVariable = errors.New("Error: Expected number or variable got something else")
	ErrorExpectedOperator         = errors.New("Error: Expected operator got something else")
	ErrorUnknownFunction          = errors.New("Error: Unknown Function")
	ErrorMissingClosingBracket    = errors.New("Error: Missing closing bracket")
	ErrorUnexpectedClosingBracket = errors.New("Error: Unexpected closing bracket")
)

// Parse parses a string to an ast
//
// Example:
//  calcgo.Parse("(1 + 2) * 3")
//
// Result:
//  calcgo.AST{
//    Node: &calcgo.Node{
//    Type:  calcgo.NMultiplication,
//    Value: "",
//    LeftChild: &calcgo.Node{
//      Type:  calcgo.NSubtraction,
//      Value: "",
//      LeftChild: &calcgo.Node{
//        Type: calcgo.NInteger,
//        Value:      "1",
//        LeftChild:  nil,
//        RightChild: nil,
//      },
//      RightChild: &calcgo.Node{
//        Type: calcgo.NInteger,
//        Value:      "2",
//        LeftChild:  nil,
//        RightChild: nil,
//      },
//    },
//    RightChild: &calcgo.Node{
//      Type:       calcgo.NInteger,
//      Value:      "3",
//      LeftChild:  nil,
//      RightChild: nil,
//    },
//  },
//
func Parse(str string) (AST, []error) {
	if len(str) == 0 {
		return AST{}, nil
	}

	l := lexer.NewBufferedLexerFromString(str)

	return ParseFromReader(l)
}

// ParseTokens parses a list of tokens to an ast
func ParseTokens(tokens []token.Token) (AST, []error) {
	if tokens == nil {
		return AST{}, nil
	}

	r := token.NewStaticReader(tokens)

	return ParseFromReader(r)
}

// ParseFromReader parses a stream of token retrieved from a token reader.
func ParseFromReader(reader token.Reader) (AST, []error) {
	p := &Parser{reader: reader}
	p.run()

	return AST{p.topNode}, p.errors
}

// run runs the parser state machine.
func (p *Parser) run() {
	for state := parseStart; state != nil; {
		ok := p.next()
		if !ok {
			break
		}

		state = state(p)
	}
}

// next retrieves the next token from the token. If the lexer is finished next
// returns false. Otherwise returns true.
func (p *Parser) next() bool {
	p.currToken = p.reader.Read()

	if p.currToken.Type == token.EOF {
		if p.nested {
			p.pushError(ErrorMissingClosingBracket)
		}
		return false
	}

	return true
}

// pushError adds an error to the parser error list.
func (p *Parser) pushError(err error) {
	p.errors = append(p.errors, err)
}

// pushErrors adds a list of errors to the parser error list.
func (p *Parser) pushErrors(errors []error) {
	for _, err := range errors {
		p.pushError(err)
	}
}

// newOperatorNode returns a new operator node.
func (p *Parser) newOperatorNode() *Node {
	nt, ok := getOperatorNodeType(p.currToken)
	if !ok {
		p.pushError(ErrorExpectedOperator)
	}

	return &Node{nt, p.currToken.Value, nil, nil}
}

// newNumberOrVariableNode returns a new number or variable node.
func (p *Parser) newNumberOrVariableNode() *Node {
	nt, ok := getNumberOrVariableNodeType(p.currToken)
	if !ok {
		p.pushError(ErrorExpectedNumberOrVariable)
	}

	return &Node{nt, p.currToken.Value, nil, nil}
}

// newFunctionNode returns a new function node.
func (p *Parser) newFunctionNode() *Node {
	nt, ok := getFunctionNodeType(p.currToken)
	if !ok {
		p.pushError(ErrorUnknownFunction)
	}

	return &Node{nt, p.currToken.Value, nil, nil}
}

// subParse creates another parser and runs it until a closing bracket appears.
func (p *Parser) subParse() (*Node, []error) {
	p2 := &Parser{
		reader:    p.reader,
		currToken: token.Token{},
		topNode:   p.topNode,
		current:   nil,
		errors:    nil,
		nested:    true,
	}

	p2.run()

	return p2.topNode, p2.errors
}

// subParseFunctionArgument parses the argument of a function.
func (p *Parser) subParseFunctionArgument(n *Node) {
	p.current = n
	n, errors := p.subParse()
	p.current.LeftChild = n
	p.pushErrors(errors)
}

// setFirstTopNode sets the first top node
func (p *Parser) setFirstTopNode(n *Node) {
	p.topNode = n
}

// setNewTopNode sets a new top node.
// The top node from before will be the left child of the new top node.
//
//    a            n
//   / \   =>     /
//  b   c        a
//              / \
//             b   c
//
func (p *Parser) setNewTopNode(n *Node) {
	n.LeftChild = p.topNode
	p.topNode = n
	p.current = n
}

// setNewRightChild sets a new right child of the current top child.
// The right child from before, will be the new left child of the new node.
//
//    a          a
//   / \   =>   / \
//  b   c      b   n
//                /
//               c
//
func (p *Parser) setNewRightChild(n *Node) {
	n.LeftChild = p.topNode.RightChild
	p.topNode.RightChild = n
	p.current = n
}

// addNewRightChild sets the right child to the given node.
//
//    a         a
//   /    =>   / \
//  b         b   n
//
func (p *Parser) addNewRightChild(n *Node) {
	p.current.RightChild = n
}

// parseStart is the start state of the parser machine.
// Behaves the same as parseValue, except, that it sets the first top node.
//
// Expects one of these tokens:
//  - TLeftBracket
//  - TFunc*
//  - TInteger
//  - TDecimal
//  - TVariable
//
// The following states can follow:
//  - parseOperator
//  - parseOperatorAfterRightBracket
//
func parseStart(p *Parser) parseState {
	if p.currToken.Type == token.ParenL {
		n, errors := p.subParse()
		p.setFirstTopNode(n)
		p.pushErrors(errors)

		return parseOperatorAfterRightBracket
	}

	if p.currToken.IsFunction() {
		n := p.newFunctionNode()
		p.setFirstTopNode(n)
		p.current = n
		p.subParseFunctionArgument(n)

		return parseOperator
	}

	n := p.newNumberOrVariableNode()
	p.setFirstTopNode(n)

	return parseOperator
}

// parseValue is the state, that parses
//
// Expects one of these tokens:
//  - TLeftBracket
//  - TFunc*
//  - TInteger
//  - TDecimal
//  - TVariable
//
// The following states can follow:
//  - parseOperator
//  - parseOperatorAfterRightBracket
//
func parseValue(p *Parser) parseState {
	if p.currToken.Type == token.ParenL {
		n, errors := p.subParse()
		p.addNewRightChild(n)
		p.pushErrors(errors)

		return parseOperator
	}

	if p.currToken.IsFunction() {
		n := p.newFunctionNode()
		p.addNewRightChild(n)
		p.subParseFunctionArgument(n)

		return parseOperator
	}

	n := p.newNumberOrVariableNode()
	p.addNewRightChild(n)
	return parseOperator
}

// parseOperator is the state, that handles an operator.
//
// Expects one of these tokens:
//  - TRightBracket
//  - TOperator*
//
// The following states can follow:
//  - parseValue
//
func parseOperator(p *Parser) parseState {
	if p.currToken.Type == token.ParenR {
		if !p.nested {
			p.pushError(ErrorUnexpectedClosingBracket)
		}
		return nil
	}

	node := p.newOperatorNode()
	// Handle 'multiplication and division before addition and subtraction' rule
	if IsOperator(p.topNode) && p.topNode.isHigherOperator(node) {
		p.setNewRightChild(node)
	} else {
		p.setNewTopNode(node)
	}

	return parseValue
}

// parseOperatorAfterRightBracket is the state, that happens after a closing
// bracket.
//
// Expects one of these tokens:
//  - TRightBracket
//  - TOperator*
//
// The following states can follow:
//  - parseValue
//
func parseOperatorAfterRightBracket(p *Parser) parseState {
	if p.currToken.Type == token.ParenR {
		if !p.nested {
			p.pushError(ErrorUnexpectedClosingBracket)
		}
		return nil
	}

	node := p.newOperatorNode()
	p.setNewTopNode(node)

	return parseValue
}
