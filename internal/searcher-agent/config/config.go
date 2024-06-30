package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Env string `yaml:"env" env-default:"local"`

	QueueName  string `yaml:"queue_name" env:"QUEUE_NAME" env-default:"searcher"`
	ConnectUrl string `yaml:"connect_url" env:"CONNECT_URL" env-required:"true"`
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
