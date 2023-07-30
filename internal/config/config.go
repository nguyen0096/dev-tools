package config

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var (
	Cfg *Config
)

type Config struct {
	AWS AWSConfig `yaml:"aws"`
	Git GitConfig `yaml:"git"`
}

func MustLoadConfig() error {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working directory. err: %v", err)
	}

	ex, err := os.Executable()
	if err != nil {
		log.Fatalf("failed to get binary dir. err: %v", err)
	}
	binDir := filepath.Dir(ex)

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName("config_data")
	v.AddConfigPath(wd)
	v.AddConfigPath(binDir)
	v.AddConfigPath(path.Join(binDir, ".."))
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("config file not found. searched in [%s, %s]", wd, binDir)
	}

	if err := v.Unmarshal(&Cfg, func(c *mapstructure.DecoderConfig) {
		c.TagName = "json"
	}); err != nil {
		log.Fatalf("failed to unmarshal config. err : %v", err)
	}
	return nil
}
