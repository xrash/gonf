package parser

type stack []state

func (st *stack) push(s state) {
	*st = append(*st, s)
}

func (st *stack) pop() state {
	l := len(*st)-1
	s := (*st)[l]
	*st = (*st)[:l]
	return s
}
