package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env                  string        `env-default:"local"`
	HealthCheckDuration  time.Duration `env-default:"10s"`
	AuthServiceAddress   string        `env-default:":3201"`
	KeeperServiceAddress string        `env-default:":3202"`
	SyncServiceAddress   string        `env-default:":3204"`
	DriverPath           string
	Secret               string
}

func NewConfig() (Config, error) {
	var cfg Config

	if err := godotenv.Load(); err != nil {
		return Config{}, err
	}

	cfg.Env = os.Getenv("ENV")
	cfg.Secret = os.Getenv("SECRET")
	dur, err := time.ParseDuration(os.Getenv("HEALTHCHECK_DURATION"))
	if err != nil {
		slog.Warn("parse healthcheck failed", slog.String("default", "10s"))
		cfg.HealthCheckDuration = time.Duration(10 * time.Second)
	} else {
		cfg.HealthCheckDuration = dur
	}

	cfg.AuthServiceAddress = os.Getenv("AUTH_GRPC_LISTEN")
	cfg.KeeperServiceAddress = os.Getenv("KEEPER_GRPC_LISTEN")
	cfg.SyncServiceAddress = os.Getenv("SYNC_GRPC_LISTEN")

	return cfg, nil
}
