package parser

func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func isLineBreak(r rune) bool {
	return r == '\n'
}

func isBlank(r rune) bool {
	return isSpace(r) || isLineBreak(r)
}

func searchingKeyState(p *Parser) state {
	r := p.next()

	if r == '#' {
		p.stack.push(p.state)
		p.state = inCommentState
		return p.state
	}

	if r == T_EOF {
		p.finish()
		return nil
	}

	if r == '}' {
		p.emit(T_TABLE_END)
		p.state = p.stack.pop()
		return p.state
	}

	if r == ']' {
		p.emit(T_ARRAY_END)
		p.state = p.stack.pop()
		return p.state
	}

	if r == '"' {
		p.ignore()
		return inQuotedKeyState
	}

	if isBlank(r) {
		p.ignore()
		return searchingKeyState
	}

	return inKeyState
}

func inQuotedKeyState(p *Parser) state {
	r := p.next()

	if r == '"' {
		p.emit(T_KEY)
		return searchingValueState
	}

	if r == '\\' {
		return inQuotedBackslashedKeyState
	}

	return inQuotedKeyState
}

func inQuotedBackslashedKeyState(p *Parser) state {
	r := p.next()

	if r == '"' || r == '\\' {
		p.backup()
		p.backup()
		p.eat()
		p.next()
	}

	return inQuotedKeyState
}

func inKeyState(p *Parser) state {
	r := p.next()

	if isBlank(r) {
		p.emit(T_KEY)
		return searchingValueState
	}

	return inKeyState
}

func searchingValueState(p *Parser) state {
	r := p.next()

	if r == '#' {
		p.stack.push(p.state)
		p.state = inCommentState
		return p.state
	}

	if isBlank(r) {
		p.ignore()
		return searchingValueState
	}

	if r == '"' {
		p.ignore()
		return inQuotedValueState
	}

	if r == '}' {
		p.emit(T_TABLE_END)
		p.state = p.stack.pop()
		return p.state
	}

	if r == ']' {
		p.emit(T_ARRAY_END)
		p.state = p.stack.pop()
		return p.state
	}

	if r == '{' {
		p.emit(T_TABLE_START)
		p.stack.push(p.state)
		p.state = searchingKeyState
		return p.state
	}

	if r == '[' {
		p.emit(T_ARRAY_START)
		p.stack.push(p.state)
		p.state = searchingValueState
		return p.state
	}

	return inValueState
}

func inValueState(p *Parser) state {
	r := p.next()

	if isBlank(r) {
		p.emit(T_VALUE)
		return p.state
	}

	if r == T_EOF {
		p.pos++
		p.emit(T_VALUE)
		p.finish()
		return nil
	}

	return inValueState
}

func inQuotedValueState(p *Parser) state {
	r := p.next()

	if r == '"' {
		p.emit(T_VALUE)
		return p.state
	}

	if r == '\\' {
		return inQuotedBackslashedValueState
	}

	return inQuotedValueState
}

func inQuotedBackslashedValueState(p *Parser) state {
	r := p.next()

	if r == '"' || r == '\\' {
		p.backup()
		p.backup()
		p.eat()
		p.next()
	}

	return inQuotedValueState
}

func inCommentState(p *Parser) state {
	r := p.next()

	if isLineBreak(r) {
		p.ignore()
		p.state = p.stack.pop()
		return p.state
	}

	if r == T_EOF {
		p.finish()
		return nil
	}

	return inCommentState
}
