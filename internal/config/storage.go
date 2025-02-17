package config

import "os"

type StorageConfig struct {
	Config
	Address    string
	DriverPath string
}

func NewStorageConfig(base Config) StorageConfig {
	return StorageConfig{
		Config: base,
	}
}

func (c *StorageConfig) Load() {
	c.Address = os.Getenv("STORAGE_GRPC_LISTEN")
	c.DriverPath = os.Getenv("STORAGE_DRIVER_PATH")
}
