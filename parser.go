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
	tokens    chan Token
	currToken Token
	topNode   *Node
	current   *Node
	errors    []error
	nested    bool
}

// Node types
const (
	NError NodeType = iota
	NInvalidNumber
	NInvalidVariable
	NInvalidOperator
	NInteger
	NDecimal
	NVariable
	NAddition
	NSubtraction
	NMultiplication
	NDivision
)

// Errors retured by the parser
var (
	ErrorExpectedNumberOrVariable = errors.New("Error: Expected number or variable got something else")
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
	if len(str) == 0 {
		return AST{}, nil
	}

	lexer := NewLexer(str)
	lexer.Start()

	return ParseTokenStream(lexer.GetChanel())
}

// ParseTokens parses a list of tokens to an ast
func ParseTokens(tokens []Token) (AST, []error) {
	if tokens == nil {
		return AST{}, nil
	}

	c := make(chan Token, len(tokens))
	for _, token := range tokens {
		c <- token
	}
	close(c)

	return ParseTokenStream(c)
}

// ParseTokenStream parses a stream of tokens
func ParseTokenStream(c chan Token) (AST, []error) {
	p := &Parser{tokens: c}
	p.run()

	return AST{p.topNode}, p.errors
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

func (p *Parser) run() {
	for state := parseStart; state != nil; {
		p.next()
		if p.currToken.Type == TEOF {
			if p.nested {
				p.pushError(ErrorMissingClosingBracket)
			}
			break
		}

		state = state(p)
	}
}

func (p *Parser) next() {
	p.currToken = <-p.tokens
}

func (p *Parser) pushError(err error) {
	p.errors = append(p.errors, err)
}

func (p *Parser) pushErrors(errors []error) {
	for _, err := range errors {
		p.pushError(err)
	}
}

func (p *Parser) newOperatorNode() *Node {
	return &Node{p.getOperatorNodeType(), p.currToken.Value, nil, nil}
}

func (p *Parser) newNumberOrVariableNode() *Node {
	return &Node{p.getNumberOrVariableNodeType(), p.currToken.Value, nil, nil}
}

func (p *Parser) getNumberOrVariableNodeType() NodeType {
	switch p.currToken.Type {
	case TInteger:
		return NInteger
	case TDecimal:
		return NDecimal
	case TVariable:
		return NVariable
	}

	p.pushError(ErrorExpectedNumberOrVariable)
	switch p.currToken.Type {
	case TInvalidCharacterInNumber:
		return NInvalidNumber
	case TInvalidCharacterInVariable:
		return NInvalidVariable
	}
	return NError
}

func (p *Parser) getOperatorNodeType() NodeType {
	switch p.currToken.Type {
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
	if p.currToken.Type == TLeftBracket {
		p2 := &Parser{p.tokens, Token{}, p.topNode, p.current, nil, true}
		p2.run()
		p.topNode = p2.topNode
		p.pushErrors(p2.errors)

		return parseOperatorAfterRightBracket
	}

	p.topNode = p.newNumberOrVariableNode()

	return parseOperator
}

func parseNumberOrLeftBracket(p *Parser) parseState {
	if p.currToken.Type == TLeftBracket {
		p2 := &Parser{p.tokens, Token{}, p.topNode, p.current, nil, true}
		p2.run()
		p.current.RightChild = p2.topNode
		p.pushErrors(p2.errors)

		return parseOperator
	}

	p.current.RightChild = p.newNumberOrVariableNode()
	return parseOperator
}

func parseOperator(p *Parser) parseState {
	if p.currToken.Type == TRightBracket {
		if !p.nested {
			p.pushError(ErrorUnexpectedClosingBracket)
		}
		return nil
	}

	node := p.newOperatorNode()
	if IsOperator(p.topNode.Type) && isHigherOperator(p.topNode.Type, node.Type) {
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
	if p.currToken.Type == TRightBracket {
		if !p.nested {
			p.pushError(ErrorUnexpectedClosingBracket)
		}
		return nil
	}

	node := p.newOperatorNode()
	node.LeftChild = p.topNode
	p.topNode = node
	p.current = node

	return parseNumberOrLeftBracket
}
