# gonf

Package _gonf_ provides an interface to a simple configuration file format.

Below is a simple example to introduce you to the format.

    # any.conf

    database {
        host 127.0.0.1
        schema test
        auth {
            user testuser
            pass testpass
        }
    }

    fruits [
        pear
        orange
        lemon
        papaya
    ]

As intuitively noted, the format supports tables (maps), arrays and string literals. This should be all you need.

Now, a simple example of code (given the above file):

    package main

    import (
        "os"
        "fmt"
        "github.com/xrash/gonf"
    )

    func main() {
        file, _ := os.Open("any.conf")
        config, _ := gonf.Read(file)

        fmt.Println(config.String("database", "host")) // 127.0.0.1
        fmt.Println(config.String("database", "auth", "user")) // testuser
        fmt.Println(config.String("fruits", 0)) // pear
        fmt.Println(config.String("fruits", 1)) // orange
    }

You can also directly map your config to a struct. Example:

    package main

    import (
        "os"
        "fmt"
        "github.com/xrash/gonf"
    )

    type Database struct {
        Host string `gonf:"host"`
        Schema string `gonf:"schema"`
        Auth struct {
            User string `gonf:"user"`
            Pass string `gonf:"pass"`
        } `gonf:"auth"`
    }

    func main() {
        database := new(Database)

        file, _ := os.Open("any.conf")
        config, _ := gonf.Read(file)
        config, _ = config.Get("database")

        config.Map(database)
        fmt.Println(database.Schema) // test
        fmt.Println(database.Auth.User) // testuser
    }

NOTE: The struct fields have to be exported so the Map function can see them through reflection

You are encouraged to see the working examples of tests/gonf_test.go.

Here is the LL(1) grammar:

    pair -> key value pair | &
    key -> string
    value -> table | array | string
    table -> { pair }
    array -> [ values ]
    values -> value values | &

[the golang string literal specification](http://golang.org/ref/spec#String_literals)

## TODO
 - Detect syntax errors and generate nice error messages
 - Write a real spec
 - Think about the possibility of extending the format to support including external files (something like @include)
