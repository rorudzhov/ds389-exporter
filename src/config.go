package src

import (
	"errors"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Server       string `yaml:"server"`
	Port         string `yaml:"port"`
	BindDN       string `yaml:"bind_dn"`
	BindPassword string `yaml:"bind_password"`
}

func GetConfig(path string) (*Config, error) {
	config := Config{}

	// Read file
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse YAML
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	// Check for empty value
	if config.Server == "" || config.Port == "" || config.BindDN == "" || config.BindPassword == "" {
		return nil, errors.New("empty or invalid values found in configuration file")
	}

	return &config, err
}
