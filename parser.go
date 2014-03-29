package gonf

import (
	"errors"
)

type parser struct {
	tokens chan token
}

func newParser(c chan token) *parser {
	return &parser{c}
}

func (p *parser) parse() (*Config, error) {
	var key string
	var cfg *Config
	cfg = new(Config)
	cfg.table = make(map[string]*Config)

	for {
		select {
		case t := <-p.tokens:
			switch t.tokenType {
			case T_KEY:
				key = t.value
			case T_VALUE:
				if cfg.table == nil {
					cfg.array[len(cfg.array)] = &Config{value:t.value}
				} else {
					cfg.table[key] = &Config{value: t.value}
				}
			case T_TABLE_START:
				if cfg.table == nil {
					cfg.array[len(cfg.array)] = new(Config)
					cfg.array[len(cfg.array)-1].parent = cfg
					cfg = cfg.array[len(cfg.array)-1]
					cfg.table = make(map[string]*Config)
				} else {
					cfg.table[key] = new(Config)
					cfg.table[key].parent = cfg
					cfg = cfg.table[key]
					cfg.table = make(map[string]*Config)
				}
			case T_TABLE_END:
				cfg = cfg.parent
			case T_ARRAY_START:
				if cfg.table == nil {
					cfg.array[len(cfg.array)] = new(Config)
					cfg.array[len(cfg.array)-1].parent = cfg
					cfg = cfg.array[len(cfg.array)-1]
					cfg.array = make(map[int]*Config)
				} else {
					cfg.table[key] = new(Config)
					cfg.table[key].parent = cfg
					cfg = cfg.table[key]
					cfg.array = make(map[int]*Config)
				}
			case T_ARRAY_END:
				cfg = cfg.parent
			case T_EOF:
				if t.value != "" {
					e := errors.New(t.value)
					return nil, e
				}
				return cfg, nil
			}
		}
	}
}
