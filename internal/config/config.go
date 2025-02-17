package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env                 string        `env-default:"local"`
	HealthCheckDuration time.Duration `env-default:"10s"`
}

func NewConfig() (Config, error) {
	var cfg Config

	if err := godotenv.Load(); err != nil {
		return Config{}, err
	}

	cfg.Env = os.Getenv("ENV")
	dur, err := time.ParseDuration(os.Getenv("HEALTHCHECK_DURATION"))
	if err != nil {
		slog.Warn("parse healthcheck failed", slog.String("default", "10s"))
		cfg.HealthCheckDuration = time.Duration(10 * time.Second)
	} else {
		cfg.HealthCheckDuration = dur
	}

	return cfg, nil
}
