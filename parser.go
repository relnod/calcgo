package calcgo

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
	index   int
	current *Node
}

// Node types
const (
	NError NodeType = iota
	NInteger
	NDecimal
	NAddition
	NSubtraction
	NMultiplication
	NDivision
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
func Parse(str string) AST {
	tokens := Lex(str)
	return ParseTokens(tokens)
}

// ParseTokens parses a list of tokens to an ast
func ParseTokens(tokens []Token) AST {
	if tokens == nil {
		return AST{}
	}

	p := &Parser{tokens, &Node{}, -1, &Node{}}

	topNode, index := p.run()
	if index > 0 {
		// @todo
	}

	ast := AST{topNode}

	return ast
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

func (p *Parser) getNumberNodeType() NodeType {
	switch p.currentType() {
	case TInteger:
		return NInteger
	case TDecimal:
		return NDecimal
	}

	return NError
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

	return NError
}

func parseStart(p *Parser) parseState {
	if p.currentType() == TLeftBracket {
		p2 := &Parser{p.tokens, p.topNode, p.index, p.current}
		p.topNode, p.index = p2.run()

		return parseOperatorAfterRightBracket
	}

	nodeType := p.getNumberNodeType()

	//@todo: handle wrong node type

	p.topNode = p.newNode(nodeType)

	return parseOperator
}

func parseNumberOrLeftBracket(p *Parser) parseState {
	if p.currentType() == TLeftBracket {
		p2 := &Parser{p.tokens, p.topNode, p.index, p.current}
		p.current.RightChild, p.index = p2.run()

		return parseOperator
	}

	nodeType := p.getNumberNodeType()

	//@todo: handle wrong node type

	p.current.RightChild = p.newNode(nodeType)
	return parseOperator
}

func parseOperator(p *Parser) parseState {
	if p.tokens[p.index].Type == TRightBracket {
		return nil
	}

	nodeType := p.getOperatorNodeType()

	// @todo: handle wrong node type

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
	if p.tokens[p.index].Type == TRightBracket {
		return nil
	}

	nodeType := p.getOperatorNodeType()

	// @todo: handle wrong node type

	node := p.newNode(nodeType)
	node.LeftChild = p.topNode
	p.topNode = node
	p.current = node

	return parseNumberOrLeftBracket
}
