package parser

import (
	"errors"

	"github.com/relnod/calcgo/lexer"
)

type parseState func(*Parser) parseState

// Parser holds state of parser
type Parser struct {
	tokens    chan lexer.Token
	currToken lexer.Token
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
	NInvalidFunction
	NInteger
	NDecimal
	NVariable
	NAddition
	NFuncSqrt
	NSubtraction
	NMultiplication
	NDivision
)

// Errors retured by the parser
var (
	ErrorExpectedNumberOrVariable = errors.New("Error: Expected number or variable got something else")
	ErrorExpectedOperator         = errors.New("Error: Expected operator got something else")
	ErrorExpectedFunction         = errors.New("Error: Expected function got something else")
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

	lexer := lexer.NewLexer(str)
	lexer.Start()

	return ParseTokenStream(lexer.GetChanel())
}

// ParseTokens parses a list of tokens to an ast
func ParseTokens(tokens []lexer.Token) (AST, []error) {
	if tokens == nil {
		return AST{}, nil
	}

	c := make(chan lexer.Token, len(tokens))
	for _, token := range tokens {
		c <- token
	}
	close(c)

	return ParseTokenStream(c)
}

// ParseTokenStream parses a stream of tokens
func ParseTokenStream(c chan lexer.Token) (AST, []error) {
	p := &Parser{tokens: c}
	p.run()

	return AST{p.topNode}, p.errors
}

func (p *Parser) run() {
	for state := parseStart; state != nil; {
		p.next()
		if p.currToken.Type == lexer.TEOF {
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
	nt, ok := getOperatorNodeType(p.currToken)
	if !ok {
		p.pushError(ErrorExpectedOperator)
	}

	return &Node{nt, p.currToken.Value, nil, nil}
}

func (p *Parser) newNumberOrVariableNode() *Node {
	nt, ok := getNumberOrVariableNodeType(p.currToken)
	if !ok {
		p.pushError(ErrorExpectedNumberOrVariable)
	}

	return &Node{nt, p.currToken.Value, nil, nil}
}

func (p *Parser) newFunctionNode() *Node {
	nt, ok := getFunctionNodeType(p.currToken)
	if !ok {
		p.pushError(ErrorExpectedFunction)
	}

	return &Node{nt, p.currToken.Value, nil, nil}
}

func parseStart(p *Parser) parseState {
	if p.currToken.Type == lexer.TLeftBracket {
		p2 := &Parser{p.tokens, lexer.Token{}, p.topNode, p.current, nil, true}
		p2.run()
		p.topNode = p2.topNode
		p.pushErrors(p2.errors)

		return parseOperatorAfterRightBracket
	}

	if p.currToken.Type == lexer.TFuncSqrt {
		p.topNode = p.newFunctionNode()
		p.current = p.topNode
		p2 := &Parser{p.tokens, lexer.Token{}, p.topNode, p.current, nil, true}
		p2.run()
		p.current.LeftChild = p2.topNode
		p.pushErrors(p2.errors)

		return parseOperator
	}

	p.topNode = p.newNumberOrVariableNode()

	return parseOperator
}

func parseNumberOrLeftBracket(p *Parser) parseState {
	if p.currToken.Type == lexer.TLeftBracket {
		p2 := &Parser{p.tokens, lexer.Token{}, p.topNode, p.current, nil, true}
		p2.run()
		p.current.RightChild = p2.topNode
		p.pushErrors(p2.errors)

		return parseOperator
	}

	if p.currToken.Type == lexer.TFuncSqrt {
		node := p.newFunctionNode()
		p.current.RightChild = node
		p.current = node
		p2 := &Parser{p.tokens, lexer.Token{}, p.topNode, p.current, nil, true}
		p2.run()
		p.current.LeftChild = p2.topNode
		p.pushErrors(p2.errors)

		return parseOperator
	}

	p.current.RightChild = p.newNumberOrVariableNode()
	return parseOperator
}

func parseOperator(p *Parser) parseState {
	if p.currToken.Type == lexer.TRightBracket {
		if !p.nested {
			p.pushError(ErrorUnexpectedClosingBracket)
		}
		return nil
	}

	node := p.newOperatorNode()
	if p.topNode.IsOperator() && p.topNode.isHigherOperator(node) {
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
	if p.currToken.Type == lexer.TRightBracket {
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
