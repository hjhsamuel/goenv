package main

import (
	"bytes"
	"fmt"
	"github.com/hjhsamuel/goenv"
	"github.com/pkg/errors"
	"net"
	"os"
	"strconv"
	"time"
)

type Node struct {
	ID   int
	Addr net.Addr
}

func (n *Node) UnmarshalText(text []byte) error {
	fmt.Println("===========> ", string(text))
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
	Version int           `env:"VERSION;required"`
	Timeout time.Duration `env:"TIMEOUT;required;gte: 5;lte: 10"`
	Peers   []*Node       `env:"PEERS;required"`
	IDs     []int         `env:"IDS"`
}

func main() {
	if err := os.Setenv("ENV_VERSION", "1"); err != nil {
		panic(err)
	}
	if err := os.Setenv("ENV_TIMEOUT", "5"); err != nil {
		panic(err)
	}
	if err := os.Setenv("ENV_PEERS", "1=127.0.0.1:8080,2=127.0.0.1:8081"); err != nil {
		panic(err)
	}
	if err := os.Setenv("ENV_IDS", "1,2,3"); err != nil {
		panic(err)
	}

	c := &Config{}

	parser := goenv.NewEnvParser()
	if err := parser.Start(&c); err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", *c)
	for _, info := range c.Peers {
		fmt.Println(info.ID, info.Addr.String())
	}
}
