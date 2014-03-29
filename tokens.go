package gonf

type tokenType int

type token struct {
	tokenType tokenType
	value string
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
