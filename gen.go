package gonf

import (
	"github.com/xrash/gonf/parser"
)

func generate(root *parser.PairNode) *Config {
	return genPairNode(root)
}

func genPairNode(node *parser.PairNode) *Config {
	c := new(Config)
	c.table = make(map[string]*Config)

	for node != nil {
		c.table[node.Key.Value.Value] = genValueNode(node.Value)
		node = node.Pair
	}

	return c
}

func genValueNode(node *parser.ValueNode) *Config {
	var c *Config

	switch {
	case node.String != nil:
		c = genStringNode(node.String)
	case node.Table != nil:
		c = genTableNode(node.Table)
	case node.Array != nil:
		c = genArrayNode(node.Array)
	}

	return c
}

func genStringNode(node *parser.StringNode) *Config {
	c := new(Config)
	c.string = node.Value
	return c
}

func genTableNode(node *parser.TableNode) *Config {
	return genPairNode(node.Pair)
}

func genArrayNode(node *parser.ArrayNode) *Config {
	return genValuesNode(node.Values)
}

func genValuesNode(node *parser.ValuesNode) *Config {
	root := node
	i := 0
	for node != nil {
		i++
		node = node.Values
	}

	c := new(Config)
	c.array = make([]*Config, i)
	node = root
	i = 0

	for node != nil {
		c.array[i] = genValueNode(node.Value)
		i++
		node = node.Values
	}

	return c
}
