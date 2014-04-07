package parser

import (
	"unicode/utf8"
	"fmt"
)

type state func(p *Parser) state

type Parser struct {
	input string
	start int
	pos int
	line int
	column int
	width int
	state state
	stack *stack
	tokens chan Token
}

func NewParser(s string, c chan Token) *Parser {
	p := new(Parser)
	p.input = s
	p.tokens = c
	p.line = 1
	p.stack = new(stack)
	return p
}

func (p *Parser) Parse() {
	p.state = searchingKeyState
	for state := p.state; state != nil; {
		state = state(p)
	}
}

// the optional string act as an error message
func (p *Parser) finish(s ...string) {
	var e string
	if len(s) > 0 {
		e = fmt.Sprintf("Syntax error on line %v column %v: %v", p.line, p.column, s[0])
	}
	p.tokens <- NewToken(T_EOF, e)
}

func (p *Parser) emit(t TokenType) {
	tok := Token{t, p.input[p.start:(p.pos-1)]}
	p.tokens <- tok
	p.start = p.pos
}

func (p *Parser) next() (r rune) {
	if p.pos >= len(p.input) {
		p.width = 0
		return T_EOF
	}

	r, p.width = utf8.DecodeRuneInString(p.input[p.pos:])
	p.pos += p.width

	if r == '\n' {
		p.column = 0
		p.line++
	} else {
		p.column++
	}

	return r
}

func (p *Parser) ignore() {
	p.start = p.pos
}

func (p *Parser) backup() {
	p.column--
	p.pos -= p.width
	var r rune
	r, p.width = utf8.DecodeRuneInString(p.input[p.pos:])
	if isLineBreak(r) {
		p.line--
		// lost the column counter :)
		p.column = 0
	}
}

func (p *Parser) eat() {
	p.input = p.input[:p.pos] + p.input[p.pos+p.width:]
}
