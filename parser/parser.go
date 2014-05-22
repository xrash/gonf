package parser

import (
	"github.com/xrash/gonf/tokens"
)

type Parser struct {
	tokens chan tokens.Token
	token tokens.Token
	stack stack
	tree *Tree
}

func NewParser(t chan tokens.Token) *Parser {
	return &Parser{
		t,
		tokens.Token{},
		newStack(),
		NewTree(),
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

	return nil
}

func (p *Parser) next() {
	p.token = <- p.tokens
}

func (p *Parser) lookup() tokens.Token {
	return p.token
}
