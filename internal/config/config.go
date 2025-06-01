package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Address          string        `yaml:"address"`
	StoragePath      string        `yaml:"storage_path"`
	PublicKeyPath    string        `yaml:"public_key_path"`
}


func Load(env string) *Config {
	var configPath string
	switch env {
	case "local":
		configPath = "./config/local.example.yaml"
	case "prod":
		configPath = "./config/prod.example.yaml"
	case "dev":
		configPath = "./config/dev.example.yaml"
	default:
		log.Fatalf("unknown environment:")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("%v Is Not Exist", configPath)
	}
	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("can't read config: %v", err)
	}

	return &config
}
