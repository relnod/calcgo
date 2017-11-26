package calcgo

import "errors"

type parseState func(*Parser) parseState

// AST stores the data of the abstract syntax tree.
// The ast is in the form of a binary tree.
type AST struct {
	Node *Node
}

// Node represents a node
type Node struct {
	Type       NodeType `json:"type"`
	Value      string   `json:"value"`
	LeftChild  *Node    `json:"left"`
	RightChild *Node    `json:"right"`
}

// NodeType defines the type of a node
type NodeType uint

// Parser holds state of parser
type Parser struct {
	tokens  []Token
	topNode *Node
	current *Node
	index   int
	errors  []error
	nested  bool
}

// Node types
const (
	NError NodeType = iota
	NInvalidNumber
	NInvalidOperator
	NInteger
	NDecimal
	NAddition
	NSubtraction
	NMultiplication
	NDivision
)

// Errors retured by the parser
var (
	ErrorExpectedNumber           = errors.New("Error: Expected number got something else")
	ErrorExpectedOperator         = errors.New("Error: Expected operator got something else")
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
func Parse(str string) (AST, []error) {
	tokens := Lex(str)
	return ParseTokens(tokens)
}

// ParseTokens parses a list of tokens to an ast
func ParseTokens(tokens []Token) (AST, []error) {
	if tokens == nil {
		return AST{}, nil
	}

	p := &Parser{tokens, nil, nil, -1, nil, false}

	topNode, _ := p.run()

	ast := AST{topNode}

	return ast, p.errors
}

// IsOperator returns true if the given nodeType is an operator.
func IsOperator(op NodeType) bool {
	return op > NDecimal && op <= NDivision
}

func isHigherOperator(op1 NodeType, op2 NodeType) bool {
	if op1 <= NSubtraction && op2 <= NSubtraction {
		return false
	}

	if op1 > NSubtraction && op2 > NSubtraction {
		return false
	}

	return op1 < op2
}

func (p *Parser) run() (*Node, int) {
	for state := parseStart; state != nil; {
		p.index++
		if p.index >= len(p.tokens) {
			if p.nested {
				p.pushError(ErrorMissingClosingBracket)
			}
			break
		}

		state = state(p)
	}

	return p.topNode, p.index
}

func (p *Parser) currentType() TokenType {
	return p.tokens[p.index].Type
}

func (p *Parser) currentValue() string {
	return p.tokens[p.index].Value
}

func (p *Parser) newNode(nodeType NodeType) *Node {
	return &Node{nodeType, p.currentValue(), nil, nil}
}

func (p *Parser) pushError(err error) {
	p.errors = append(p.errors, err)
}

func (p *Parser) pushErrors(errors []error) {
	for _, err := range errors {
		p.pushError(err)
	}
}

func (p *Parser) getNumberNodeType() NodeType {
	switch p.currentType() {
	case TInteger:
		return NInteger
	case TDecimal:
		return NDecimal
	}

	p.pushError(ErrorExpectedNumber)
	return NInvalidNumber
}

func (p *Parser) getOperatorNodeType() NodeType {
	switch p.currentType() {
	case TOperatorPlus:
		return NAddition
	case TOperatorMinus:
		return NSubtraction
	case TOperatorMult:
		return NMultiplication
	case TOperatorDiv:
		return NDivision
	}

	p.pushError(ErrorExpectedOperator)
	return NInvalidOperator
}

func parseStart(p *Parser) parseState {
	if p.currentType() == TLeftBracket {
		p2 := &Parser{p.tokens, p.topNode, p.current, p.index, nil, true}
		p.topNode, p.index = p2.run()
		p.pushErrors(p2.errors)

		return parseOperatorAfterRightBracket
	}

	nodeType := p.getNumberNodeType()

	p.topNode = p.newNode(nodeType)

	return parseOperator
}

func parseNumberOrLeftBracket(p *Parser) parseState {
	if p.currentType() == TLeftBracket {
		p2 := &Parser{p.tokens, p.topNode, p.current, p.index, nil, true}
		p.current.RightChild, p.index = p2.run()
		p.pushErrors(p2.errors)

		return parseOperator
	}

	nodeType := p.getNumberNodeType()

	p.current.RightChild = p.newNode(nodeType)
	return parseOperator
}

func parseOperator(p *Parser) parseState {
	if p.currentType() == TRightBracket {
		if !p.nested {
			p.pushError(ErrorUnexpectedClosingBracket)
		}
		return nil
	}

	nodeType := p.getOperatorNodeType()

	node := p.newNode(nodeType)
	if IsOperator(p.topNode.Type) && isHigherOperator(p.topNode.Type, nodeType) {
		node.LeftChild = p.topNode.RightChild
		p.topNode.RightChild = node
	} else {
		node.LeftChild = p.topNode
		p.topNode = node
	}
	p.current = node

	return parseNumberOrLeftBracket
}

func parseOperatorAfterRightBracket(p *Parser) parseState {
	if p.currentType() == TRightBracket {
		if !p.nested {
			p.pushError(ErrorUnexpectedClosingBracket)
		}
		return nil
	}

	nodeType := p.getOperatorNodeType()

	node := p.newNode(nodeType)
	node.LeftChild = p.topNode
	p.topNode = node
	p.current = node

	return parseNumberOrLeftBracket
}
