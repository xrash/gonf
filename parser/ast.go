package parser

const (
	VALUE_NODE_TYPE_STRING = iota
	VALUE_NODE_TYPE_ARRAY
	VALUE_NODE_TYPE_TABLE
)

type Node interface {}

type ValueType int

type ValueNode struct {
	string *StringNode
	table *TableNode
	array *ArrayNode
	valueType ValueType
}

func (v *ValueNode) ValueType() ValueType {
	return v.valueType
}

type StringNode struct {
	value string
}

type PairNode struct {
	key string
	value *ValueNode
	pair *PairNode
}

type ValuesNode struct {
	value *ValueNode
	values *ValuesNode
}

type TableNode struct {
	pair *PairNode
}

type ArrayNode struct {
	values *ValuesNode
}

type Tree struct {
	
}

func NewTree() *Tree {
	return &Tree{}
}
