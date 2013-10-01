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

	for {
		select {
		case t := <-p.tokens:
			switch t.tokenType {
			case t_KEY:
				key = t.value
			case t_VALUE:
				if cfg.map_ == nil {
					cfg.array[len(cfg.array)] = &Config{value:t.value}
				} else {
					cfg.map_[key] = &Config{value: t.value}
				}
			case t_MAP_START:
				if cfg.map_ == nil {
					cfg.array[len(cfg.array)] = new(Config)
					cfg.array[len(cfg.array)-1].parent = cfg
					cfg = cfg.array[len(cfg.array)-1]
					cfg.map_ = make(map[string]*Config)
				} else {
					cfg.map_[key] = new(Config)
					cfg.map_[key].parent = cfg
					cfg = cfg.map_[key]
					cfg.map_ = make(map[string]*Config)
				}
			case t_MAP_END:
				cfg = cfg.parent
			case t_ARRAY_START:
				if cfg.map_ == nil {
					cfg.array[len(cfg.array)] = new(Config)
					cfg.array[len(cfg.array)-1].parent = cfg
					cfg = cfg.array[len(cfg.array)-1]
					cfg.array = make(map[int]*Config)
				} else {
					cfg.map_[key] = new(Config)
					cfg.map_[key].parent = cfg
					cfg = cfg.map_[key]
					cfg.array = make(map[int]*Config)
				}
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
