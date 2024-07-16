package goenv

import (
	"bytes"
	"github.com/pkg/errors"
	"net"
	"os"
	"strconv"
	"testing"
	"time"
)

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
	InlineConfig `env:"inline"`
	ProjectName  string         `env:"PROJNAME;required"`
	ProjectID    int64          `env:"PROJID;default:1"`
	Version      int            `env:"VERSION"`
	VerPtr       *int           `env:"VERSION"`
	Timeout      time.Duration  `env:"TIMEOUT;gte: 5;lte: 10"`
	Peers        []*Node        `env:"PEERS;required"`
	IDs          []int          `env:"IDS;default: 1,2,3"`
	IDMap        map[int]string `env:"IDMAP"`
	PeersMap     map[int]*Node  `env:"PEERSMAP"`
	Server       *ServerConfig  `env:"SERVER"`
	Log          LogConfig      `env:"LOG"`
}

type InlineConfig struct {
	Node int `env:"NODE"`
}

type ServerConfig struct {
	ID          uint `env:"ID"`
	Port        int  `env:"PORT"`
	MultiServer bool `env:"MULTISERVER"`
}

type LogConfig struct {
	Level string  `env:"LEVEL;default:debug"`
	Path  *string `env:"PATH"`
}

func TestNewEnvParser(t *testing.T) {
	if err := os.Setenv("ENV_NODE", "1"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("ENV_PROJNAME", "goenv"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("ENV_LOG_PATH", "data"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("ENV_IDMAP", "1:127.0.0.1:8000|2:127.0.0.1:8001"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("ENV_PEERSMAP", "1:1=127.0.0.1:8000|2:2=127.0.0.1:8001"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("ENV_TIMEOUT", "7"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("ENV_PEERS", "1=127.0.0.1:8080,2=127.0.0.1:8081"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("ENV_VERSION", "1"); err != nil {
		t.Fatal(err)
	}

	c := &Config{
		ProjectID: 10,
		Server: &ServerConfig{
			ID:          1,
			MultiServer: true,
		},
		IDs: []int{5, 6, 7},
	}
	parser := NewEnvParser()
	if err := parser.Start(c); err != nil {
		t.Fatal(err)
	}

	if c.Node != 1 {
		t.Fatal("node")
	}
	if c.ProjectName != "goenv" {
		t.Fatal("ProjectName")
	}
	if c.ProjectID != 10 {
		t.Fatal("ProjectID")
	}
	if c.Server == nil || c.Server.ID != 1 || !c.Server.MultiServer {
		t.Fatal("Server")
	}
	if c.Log.Level != "debug" || *c.Log.Path != "data" {
		t.Fatal("Log")
	}

	if len(c.IDs) != 3 {
		t.Fatal("IDs")
	}
	for _, id := range c.IDs {
		if id != 5 && id != 6 && id != 7 {
			t.Fatal("IDs")
		}
	}

	if len(c.IDMap) != 2 {
		t.Fatal("IDMap")
	}
	for k, v := range c.IDMap {
		switch k {
		case 1:
			if v != "127.0.0.1:8000" {
				t.Fatalf("IDMap %d", k)
			}
		case 2:
			if v != "127.0.0.1:8001" {
				t.Fatalf("IDMap %d", k)
			}
		default:
			t.Fatal("IDMap")
		}
	}

	if len(c.PeersMap) != 2 {
		t.Fatal("PeersMap")
	}
	for k, v := range c.PeersMap {
		switch k {
		case 1:
			if v.ID != 1 || v.Addr.String() != "127.0.0.1:8000" {
				t.Fatalf("PeersMap %d", k)
			}
		case 2:
			if v.ID != 2 || v.Addr.String() != "127.0.0.1:8001" {
				t.Fatalf("PeersMap %d", k)
			}
		default:
			t.Fatal("PeersMap")
		}
	}

	if len(c.Peers) != 2 {
		t.Fatal("Peers")
	}
	for _, info := range c.Peers {
		switch info.ID {
		case 1:
			if info.Addr.String() != "127.0.0.1:8080" {
				t.Fatalf("Peers %d", info.ID)
			}
		case 2:
			if info.Addr.String() != "127.0.0.1:8081" {
				t.Fatalf("Peers %d", info.ID)
			}
		default:
			t.Fatal("Peers")
		}
	}

	if c.Version != 1 || c.VerPtr == nil || *c.VerPtr != 1 {
		t.Fatal("Version")
	}
}

func Test_Required(t *testing.T) {
	type InnerConfig struct {
		ProjectName string `env:"PROJNAME;required"`
	}

	parser := NewEnvParser()

	c1 := &InnerConfig{}
	if err := parser.Start(&c1); err == nil {
		t.Fatal("ProjectName should be required")
	}

	c2 := &InnerConfig{
		ProjectName: "aaa",
	}
	if err := parser.Start(&c2); err != nil || c2.ProjectName != "aaa" {
		t.Fatal(err)
	}

	c3 := &InnerConfig{}
	if err := os.Setenv("ENV_PROJNAME", "aaa"); err != nil {
		t.Fatal(err)
	}
	if err := parser.Start(&c3); err != nil || c3.ProjectName != "aaa" {
		t.Fatal(err)
	}
}

func Test_Name(t *testing.T) {
	type InnerConfig struct {
		Name string `env:"name: NAME"`
		Age  int    `env:"AGE"`
	}

	if err := os.Setenv("ENV_NAME", "aaa"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("ENV_AGE", "18"); err != nil {
		t.Fatal(err)
	}

	c := &InnerConfig{}
	parser := NewEnvParser()
	if err := parser.Start(&c); err != nil {
		t.Fatal(err)
	}

	if c.Name != "aaa" || c.Age != 18 {
		t.Fatal("failed")
	}
}
