package lexer

func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func isLineBreak(r rune) bool {
	return r == '\n'
}

func isBlank(r rune) bool {
	return isSpace(r) || isLineBreak(r)
}

func searchingKeyState(l *Lexer) state {
	r := l.next()

	if r == '#' {
		l.stack.push(l.state)
		l.state = inCommentState
		return l.state
	}

	if r == T_EOF {
		l.finish()
		return nil
	}

	if r == '}' {
		l.emit(T_TABLE_END)
		l.state = l.stack.pop()
		return l.state
	}

	if r == ']' {
		l.emit(T_ARRAY_END)
		l.state = l.stack.pop()
		return l.state
	}

	if r == '"' {
		l.ignore()
		return inQuotedKeyState
	}

	if isBlank(r) {
		l.ignore()
		return searchingKeyState
	}

	return inKeyState
}

func inQuotedKeyState(l *Lexer) state {
	r := l.next()

	if r == '"' {
		l.emit(T_KEY)
		return searchingValueState
	}

	if r == '\\' {
		return inQuotedBackslashedKeyState
	}

	return inQuotedKeyState
}

func inQuotedBackslashedKeyState(l *Lexer) state {
	r := l.next()

	if r == '"' || r == '\\' {
		l.backup()
		l.backup()
		l.eat()
		l.next()
	}

	return inQuotedKeyState
}

func inKeyState(l *Lexer) state {
	r := l.next()

	if isBlank(r) {
		l.emit(T_KEY)
		return searchingValueState
	}

	return inKeyState
}

func searchingValueState(l *Lexer) state {
	r := l.next()

	if r == '#' {
		l.stack.push(l.state)
		l.state = inCommentState
		return l.state
	}

	if isBlank(r) {
		l.ignore()
		return searchingValueState
	}

	if r == '"' {
		l.ignore()
		return inQuotedValueState
	}

	if r == '}' {
		l.emit(T_TABLE_END)
		l.state = l.stack.pop()
		return l.state
	}

	if r == ']' {
		l.emit(T_ARRAY_END)
		l.state = l.stack.pop()
		return l.state
	}

	if r == '{' {
		l.emit(T_TABLE_START)
		l.stack.push(l.state)
		l.state = searchingKeyState
		return l.state
	}

	if r == '[' {
		l.emit(T_ARRAY_START)
		l.stack.push(l.state)
		l.state = searchingValueState
		return l.state
	}

	return inValueState
}

func inValueState(l *Lexer) state {
	r := l.next()

	if isBlank(r) {
		l.emit(T_VALUE)
		return l.state
	}

	if r == T_EOF {
		l.pos++
		l.emit(T_VALUE)
		l.finish()
		return nil
	}

	return inValueState
}

func inQuotedValueState(l *Lexer) state {
	r := l.next()

	if r == '"' {
		l.emit(T_VALUE)
		return l.state
	}

	if r == '\\' {
		return inQuotedBackslashedValueState
	}

	return inQuotedValueState
}

func inQuotedBackslashedValueState(l *Lexer) state {
	r := l.next()

	if r == '"' || r == '\\' {
		l.backup()
		l.backup()
		l.eat()
		l.next()
	}

	return inQuotedValueState
}

func inCommentState(l *Lexer) state {
	r := l.next()

	if isLineBreak(r) {
		l.ignore()
		l.state = l.stack.pop()
		return l.state
	}

	if r == T_EOF {
		l.finish()
		return nil
	}

	return inCommentState
}
