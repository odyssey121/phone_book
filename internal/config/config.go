package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string  `yaml: "env" env: "ENV" env-default: "development"`
	Storage    Storage `yaml: "storage"`
	HTTPServer `yaml:"http_server"`
}

type Storage struct {
	Host        string `yaml:"host" env: "host" env-default: "localhost"`
	Port        int    `yaml:"port" env: "port" env-default: "5432"`
	Driver      string `yaml:"driver" env: "driver" env-default: "json_driver"`
	Username    string `yaml:"username" env: "username"`
	Password    string `yaml:"password" env: "password"`
	Database    string `yaml:"database" env: "database"`
	StoragePath string `yaml:"storage_path" env: "storage_path" env-default:"internal/store/data/store.json"`
	IndexesPath string `yaml:"indexes_path" env: "indexes_path" env-default:"internal/store/data/indexes.json"`
}

type HTTPServer struct {
	Address     string        `yaml: "address" env-default: "localhost:1234"`
	Timeout     time.Duration `yaml: "timeout" env-default: "4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default: "60s"`
}

func MustLoad() Config {
	configPath := "config/local.yaml"
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return cfg
}
