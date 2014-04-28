package tokens

type TokenType int

type Token struct {
	tokenType TokenType
	value string
	line int
	column int
}

func NewToken(t TokenType, v string, l, c int) Token {
	return Token{t, v, l, c}
}

func (t Token) Type() TokenType {
	return t.tokenType
}

func (t Token) Value() string {
	return t.value
}

func (t Token) Line() int {
	return t.line
}

func (t Token) Column() int {
	return t.column
}

func (t Token) String() string {
	switch t.Type() {
	case T_EOF:
		return "EOF"
	case T_STRING:
		return "'" + t.value + "'"
	case T_ARRAY_START:
		return "'['"
	case T_ARRAY_END:
		return "']'"
	case T_TABLE_START:
		return "'{'"
	case T_TABLE_END:
		return "'}'"
	}

	return ""
}
