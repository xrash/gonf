package gonf

import (
	"io"
	"io/ioutil"
	"reflect"
	"strconv"
	"errors"
	"github.com/xrash/gonf/tokens"
	"github.com/xrash/gonf/lexer"
	"github.com/xrash/gonf/parser"
)

type Config struct {
	parent *Config
	value string
	table map[string]*Config
	array map[int]*Config
}

func Read(r io.Reader) (*Config, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	input := string(b)
	tokens := make(chan tokens.Token)

	l := lexer.NewLexer(input, tokens)
	p := parser.NewParser(tokens)

	go l.Lex()
	err = p.Parse()

	return nil, err
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
	return c.table != nil;
}

func (c *Config) IsArray() bool {
	return c.array != nil;
}

func (c *Config) IsString() bool {
	return !c.IsTable() && !c.IsArray();
}

func (c *Config) Parent() *Config {
	return c.parent
}

func (c *Config) TraverseTable(visit func(string, *Config)) {
	for key, value := range(c.table) {
		visit(key, value)
	}
}

func (c *Config) TraverseArray(visit func(int, *Config)) {
	for key, value := range(c.array) {
		visit(key, value)
	}
}

func (c *Config) Value() string {
	return c.value
}

func (c *Config) Table() map[string]*Config {
	return c.table
}

func (c *Config) Array() map[int]*Config {
	return c.array
}

func (c *Config) String(args ...interface{}) (string, error) {
	c, err := c.Get(args...)
	if err != nil {
		return "", err
	}
	return c.value, nil
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
			if c = c.array[a.(int)]; c == nil {
				return nil, errors.New("index " + strconv.Itoa(a.(int)) + " not found")
			}
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
					if value, err := c.String(tag); err == nil {
						f.SetString(value)
					}
				case reflect.Int:
					if value, err := c.Int(tag); err == nil {
						f.SetInt(int64(value))
					}
				case reflect.Struct:
					if c, err := c.Get(tag); err == nil {
						c.rmap(f.Type(), f)
					}
				case reflect.Slice:
					if c, err := c.Get(tag); err == nil {
						c.rmap(f.Type(), f)
					}
				}
			}
		}
	} else if t.Kind() == reflect.Slice {
		v.Set(reflect.MakeSlice(v.Type(), c.Length(), c.Length()))
		for i := 0; i < v.Len(); i++ {
			c, _ := c.Get(i)
			c.rmap(v.Index(i).Type(), v.Index(i))
		}
	}
}
