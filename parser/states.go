package parser

import (
	"fmt"
	"errors"
	"strings"
	"github.com/xrash/gonf/tokens"
)

type state func(p *Parser) error

func err(got tokens.Token, expected ...string) error {
	msg := "Expected %s at line %d:%d. Got %s."

	for k, e := range(expected) {
		expected[k] = "'" + e + "'"
	}

	return errors.New(fmt.Sprintf(msg, strings.Join(expected, " OR "), got.Line(), got.Column(), got))
}

func gonfState(p *Parser) error {
	token := p.lookup()

	switch token.Type() {
	case tokens.T_STRING:
		p.stack.push(pairState)
	case tokens.T_TABLE_END:
	case tokens.T_EOF:
	default:
		return err(token, "STRING", "{", "EOF")
	}

	return nil
}

func pairState(p *Parser) error {
	token := p.lookup()

	switch token.Type() {
	case tokens.T_STRING:
		p.stack.push(gonfState)
		p.stack.push(valueState)
		p.stack.push(keyState)
	default:
		return err(token, "STRING")
	}

	return nil
}

func keyState(p *Parser) error {
	token := p.lookup()

	switch token.Type() {
	case tokens.T_STRING:
		p.stack.push(stringState)
	default:
		return err(token, "STRING")
	}

	return nil
}

func valueState(p *Parser) error {
	token := p.lookup()

	switch token.Type() {
	case tokens.T_STRING:
		p.stack.push(stringState)
	case tokens.T_ARRAY_START:
		p.stack.push(arrayState)
	case tokens.T_TABLE_START:
		p.stack.push(tableState)
	default:
		return err(token, "STRING", "[", "{")
	}

	return nil
}

func arrayState(p *Parser) error {
	token := p.lookup()

	switch token.Type() {
		case tokens.T_ARRAY_START:
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
		p.stack.push(valuesState)
		p.stack.push(valueState)
	case tokens.T_ARRAY_START:
		p.stack.push(valuesState)
		p.stack.push(valueState)
	case tokens.T_TABLE_START:
		p.stack.push(valuesState)
		p.stack.push(valueState)
	case tokens.T_ARRAY_END:
	default:
		return err(token, "STRING", "[", "{", "]")
	}

	return nil
}

func tableState(p *Parser) error {
	token := p.lookup()

	switch token.Type() {
	case tokens.T_TABLE_START:
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
		p.next()
	default:
		return err(token, "STRING")
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
