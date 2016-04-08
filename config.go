package gonf

import (
	"errors"
	"github.com/xrash/gonf/lexer"
	"github.com/xrash/gonf/parser"
	"github.com/xrash/gonf/tokens"
	"io"
	"io/ioutil"
	"reflect"
	"strconv"
)

type Config struct {
	parent *Config
	string string
	table  map[string]*Config
	array  []*Config
}

func Read(r io.Reader) (*Config, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	input := string(b)
	tokens := make(chan tokens.Token, len(input)/8)

	l := lexer.NewLexer(input, tokens)
	p := parser.NewParser(tokens)

	go l.Lex()

	var tree *parser.PairNode
	tree, err = p.Parse()

	if err != nil {
		return nil, err
	}

	return generate(tree), nil
}

func (c *Config) Length() int {
	if c.IsArray() {
		return len(c.array)
	} else if c.IsTable() {
		return len(c.table)
	}

	return 0
}

func (c *Config) IsTable() bool {
	return c.table != nil
}

func (c *Config) IsArray() bool {
	return c.array != nil
}

func (c *Config) IsString() bool {
	return !c.IsTable() && !c.IsArray()
}

func (c *Config) Parent() *Config {
	return c.parent
}

func (c *Config) TraverseTable(visit func(string, *Config)) {
	for key, value := range c.table {
		visit(key, value)
	}
}

func (c *Config) TraverseArray(visit func(int, *Config)) {
	for key, value := range c.array {
		visit(key, value)
	}
}

func (c *Config) Table(args ...interface{}) (map[string]*Config, error) {
	c, err := c.Get(args...)
	if err != nil {
		return nil, err
	}
	return c.table, nil
}

func (c *Config) Array(args ...interface{}) ([]*Config, error) {
	c, err := c.Get(args...)
	if err != nil {
		return nil, err
	}
	return c.array, nil
}

func (c *Config) String(args ...interface{}) (string, error) {
	c, err := c.Get(args...)
	if err != nil {
		return "", err
	}
	return c.string, nil
}

func (c *Config) Int(args ...interface{}) (int, error) {
	s, err := c.String(args...)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(s)
}

func (c *Config) Get(args ...interface{}) (*Config, error) {
	var ok bool
	for _, a := range args {
		switch reflect.TypeOf(a).Kind() {
		case reflect.String:
			if c, ok = c.table[a.(string)]; !ok {
				return nil, errors.New("key " + a.(string) + " not found")
			}
		case reflect.Int:
			if a.(int) >= len(c.array) {
				return nil, errors.New("index " + strconv.Itoa(a.(int)) + " not found")
			}
			c = c.array[a.(int)]
		}
	}
	return c, nil
}

func (c *Config) Map(s interface{}) error {
	t := reflect.TypeOf(s)

	if t.Kind() != reflect.Ptr {
		return errors.New("The argument to Map must be a pointer")
	}

	t = t.Elem()

	if t.Kind() != reflect.Struct {
		return errors.New("The argument to Map must be a struct")
	}

	v := reflect.ValueOf(s).Elem()

	c.rmap(t, v)

	return nil
}

func (c *Config) rmap(t reflect.Type, v reflect.Value) {
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if tag := field.Tag.Get("gonf"); tag != "" {
				f := v.FieldByName(field.Name)

				switch field.Type.Kind() {

				case reflect.String:
					if c, err := c.Get(tag); err == nil {
						if c.IsArray() {
							if c, err := c.Get(c.Length() - 1); err == nil {
								if value, err := c.String(); err == nil {
									f.SetString(value)
								}
							}
						} else {
							if value, err := c.String(); err == nil {
								f.SetString(value)
							}
						}
					}

				case reflect.Int:
					if c, err := c.Get(tag); err == nil {
						if c.IsArray() {
							if c, err := c.Get(c.Length() - 1); err == nil {
								if value, err := c.Int(); err == nil {
									f.SetInt(int64(value))
								}
							}
						} else {
							if value, err := c.Int(); err == nil {
								f.SetInt(int64(value))
							}
						}
					}

				case reflect.Struct:
					if c, err := c.Get(tag); err == nil {
						if c.IsArray() {
							if c, err := c.Get(c.Length() - 1); err == nil {
								if c, err := c.Get(); err == nil {
									c.rmap(f.Type(), f)
								}
							}
						} else {
							if c, err := c.Get(); err == nil {
								c.rmap(f.Type(), f)
							}
						}
					}

				case reflect.Slice:
					if c, err := c.Get(tag); err == nil {
						c.rmap(f.Type(), f)
					}

				case reflect.Map:
					if c, err := c.Get(tag); err == nil {
						c.rmap(f.Type(), f)
					}
				}
			}
		}

	} else if t.Kind() == reflect.Slice {
		if c.IsArray() {
			v.Set(reflect.MakeSlice(v.Type(), c.Length(), c.Length()))
			for i := 0; i < v.Len(); i++ {
				c, _ := c.Get(i)
				mapArrayElement(t, v, c, i)
			}
		} else {
			v.Set(reflect.MakeSlice(v.Type(), 1, 1))
			mapArrayElement(t, v, c, 0)
		}

	} else if t.Kind() == reflect.Map {
		if t.Key() != reflect.TypeOf("") {
			return
		}

		if c.IsTable() {
			v.Set(reflect.MakeMap(v.Type()))
			table, _ := c.Table()
			for key := range table {
				c, _ := c.Get(key)
				switch t.Elem() {

				case reflect.TypeOf(""):
					if value, err := c.String(); err == nil {
						v.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
					}

				case reflect.TypeOf(0):
					if value, err := c.Int(); err == nil {
						v.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
					}

				default:
					v.SetMapIndex(reflect.ValueOf(key), reflect.New(t.Elem()).Elem())
					//					c.rmap(v.MapIndex(reflect.ValueOf(key)).Type(), v.MapIndex(reflect.ValueOf(key)))
				}
			}
		}
	}

}

func mapArrayElement(t reflect.Type, v reflect.Value, c *Config, i int) {
	switch t.Elem() {
	case reflect.TypeOf(""):
		if value, err := c.String(); err == nil {
			v.Index(i).SetString(value)
		}
	case reflect.TypeOf(0):
		if value, err := c.Int(); err == nil {
			v.Index(i).SetInt(int64(value))
		}
	default:
		c.rmap(v.Index(i).Type(), v.Index(i))
	}
}
