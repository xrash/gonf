package gonf

type tokenType int

type token struct {
	tokenType tokenType
	value string
}

const (
	t_EOF            = iota
	t_KEY
	t_VALUE
	t_MAP_START
	t_MAP_END
	t_ARRAY_START
	t_ARRAY_END
)
