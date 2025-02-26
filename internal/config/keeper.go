package config

import (
	"os"
)

type KeeperConfig struct {
	Config
}

func NewKeeperConfig(base Config) KeeperConfig {
	return KeeperConfig{
		Config: base,
	}
}

func (c *KeeperConfig) Load() {
	c.DriverPath = os.Getenv("KEEPER_DRIVER_PATH")
}
