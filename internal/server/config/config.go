package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config represents the application's configuration.

type Config struct {
	Env string `yaml:"env" env-default:"local"`

	Server struct {
		Port           int      `yaml:"port" env-default:"8080"`
		AllowedOrigins []string `yaml:"allowed_origins"`
	} `yaml:"server"`

	Database struct {
		URI  string `yaml:"uri" env-required:"true"`
		Name string `yaml:"name" env-required:"true"`
	} `yaml:"database"`

	Auth struct {
		ConfigPath string `yaml:"config_path" env-required:"true"`
	} `yaml:"auth"`

	GRPC struct {
		Addr string `yaml:"addr" env-default:"localhost:44044"`
	} `yaml:"grpc"`
}

// MustLoad loads the configuration from the specified path and environment variables.
func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file doesn't exist: " + configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(fmt.Errorf("failed to read config: %w", err))
	}

	return &cfg
}

// fetchConfigPath fetches the config path from command-line flags or environment variables.
// Priority: command-line flag > environment variable.
func fetchConfigPath() string {
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to the config file")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	if configPath == "" {
		panic("no config file path provided")
	}

	return configPath
}
