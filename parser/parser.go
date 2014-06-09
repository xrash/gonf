package parser

import (
	"github.com/xrash/gonf/tokens"
)

type Parser struct {
	tokens chan tokens.Token
	token tokens.Token
	stack stateStack
	nodeStack nodeStack
	tree *PairNode
}

func NewParser(t chan tokens.Token) *Parser {
	return &Parser{
		t,
		tokens.Token{},
		newStateStack(),
		newNodeStack(),
		nil,
	}
}

func (p *Parser) Parse() error {
	p.stack.push(pairState)
	p.next()

	for !p.stack.empty() {
		state := p.stack.pop()
		if error := state(p); error != nil {
			return error
		}
	}

	p.tree = p.nodeStack.pop().(*PairNode)

	return nil
}

func (p *Parser) Tree() *PairNode {
	return p.tree
}

func (p *Parser) next() {
	p.token = <- p.tokens
}

func (p *Parser) lookup() tokens.Token {
	return p.token
}
