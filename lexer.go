package gonf

import (
	"unicode/utf8"
	"fmt"
)

type state func(l *lexer) state

type stateStack []state

func (st *stateStack) push(s state) {
	*st = append(*st, s)
}

func (st *stateStack) pop() state {
	l := len(*st)-1
	s := (*st)[l]
	*st = (*st)[:l]
	return s
}

type lexer struct {
	input string
	start int
	pos int
	line int
	column int
	width int
	state state
	stack *stateStack
	tokens chan token
}

func newLexer(s string, c chan token) *lexer {
	l := new(lexer)
	l.input = s
	l.tokens = c
	l.line = 1
	l.stack = new(stateStack)
	return l
}

func (l *lexer) lex() {
	l.state = searchingKeyState
	for state := l.state; state != nil; {
		state = state(l)
	}
}

// the optional string act as an error message
func (l *lexer) finish(s ...string) {
	var e string
	if len(s) > 0 {
		e = fmt.Sprintf("Syntax error on line %v column %v: %v", l.line, l.column, s[0])
	}
	l.tokens <- token{T_EOF, e}
}

func (l *lexer) emit(t tokenType) {
	tok := token{t, l.input[l.start:(l.pos-1)]}
	l.tokens <- tok
	l.start = l.pos
}

func (l *lexer) next() (r rune) {
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

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) backup() {
	l.column--
	l.pos -= l.width
	var r rune
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	if r == '\n' {
		l.line--
		// lost the column counter :)
		l.column = 0
	}
}

func (l *lexer) eat() {
	l.input = l.input[:l.pos] + l.input[l.pos+l.width:]
}
