package lexer

import (
	"unicode/utf8"
	"fmt"
)

type state func(l *Lexer) state

type Lexer struct {
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

func NewLexer(s string, c chan Token) *Lexer {
	l := new(Lexer)
	l.input = s
	l.tokens = c
	l.line = 1
	l.stack = new(stack)
	return l
}

func (l *Lexer) Lex() {
	l.state = searchingKeyState
	for state := l.state; state != nil; {
		state = state(l)
	}
}

// the optional string act as an error message
func (l *Lexer) finish(s ...string) {
	var e string
	if len(s) > 0 {
		e = fmt.Sprintf("Syntax error on line %v column %v: %v", l.line, l.column, s[0])
	}
	l.tokens <- NewToken(T_EOF, e)
}

func (l *Lexer) emit(t TokenType) {
	tok := Token{t, l.input[l.start:(l.pos-1)]}
	l.tokens <- tok
	l.start = l.pos
}

func (l *Lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return T_EOF
	}

	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width

	if r == '\n' {
		l.column = 0
		l.line++
	} else {
		l.column++
	}

	return r
}

func (l *Lexer) ignore() {
	l.start = l.pos
}

func (l *Lexer) backup() {
	l.column--
	l.pos -= l.width
	var r rune
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	if isLineBreak(r) {
		l.line--
		// lost the column counter :)
		l.column = 0
	}
}

func (l *Lexer) eat() {
	l.input = l.input[:l.pos] + l.input[l.pos+l.width:]
}
