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

testmerge "merge zero"
testmerge "merge one"

testsingleint 1
testsingleint 2

testsinglestruct {
    bozo bozo
}

testsinglestruct {
    bozo wachacha
}

yatm {
    bozo bozoca
    nariz "de pipoca"
    nariz "de papel"
}

yatm {
    naked tongues
}

yatm [
    Perturbator
    ZeroCall
    "Battle of the Young"
]

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

stringarray [
    datum1
    datum2
]

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
	StringArray []string `gonf:"stringarray"`
	TestMerge []string `gonf:"testmerge"`
	TestIntegers []int `gonf:"testinteger"`
	TestSingleInteger int `gonf:"testsingleint"`
	TestSingleStruct struct {
		Bozo string `gonf:"bozo"`
	} `gonf:"testsinglestruct"`
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
		s.StringArray[0]: "datum1",
		s.StringArray[1]: "datum2",
		s.TestMerge[0]: "merge zero",
		s.TestMerge[1]: "merge one",
		s.TestSingleStruct.Bozo: "wachacha",
	}

	inttests := map[int]int{
		s.TestInteger: 22,
		s.TestIntegers[0]: 22,
		s.TestSingleInteger: 2,
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

func testMerge(t *testing.T) {
	testmerge0, _ := config.String("testmerge", 0)
	testmerge1, _ := config.String("testmerge", 1)
	yatm0, _ := config.String("yatm", 0, "bozo")
	yatm1, _ := config.String("yatm", 0, "nariz", 0)
	yatm2, _ := config.String("yatm", 0, "nariz", 1)
	yatm3, _ := config.String("yatm", 1, "naked")
	yatm4, _ := config.String("yatm", 2, 0)
	yatm5, _ := config.String("yatm", 2, 1)
	yatm6, _ := config.String("yatm", 2, 2)

	stringtests := map[string]string{
		testmerge0: "merge zero",
		testmerge1: "merge one",
		yatm0: "bozoca",
		yatm1: "de pipoca",
		yatm2: "de papel",
		yatm3: "tongues",
		yatm4: "Perturbator",
		yatm5: "ZeroCall",
		yatm6: "Battle of the Young",
	}

	for i, v := range stringtests {
		if i != v {
			fmt.Printf("%v != %v\n", i, v)
			t.Fail()
		}
	}
}

func testArray(t *testing.T) {
	a, _ := config.Array("anarraywithmaps")
	a0, _ := a[0].String("name")
	a1, _ := a[1].String("command")

	tests := map[string]string{
		a0: "first",
		a1: "fourth",
	}

	for i, v := range tests {
		if i != v {
			fmt.Printf("%v != %v\n", i, v)
			t.Fail()
		}
	}

	moretests := map[int]string{
		0: "arrays",
		1: "like this",
	}

	awesome, _ := config.Get("awesome")
	awesome.TraverseArray(func(i int, c *gonf.Config) {
		s, _ := c.String()
		if moretests[i] != s {
			fmt.Printf("%v != %v\n", moretests[i], s)
			t.Fail()
		}
	})
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
	testArray(t)
	testMap(t)
	testMerge(t)
	testLastLine(t)
}
