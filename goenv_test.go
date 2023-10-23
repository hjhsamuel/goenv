package goenv

import (
	"os"
	"testing"
)

func TestNewEnvParser(t *testing.T) {
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

	if err := os.Setenv("ENV_ProjName", "goenv"); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("ENV_Log_Path", "data"); err != nil {
		t.Fatal(err)
	}

	c := &Config{
		ProjectID: 10,
		Server: &ServerConfig{
			ID:          1,
			MultiServer: true,
		},
	}
	parser := NewEnvParser()
	if err := parser.Start(&c); err != nil {
		t.Fatal(err)
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
}
