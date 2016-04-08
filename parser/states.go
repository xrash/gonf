package parser

import (
	"errors"
	"fmt"
	"github.com/xrash/gonf/tokens"
	"strings"
)

type state func(p *Parser) error

func err(got tokens.Token, expected ...string) error {
	msg := "Expected %s at line %d:%d. Got %s."

	for k, e := range expected {
		expected[k] = "'" + e + "'"
	}

	return errors.New(fmt.Sprintf(msg, strings.Join(expected, " OR "), got.Line(), got.Column(), got))
}

func pairState(p *Parser) error {
	token := p.lookup()

	switch token.Type() {
	case tokens.T_STRING:
		p.stack.push(buildPairNode)
		p.stack.push(pairState)
		p.stack.push(valueState)
		p.stack.push(keyState)
	case tokens.T_QUOTE:
		p.stack.push(buildPairNode)
		p.stack.push(pairState)
		p.stack.push(valueState)
		p.stack.push(keyState)
	case tokens.T_TABLE_END:
	case tokens.T_EOF:
	default:
		return err(token, "STRING", "{", "EOF", "\"")
	}

	return nil
}

func keyState(p *Parser) error {
	token := p.lookup()

	switch token.Type() {
	case tokens.T_STRING:
		p.stack.push(buildKeyNode)
		p.stack.push(stringState)
	case tokens.T_QUOTE:
		p.stack.push(buildKeyNode)
		p.stack.push(stringState)
	default:
		return err(token, "STRING", "\"")
	}

	return nil
}

func valueState(p *Parser) error {
	token := p.lookup()

	switch token.Type() {
	case tokens.T_STRING:
		p.stack.push(buildValueStringNode)
		p.stack.push(stringState)
	case tokens.T_QUOTE:
		p.stack.push(buildValueStringNode)
		p.stack.push(stringState)
	case tokens.T_ARRAY_START:
		p.stack.push(buildValueArrayNode)
		p.stack.push(arrayState)
	case tokens.T_TABLE_START:
		p.stack.push(buildValueTableNode)
		p.stack.push(tableState)
	default:
		return err(token, "STRING", "[", "{", "\"")
	}

	return nil
}

func arrayState(p *Parser) error {
	token := p.lookup()

	switch token.Type() {
	case tokens.T_ARRAY_START:
		p.stack.push(buildArrayNode)
		p.stack.push(arrayEndState)
		p.stack.push(valuesState)
		p.stack.push(arrayStartState)
	default:
		return err(token, "[")
	}

	return nil
}

func valuesState(p *Parser) error {
	token := p.lookup()

	switch token.Type() {
	case tokens.T_STRING:
		p.stack.push(buildValuesNode)
		p.stack.push(valuesState)
		p.stack.push(valueState)
	case tokens.T_QUOTE:
		p.stack.push(buildValuesNode)
		p.stack.push(valuesState)
		p.stack.push(valueState)
	case tokens.T_ARRAY_START:
		p.stack.push(buildValuesNode)
		p.stack.push(valuesState)
		p.stack.push(valueState)
	case tokens.T_TABLE_START:
		p.stack.push(buildValuesNode)
		p.stack.push(valuesState)
		p.stack.push(valueState)
	case tokens.T_ARRAY_END:
	default:
		return err(token, "STRING", "[", "{", "]", "\"")
	}

	return nil
}

func tableState(p *Parser) error {
	token := p.lookup()

	switch token.Type() {
	case tokens.T_TABLE_START:
		p.stack.push(buildTableNode)
		p.stack.push(tableEndState)
		p.stack.push(pairState)
		p.stack.push(tableStartState)
	default:
		return err(token, "{")
	}

	return nil
}

func stringState(p *Parser) error {
	token := p.lookup()

	switch token.Type() {
	case tokens.T_STRING:
		p.stack.push(unquotedStringState)
	case tokens.T_QUOTE:
		p.stack.push(quotedStringState)
	default:
		return err(token, "STRING", "\"")
	}

	return nil
}

func unquotedStringState(p *Parser) error {
	token := p.lookup()

	switch token.Type() {
	case tokens.T_STRING:
		p.stack.push(literalState)
	default:
		return err(token, "STRING")
	}

	return nil
}

func quotedStringState(p *Parser) error {
	token := p.lookup()

	switch token.Type() {
	case tokens.T_QUOTE:
		p.stack.push(quoteState)
		p.stack.push(literalState)
		p.stack.push(quoteState)
	default:
		return err(token, "STRING")
	}

	return nil
}

func literalState(p *Parser) error {
	token := p.lookup()

	switch token.Type() {
	case tokens.T_STRING:
		p.next()
	default:
		return err(token, "STRING")
	}

	p.nodeStack.push(NewStringNode(token.Value()))

	return nil
}

func quoteState(p *Parser) error {
	token := p.lookup()

	switch token.Type() {
	case tokens.T_QUOTE:
		p.next()
	default:
		return err(token, "\"")
	}

	return nil
}

func arrayStartState(p *Parser) error {
	token := p.lookup()

	switch token.Type() {
	case tokens.T_ARRAY_START:
		p.next()
	default:
		return err(token, "[")
	}

	return nil
}

func arrayEndState(p *Parser) error {
	token := p.lookup()

	switch token.Type() {
	case tokens.T_ARRAY_END:
		p.next()
	default:
		return err(token, "]")
	}

	return nil
}

func tableStartState(p *Parser) error {
	token := p.lookup()

	switch token.Type() {
	case tokens.T_TABLE_START:
		p.next()
	default:
		return err(token, "{")
	}

	return nil
}

func tableEndState(p *Parser) error {
	token := p.lookup()

	switch token.Type() {
	case tokens.T_TABLE_END:
		p.next()
	default:
		return err(token, "}")
	}

	return nil
}

func buildPairNode(p *Parser) error {
	var pn *PairNode

	if node := p.nodeStack.pop(); node.Kind() == PAIR_NODE {
		pn = node.(*PairNode)
	} else {
		p.nodeStack.push(node)
	}

	vn := p.nodeStack.pop().(*ValueNode)
	kn := p.nodeStack.pop().(*KeyNode)

	p.nodeStack.push(NewPairNode(kn, vn, pn))

	return nil
}

func buildValueStringNode(p *Parser) error {
	sn := p.nodeStack.pop().(*StringNode)
	p.nodeStack.push(NewValueNode(sn, nil, nil))

	return nil
}

func buildValueArrayNode(p *Parser) error {
	vn := p.nodeStack.pop().(*ValuesNode)
	an := NewArrayNode(vn)
	p.nodeStack.push(NewValueNode(nil, nil, an))

	return nil
}

func buildKeyNode(p *Parser) error {
	sn := p.nodeStack.pop().(*StringNode)
	p.nodeStack.push(NewKeyNode(sn))

	return nil
}

func buildValuesNode(p *Parser) error {
	var values *ValuesNode

	if node := p.nodeStack.pop(); node.Kind() == VALUES_NODE {
		values = node.(*ValuesNode)
	} else {
		p.nodeStack.push(node)
	}

	vn := p.nodeStack.pop().(*ValueNode)
	p.nodeStack.push(NewValuesNode(vn, values))

	return nil
}

func buildArrayNode(p *Parser) error {
	return nil
}

func buildValueTableNode(p *Parser) error {
	tn := p.nodeStack.pop().(*TableNode)
	p.nodeStack.push(NewValueNode(nil, tn, nil))

	return nil
}

func buildTableNode(p *Parser) error {
	pn := p.nodeStack.pop().(*PairNode)
	p.nodeStack.push(NewTableNode(pn))

	return nil
}
