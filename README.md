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

# Mapping

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

> NOTE: The struct fields have to be exported so the Map function can see them through reflection

# Semantic key merging

One nice feature is the automatic merge of multiple equal keys into an array. Consider the following example:

    song {
	    name "Naked Tongues"
		artist Perturbator
	}
	
	song {
	    name "Battle of the Young"
		artist ZeroCall
	}

This will be translated in a semantic analyzing phase to:

    song [
	    {
		    name "Naked Tongues"
			artist Perturbator
		}
		{
		    name "Battle of the Young"
			artist ZeroCall
		}
	]

And can therefore be accessed like this:

    config.String("song", 0, "name") // Naked Tongues
    config.String("song", 1, "artist") // ZeroCall

# Traversing non-scalar types

A problem that arises in practice is the need to traverse through non-scalar types. In gonf, we got tables and arrays, and both can be traversed. The order of elements in a table may not be guaranteed by the implementation, but the order in an array is expected to be guaranteed in any implementation. There are two ways to traverse through these types:

### Using traversing functions

    config.TraverseTable(func(key string, value gonf.*Config) {
	    fmt.Println(key, value)
	})

    config.TraverseArray(func(index int, value gonf.*Config) {
	    fmt.Println(index, value)
	})

### Using the underlying implemented data structure

    a := config.Array()
	
	for key, value := range a {
	    fmt.Println(key, value)
	}

    t := config.Table()

	for key, value := range t {
	    fmt.Println(key, value)
	}

If you need to check which type the Config object holds, you can use the functions:

    config.IsString()
    config.IsArray()
    config.IsTable()

# More examples

You are encouraged to see the working examples of tests/gonf_test.go.

# Help for implementers

Here is the LL(1) grammar:

    pair -> key value pair | &
    key -> string
    value -> table | array | string
    table -> { pair }
    array -> [ values ]
    values -> value values | &
    string -> quoted-string | unquoted-string
    quoted-string -> " LITERAL "
    unquoted-string -> NONSPACED-LITERAL

    LITERAL => <ANYTHING SUPPORTED BY THE IMPLEMENTATION>
    NONSPACED-LITERAL => <LITERAL WITHOUT SPACES>

[the golang string specification](http://golang.org/ref/spec#String_literals)

Below is the predict table:

|              production              |          stack          |
|:------------------------------------:|:-----------------------:|
|        pair -> key value pair        |   " NONSPACED-LITERAL   |
|               pair -> &              |            &            |
|             key -> string            |   " NONSPACED-LITERAL   |
|            value -> table            |            {            |
|            value -> array            |            [            |
|            value -> string           |   " NONSPACED LITERAL   |
|           table -> { pair }          |            {            |
|          array -> [ values ]         |            [            |
|        values -> value values        | { [ " NONSPACED-LITERAL |
|             values ->  &             |            &            |
|        string -> quoted-string       |            "            |
|       string -> unquoted-string      |    NONSPACED-LITERAL    |
|     quoted-string -> " LITERAL "     |            "            |
| unquoted-string -> NONSPACED-LITERAL |    NONSPACED-LITERAL    |

# TODO
 - Study implicit semi-colons to support unquoted long strings separated by spaces.
 - Write a real spec.
