# goenv

Reflect environment values to config struct

## Install

```shell
go get github.com/hjhsamuel/goenv
```

## Usage

Support tags:

- `required`
- `default:`

The default config like:

```go
type ServerConfig struct {
	ID          uint `env:"ID"`
	Port        int  `env:"Port"`
	MultiServer bool `env:"MultiServer"`
}

type LogConfig struct {
	Level string  `env:"Level;default:debug"`
	Path  *string `env:"Path"`
}

type Config struct {
	ProjectName string        `env:"ProjName;required"`
	ProjectID   int64         `env:"ProjID;default:1"`
	Server      *ServerConfig `env:"Server"`
	Log         LogConfig     `env:"Log"`
}
```

So, an environment value could be found like `ENV_ProjName`

Some function could be called to set the custom tag:

- `SetPrefix`

    the default prefix is `ENV`
- `SetSplitChar`

    the default split char is `_`
- `SetTag`

    the default tag name is `env`

We could parse config struct like:

```go
c := &Config{ProjectName: "goenv"}
parser := NewEnvParser()
if err := parser.Start(&c); err != nil {
    log.Fatal(err)
}
```