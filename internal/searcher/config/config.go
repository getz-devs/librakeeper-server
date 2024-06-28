package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env string `yaml:"env" env-default:"local"`

	StoragePath string `yaml:"storage_path" env-required:"true"`

	GRPC GRPCConfig `yaml:"grpc"`

	DatabaseMongo DatabaseMongoConfig `yaml:"database_mongo"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env-default:"44044"`
	Timeout time.Duration `yaml:"timeout" env-default:"10h"`
}

type DatabaseMongoConfig struct {
	ConnectURL     string `yaml:"connect_url" env-required:"true"`
	DatabaseName   string `yaml:"database_name" env-required:"true"`
	CollectionName string `yaml:"collection_name_books" env-required:"true"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file doesn't exist: " + path)
	}

	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable
// Priority: command line flag > environment variable > default
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	//if res == "" {
	//	res = "config.yml"
	//}

	return res
}
