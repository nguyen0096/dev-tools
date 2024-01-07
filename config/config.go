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

type AWSConfig struct {
	MFAs []MFA `yaml:"mfas"`
}

type MFA struct {
	// Name of initial profile to do the MFA
	Profile         string `json:"profile"`
	Device          string `json:"device"`
	SessionDuration int    `json:"session_duration"`
	// Name of the profile with the MFA session
	OutputProfile string `json:"output_profile"`
	// The name of the cluster for which to create a kubeconfig entry
	ClusterName string `json:"cluster_name"`
}

type GitConfig struct {
	Repos []Repository `json:"repositories"`
}

type Repository struct {
	Path            string `json:"path"`
	MessageTemplate string `json:"message_template"`
	DateFormat      string `json:"date_format"`
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
