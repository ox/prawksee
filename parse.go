package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Server struct {
	Bind    int
	Network string
	Address string
}

type Config struct {
	Servers map[string]Server
}

func ParseConfig(filepath string) (*Config, error) {
	stat, err := os.Stat(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not open config: %w", err)
	}

	if stat.IsDir() {
		return nil, fmt.Errorf("config is not a file")
	}

	var config Config
	if _, err := toml.DecodeFile(filepath, &config); err != nil {
		return nil, fmt.Errorf("could not decode config file: %w", err)
	}

	return &config, nil
}
