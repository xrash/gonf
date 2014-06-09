package parser

const (
	PAIR_NODE = iota
	KEY_NODE
	VALUE_NODE
	STRING_NODE
	VALUES_NODE
	TABLE_NODE
	ARRAY_NODE
)

type Node interface {
	Kind() int
}

type PairNode struct {
	key *KeyNode
	value *ValueNode
	pair *PairNode
}

func NewPairNode(key *KeyNode, value *ValueNode, pair *PairNode) *PairNode {
	return &PairNode{
		key,
		value,
		pair,
	}
}

func (p *PairNode) Kind() int {
	return PAIR_NODE
}

type KeyNode struct {
	value *StringNode
}

func NewKeyNode(s *StringNode) *KeyNode {
	return &KeyNode{
		s,
	}
}

func (k *KeyNode) Kind() int {
	return KEY_NODE
}

type ValueType int

type ValueNode struct {
	string *StringNode
	table *TableNode
	array *ArrayNode
}

func NewValueNode(s *StringNode, t *TableNode, a *ArrayNode) *ValueNode {
	return &ValueNode{
		s,
		t,
		a,
	}
}

func (v *ValueNode) Kind() int {
	return VALUE_NODE
}

type ValuesNode struct {
	value *ValueNode
	values *ValuesNode
}

func NewValuesNode(value *ValueNode, values *ValuesNode) *ValuesNode {
	return &ValuesNode{
		value,
		values,
	}
}

func (v *ValuesNode) Kind() int {
	return VALUES_NODE
}

type StringNode struct {
	value string
}

func (s *StringNode) Kind() int {
	return STRING_NODE
}

func NewStringNode(s string) *StringNode {
	return &StringNode{
		s,
	}
}

type TableNode struct {
	pair *PairNode
}

func NewTableNode(pair *PairNode) *TableNode {
	return &TableNode{
		pair,
	}
}

func (t *TableNode) Kind() int {
	return TABLE_NODE
}

type ArrayNode struct {
	values *ValuesNode
}

func NewArrayNode(vn *ValuesNode) *ArrayNode {
	return &ArrayNode{
		vn,
	}
}

func (a *ArrayNode) Kind() int {
	return ARRAY_NODE
}
