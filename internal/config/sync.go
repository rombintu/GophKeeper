package config

import (
	"os"
)

type SyncConfig struct {
	Config
}

func NewSyncConfig(base Config) SyncConfig {
	return SyncConfig{
		Config: base,
	}
}

func (c *SyncConfig) Load() {
	c.DriverPath = os.Getenv("SYNC_DRIVER_PATH")
}
