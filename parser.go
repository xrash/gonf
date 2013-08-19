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
	cfg.map_ = make(map[string]*Config)
	var i int

	for {
		select {
		case t := <-p.tokens:
			switch t.tokenType {
			case t_KEY:
				key = t.value
			case t_VALUE:
				if cfg.map_ == nil {
					cfg.array[i] = &Config{value:t.value}
					i++
				} else {
					cfg.map_[key] = new(Config)
					cfg.map_[key].value = t.value
				}
			case t_MAP_START:
				cfg.map_[key] = new(Config)
				cfg.map_[key].parent = cfg
				cfg = cfg.map_[key]
				cfg.map_ = make(map[string]*Config)
			case t_MAP_END:
				cfg = cfg.parent
			case t_ARRAY_START:
				i = 0
				cfg.map_[key] = new(Config)
				cfg.map_[key].parent = cfg
				cfg = cfg.map_[key]
				cfg.array = make([]*Config, 255)
			case t_ARRAY_END:
				cfg = cfg.parent
			case t_EOF:
				if t.value != "" {
					e := errors.New(t.value)
					return nil, e
				}
				return cfg, nil
			}
		}
	}
}
