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
	Key *KeyNode
	Value *ValueNode
	Pair *PairNode
}

func NewPairNode(key *KeyNode, value *ValueNode, pair *PairNode) *PairNode {
	return &PairNode{
		key,
		value,
		pair,
	}
}

func (n *PairNode) Kind() int {
	return PAIR_NODE
}

type KeyNode struct {
	Value *StringNode
}

func NewKeyNode(value *StringNode) *KeyNode {
	return &KeyNode{
		value,
	}
}

func (n *KeyNode) Kind() int {
	return KEY_NODE
}

type ValueNode struct {
	String *StringNode
	Table *TableNode
	Array *ArrayNode
}

func NewValueNode(string *StringNode, table *TableNode, array *ArrayNode) *ValueNode {
	return &ValueNode{
		string,
		table,
		array,
	}
}

func (n *ValueNode) Kind() int {
	return VALUE_NODE
}

type ValuesNode struct {
	Value *ValueNode
	Values *ValuesNode
}

func NewValuesNode(value *ValueNode, values *ValuesNode) *ValuesNode {
	return &ValuesNode{
		value,
		values,
	}
}

func (n *ValuesNode) Kind() int {
	return VALUES_NODE
}

type StringNode struct {
	Value string
}

func NewStringNode(value string) *StringNode {
	return &StringNode{
		value,
	}
}

func (n *StringNode) Kind() int {
	return STRING_NODE
}

type TableNode struct {
	Pair *PairNode
}

func NewTableNode(pair *PairNode) *TableNode {
	return &TableNode{
		pair,
	}
}

func (n *TableNode) Kind() int {
	return TABLE_NODE
}

type ArrayNode struct {
	Values *ValuesNode
}

func NewArrayNode(values *ValuesNode) *ArrayNode {
	return &ArrayNode{
		values,
	}
}

func (n *ArrayNode) Kind() int {
	return ARRAY_NODE
}
