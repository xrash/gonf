package parser

import (
	"github.com/xrash/gonf/tokens"
	"fmt"
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

	fmt.Println()
	fmt.Println("=============")
	fmt.Println()

	fmt.Println(p.tree.key.value)
	fmt.Println(p.tree.value.string.value)
	fmt.Println(p.tree.pair.key.value)
	fmt.Println(p.tree.pair.value.array.values.value.table.pair.pair.value.string.value)
	fmt.Println(p.tree.pair.value.array.values.values.value.table.pair.value.string.value)

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
