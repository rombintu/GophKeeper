package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/rombintu/GophKeeper/lib/jwt"
)

type AuthConfig struct {
	Config
	TokenExp time.Duration
}

func NewAuthConfig(base Config) AuthConfig {
	return AuthConfig{
		Config: base,
	}
}

func (c *AuthConfig) Load() {
	c.DriverPath = os.Getenv("AUTH_DRIVER_PATH")

	dur, err := time.ParseDuration(os.Getenv("AUTH_TOKEN_EXPIRE"))
	if err != nil {
		slog.Warn("parse healthcheck failed", slog.String("default", "10m"))
		c.TokenExp = time.Duration(10 * time.Minute)
	} else {
		c.TokenExp = dur
	}
	c.Secret, err = jwt.GenerateHMACSecret(32)
	if err != nil {
		c.Secret = ""
		slog.Warn("generate secret failed", slog.String("default", "empty string"))
	}
}
