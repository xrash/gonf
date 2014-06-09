package tests

import (
	"fmt"
	"strings"
	"testing"
	"github.com/xrash/gonf"
)

const (
	teststring = `
testinteger 22

anarraywithmaps [
    {
        name first
        command second
    }
    {
        name third
        command fourth
    }
]

damn [
    {
        one 1
        two 2
    }
    {
        one 10
        two 20
    }
    [
        42
    ]
]

awesome [
    arrays
    "like this"
    # comment inside an array
]

"keys are any string literal" "just as values"

username www-data
group www-data

anything "any thing"
backslashed "wa\\"
double-backslashed "wa\\\\"
escaped "\""

database {
    host 127.0.0.1
    auth {
        # comment over here
        user testuser
        pass testpass
    }
}

# finish it up
lastline lastvalue`
)

type teststruct struct {
	AnArrayWithMaps []struct {
		Name string `gonf:"name"`
		Command string `gonf:"command"`
	} `gonf:"anarraywithmaps"`
	Username string `gonf:"username"`
	Group string `gonf:"group"`
	Database struct {
		Host string `gonf:"host"`
		Auth struct {
			User string `gonf:"user"`
			Pass string `gonf:"pass"`
		} `gonf:"auth"`
	} `gonf:"database"`
	TestInteger int `gonf:"testinteger"`
	DoNotExists string `gonf:"donotexists"`
}

var config *gonf.Config

func readConfig(t *testing.T) {
	var err error

	config, err = gonf.Read(strings.NewReader(teststring))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
}

func testString(t *testing.T) {
	username, _ := config.String("username")
	group, _ := config.String("group")
	anything, _ := config.String("anything")
	backslashed, _ := config.String("backslashed")
	dblbackslashed, _ := config.String("double-backslashed")
	escaped, _ := config.String("escaped")
	host, _ := config.String("database", "host")
	user, _ := config.String("database", "auth", "user")
	arrays, _ := config.String("awesome", 0)
	likethis, _ := config.String("awesome", 1)
	wololo, _ := config.String("keys are any string literal")
	testinteger, _ := config.Int("testinteger")
	damn0one, _ := config.Int("damn", 0, "one")
	damn1two, _ := config.Int("damn", 1, "two")
	damn20, _ := config.Int("damn", 2, 0)
	donotexists, _ := config.String("donotexists")

	stringtests := map[string]string{
		username: "www-data",
		group: "www-data",
		anything: "any thing",
		backslashed: "wa\\",
		dblbackslashed: "wa\\\\",
		escaped: "\"",
		host: "127.0.0.1",
		user: "testuser",
		arrays: "arrays",
		likethis: "like this",
		wololo: "just as values",
		donotexists: "",
	}

	inttests := map[int]int{
		testinteger: 22,
		damn0one: 1,
		damn1two: 20,
		damn20: 42,
	}

	for i, v := range stringtests {
		if i != v {
			fmt.Printf("%v != %v\n", i, v)
			t.Fail()
		}
	}

	for i, v := range inttests {
		if i != v {
			fmt.Printf("%v != %v\n", i, v)
			t.Fail()
		}
	}
}

func testMap(t *testing.T) {
	s := new(teststruct)

	config.Map(s)

	stringtests := map[string]string{
		s.AnArrayWithMaps[0].Name: "first",
		s.AnArrayWithMaps[1].Command: "fourth",
		s.Username: "www-data",
		s.Group: "www-data",
		s.Database.Host: "127.0.0.1",
		s.Database.Auth.User: "testuser",
	}

	inttests := map[int]int{
		s.TestInteger: 22,
	}

	for i, v := range stringtests {
		if i != v {
			t.Fail()
		}
	}

	for i, v := range inttests {
		if i != v {
			t.Fail()
		}
	}
}

func testLastLine(t *testing.T) {
	v, _ := config.String("lastline")
	if v != "lastvalue" {
		t.Fail()
	}
}

func TestAll(t *testing.T) {
	readConfig(t)
	testString(t)
	testMap(t)
	testLastLine(t)
}
