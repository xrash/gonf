package parser

import (
	"github.com/xrash/gonf/tokens"
)

type Parser struct {
	tokens chan tokens.Token
	token tokens.Token
	stack stateStack
	nodeStack nodeStack
}

func NewParser(t chan tokens.Token) *Parser {
	return &Parser{
		t,
		tokens.Token{},
		newStateStack(),
		newNodeStack(),
	}
}

func (p *Parser) Parse() (*PairNode, error) {
	p.stack.push(pairState)
	p.next()

	for !p.stack.empty() {
		state := p.stack.pop()
		if error := state(p); error != nil {
			return nil, error
		}
	}

	root := p.nodeStack.pop().(*PairNode)

	return root, nil
}

func (p *Parser) next() {
	p.token = <- p.tokens
}

func (p *Parser) lookup() tokens.Token {
	return p.token
}
