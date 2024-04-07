package config

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server   HTTPConfig     `yaml:"server"`
	Database PostgresConfig `yaml:"postgres"`
}

type HTTPConfig struct {
	Port string `yaml:"port"`
}

type PostgresConfig struct {
	Host     string        `yaml:"host"`
	Port     string        `yaml:"port"`
	User     string        `yaml:"user"`
	Password string        `yaml:"password"`
	Database string        `yaml:"database"`
	Timeout  time.Duration `yaml:"timeout"`
}

func Load() (Config, error) {
	log.Println("read configuration file")
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return Config{}, errors.New("cfg path is not set")
	}
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return Config{}, err
	}
	log.Println("configuration loaded")
	return cfg, nil
}
