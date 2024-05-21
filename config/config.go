package config

import (
	"io"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type (
	Config struct {
		Application Application `toml:"application"`
		Server      Server      `toml:"server"`
	}

	Application struct {
		EnvMode  string `toml:"env_mode"`
		LogLevel string `toml:"log_level"`
	}

	Server struct {
		Host string `toml:"host"`
		Port string `toml:"port"`
	}
)

func Load(filepath string) (*Config, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = toml.Unmarshal(b, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
