package lexer

type TokenType int

type Token struct {
	tokenType TokenType
	value string
}

func NewToken(t TokenType, s string) Token {
	return Token{t, s}
}

func (t *Token) Type() TokenType {
	return t.tokenType
}

func (t *Token) Value() string {
	return t.value
}

const (
	T_EOF            = iota
	T_KEY
	T_VALUE
	T_TABLE_START
	T_TABLE_END
	T_ARRAY_START
	T_ARRAY_END
)
