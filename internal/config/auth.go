package config

type AuthConfig struct {
	Config
}

func NewAuthConfig(base Config) AuthConfig {
	return AuthConfig{
		Config: base,
	}
}

func (c *AuthConfig) Load() {
}
