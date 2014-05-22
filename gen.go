package gonf

import (
	"github.com/xrash/gonf/tokens"
)

type Generator struct {
	tokens chan tokens.Token
}

func NewGenerator(tokens chan tokens.Token) *Generator {
	return &Generator{tokens}
}

func (g *Generator) Gen() (*Config, error) {
	return nil, nil
}
