package run

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port string `yaml:"port" env:"PORT"`
}

func newDefaultConfig() *Config {
	return &Config{
		Port: "8080",
	}
}

func loadConfig(path string, getenv func(string) string) (*Config, error) {
	config := newDefaultConfig()

	// Load from yaml.
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("reading config file: %w", err)
	}
	if err == nil {
		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("parsing config file: %w", err)
		}
	}

	// Override with env vars.
	if envPort := getenv("PORT"); envPort != "" {
		config.Port = envPort
	}

	return config, nil
}
