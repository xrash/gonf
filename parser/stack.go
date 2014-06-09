package parser

type stack []interface{}

func newStack() stack {
	return stack{}
}

func (s *stack) push(i interface{}) {
	*s = append(*s, i)
}

func (s *stack) pop() interface{} {
	l := len(*s)-1
	st := (*s)[l]
	*s = (*s)[:l]
	return st
}

func (s *stack) empty() bool {
	return len(*s) == 0
}

type stateStack struct {
	stack stack
}

func newStateStack() stateStack {
	return stateStack{
		newStack(),
	}
}

func (s *stateStack) push(st state) {
	s.stack.push(st)
}

func (s *stateStack) pop() state {
	return s.stack.pop().(state)
}

func (s *stateStack) empty() bool {
	return s.stack.empty()
}

type nodeStack struct {
	stack stack
}

func newNodeStack() nodeStack {
	return nodeStack{
		newStack(),
	}
}

func (s *nodeStack) push(st Node) {
	s.stack.push(st)
}

func (s *nodeStack) pop() Node {
	return s.stack.pop().(Node)
}

func (s *nodeStack) empty() bool {
	return s.stack.empty()
}
