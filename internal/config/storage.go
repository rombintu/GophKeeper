package config

import "os"

type StorageConfig struct {
	Config
	DriverPath string
}

func NewStorageConfig(base Config) StorageConfig {
	return StorageConfig{
		Config: base,
	}
}

func (c *StorageConfig) Load() {
	c.DriverPath = os.Getenv("STORAGE_DRIVER_PATH")
}
