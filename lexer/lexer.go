package lexer

import (
	"unicode/utf8"
	"github.com/xrash/gonf/tokens"
)

type Lexer struct {
	input string
	start int
	pos int
	line int
	column int
	width int
	state state
	tokens chan tokens.Token
}

func NewLexer(s string, c chan tokens.Token) *Lexer {
	l := new(Lexer)
	l.input = s + " "
	l.tokens = c
	l.line = 1
	return l
}

func (l *Lexer) Lex() {
	for state := searchingState; state != nil; {
		state = state(l)
	}
}

func (l *Lexer) emit(t tokens.TokenType) {
	token := tokens.NewToken(t, l.input[l.start:(l.pos-1)], l.line, l.column)
	l.tokens <- token
	l.start = l.pos
}

func (l *Lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.tokens <- tokens.NewToken(tokens.T_EOF, "", 0, 0)
		l.width = 0
		return 0
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
		l.column = 666
	}
}

func (l *Lexer) eat() {
	l.input = l.input[:l.pos] + l.input[l.pos+l.width:]
}
