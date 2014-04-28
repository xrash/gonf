package parser

import (
	"fmt"
	"runtime"
	"reflect"
	"github.com/xrash/gonf/tokens"
)

type Parser struct {
	tokens chan tokens.Token
	token tokens.Token
	stack stack
}

func NewParser(t chan tokens.Token) *Parser {
	return &Parser{
		t,
		tokens.Token{},
		newStack(),
	}
}

func (p *Parser) Parse() error {
	p.stack.push(pairState)
	p.next()

	for !p.stack.empty() {
		state := p.stack.pop()
		fmt.Println(fname(state))
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

func fname(i interface{}) string {
    return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
