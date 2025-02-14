package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env string `env-default:"local"`
}

func NewConfig() (Config, error) {
	var cfg Config

	if err := godotenv.Load(); err != nil {
		return Config{}, err
	}

	cfg.Env = os.Getenv("ENV")

	return cfg, nil
}
