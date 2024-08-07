# goenv

Reflect environment values to config struct

## Install

```shell
go get github.com/hjhsamuel/goenv
```

## Usage

Support type:

- `string`
- `bool`
- `int` `int8` `int16` `int32` `int64`
- `uint` `uint8` `uint16` `uint32` `uint64`
- `float32` `float64`
- `complex64` `complex128`
- `pointer`
- `struct`
- `slice`
- `map`

Support tags:

- `required`
- `default`
- `inline`
- `-`
- `name`
- `lt` and `lte`
- `gt` and `gte`

The default config like:

```go
type ServerConfig struct {
	ID          uint            `env:"ID"`
	Port        int             `env:"PORT;gte:1000;lte:60000"`
	MultiServer bool            `env:"MULTISERVER"`
	Peers       map[int]string  `env:"PEERS"`
	IDs         []int           `env:IDS`
}

type LogConfig struct {
	Level string  `env:"LEVEL;default:debug"`
	Path  *string `env:"PATH"`
}

type IPAddress struct {
	IP      string `env:"name:IP"`
	Gateway string `env:"name:GATEWAY"`
}

type IPv4 struct {
	IPAddress   `env:inline`
}

type Config struct {
	ProjectName string        `env:"PROJNAME;required"`
	ProjectID   int64         `env:"PROJID;default:1"`
	Server      *ServerConfig `env:"SERVER"`
	Log         LogConfig     `env:"LOG"`
	Address     IPv4          `env:"name:IPV4"`
}
```

So, an environment value could be found like `ENV_PROJNAME`

Some function could be called to set the custom tag:

- `SetPrefix` or `WithPrefix`

    the default prefix is `ENV`

- `SetSplitChar` or `WithSplitChar`

    the default split char is `_`

- `SetTag` or `WithTag`

    the default tag name is `env`

We could parse config struct like:

```shell
ENV_PROJNAME=goenv
ENV_SERVER_PEERS="1:127.0.0.1:8000|2:127.0.0.1:8001"
ENV_SERVER_PORT=8000
ENV_SERVER_MULTISERVER=true
ENV_SERVER_IDS="1,2"
```

```go
c := &Config{ProjectName: "goenv"}
parser := NewEnvParser()
if err := parser.Start(&c); err != nil {
    log.Fatal(err)
}
```

## Custom Unmarshal

You can use your own function to format fields by implement the `UnmarshalText` interface

```go
type Node struct {
	ID   int
	Addr net.Addr
}

func (n *Node) UnmarshalText(text []byte) error {
	l := bytes.Split(text, []byte("="))
	if len(l) != 2 {
		return errors.New("invalid node format")
	}

	if id, err := strconv.Atoi(string(l[0])); err != nil {
		return err
	} else {
		n.ID = id
	}

	if host, port, err := net.SplitHostPort(string(l[1])); err != nil {
		return err
	} else {
		iport, err := strconv.Atoi(port)
		if err != nil {
			return err
		}
		n.Addr = &net.TCPAddr{IP: net.ParseIP(host), Port: iport}
	}
	return nil
}

type Config struct {
	Node Node `env:"NODE"`
}
```

```shell
ENV_NODE="1=127.0.0.1:8000"
```
