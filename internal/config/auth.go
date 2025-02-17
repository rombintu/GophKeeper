package config

import "os"

type AuthConfig struct {
	Config
	Address string
}

func NewAuthConfig(base Config) AuthConfig {
	return AuthConfig{
		Config: base,
	}
}

func (c *AuthConfig) Load() {
	c.Address = os.Getenv("AUTH_GRPC_LISTEN")
}
