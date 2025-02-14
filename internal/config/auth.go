package config

import "os"

type AuthConfig struct {
	Config
	AuthGRPCPort string
}

func NewAuthConfig(base Config) AuthConfig {
	return AuthConfig{
		Config: base,
	}
}

func (c *AuthConfig) Load() {
	c.AuthGRPCPort = os.Getenv("AUTH_GRPC_PORT")
}
