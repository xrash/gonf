package gonf
/*
import (
	"errors"
	"github.com/xrash/gonf/parser"
)

type generator struct {
	tokens chan parser.Token
}

func newGenerator(c chan parser.Token) *generator {
	return &generator{c}
}

func (p *generator) generate() (*Config, error) {
	var key string
	var cfg *Config
	cfg = new(Config)
	cfg.table = make(map[string]*Config)

	for {
		select {
		case t := <-p.tokens:
			switch t.Type() {
			case parser.T_KEY:
				key = t.Value()
			case parser.T_VALUE:
				if cfg.table == nil {
					cfg.array[len(cfg.array)] = &Config{value:t.Value()}
				} else {
					cfg.table[key] = &Config{value: t.Value()}
				}
			case parser.T_TABLE_START:
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
			case parser.T_TABLE_END:
				cfg = cfg.parent
			case parser.T_ARRAY_START:
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
			case parser.T_ARRAY_END:
				cfg = cfg.parent
			case parser.T_EOF:
				if t.Value() != "" {
					e := errors.New(t.Value())
					return nil, e
				}
				return cfg, nil
			}
		}
	}
}
*/