package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	Cfg *Config
)

type Config struct {
	AWS AWSConfig `yaml:"aws"`
}

const (
	configFile = "config_data.yml"
)

func LoadConfig() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	configPath := filepath.Join(wd, configFile)
	f, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(f, &Cfg); err != nil {
		return err
	}

	return nil
}
