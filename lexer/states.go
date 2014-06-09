package lexer

import (
	"github.com/xrash/gonf/tokens"
)

type state func(l *Lexer) state

func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func isLineBreak(r rune) bool {
	return r == '\n'
}

func isEOF(r rune) bool {
	return r == 0
}

func isBlank(r rune) bool {
	return isSpace(r) || isLineBreak(r) || isEOF(r)
}

func searchingState(l *Lexer) state {
	r := l.next()

	if isEOF(r) {
		return nil
	}

	if isBlank(r) {
		l.ignore()
		return searchingState
	}

	if r == '#' {
		return inCommentState
	}

	if r == '{' {
		l.emit(tokens.T_TABLE_START)
		return searchingState
	}

	if r == '[' {
		l.emit(tokens.T_ARRAY_START)
		return searchingState
	}

	if r == '}' {
		l.emit(tokens.T_TABLE_END)
		return searchingState
	}

	if r == ']' {
		l.emit(tokens.T_ARRAY_END)
		return searchingState
	}

	if r == '"' {
		l.ignore()
		return inQuotedStringState
	}

	return inStringState
}

func inCommentState(l *Lexer) state {
	r := l.next()

	if isLineBreak(r) {
		l.ignore()
		return searchingState
	}

	return inCommentState
}

func inQuotedStringState(l *Lexer) state {
	r := l.next()

	if r == '"' {
		l.emit(tokens.T_STRING)
		return searchingState
	}

	if r == '\\' {
		return inBackslashedStringState
	}

	return inQuotedStringState
}

func inBackslashedStringState(l *Lexer) state {
	l.backup()
	l.eat()
	l.next()
	return inQuotedStringState
}

func inStringState(l *Lexer) state {
	r := l.next()

	if isBlank(r) {
		l.emit(tokens.T_STRING)
		return searchingState
	}

	return inStringState
}
