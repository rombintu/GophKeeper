package config

import (
	"encoding/json"
	"log/slog"
	"os"
)

type ClientConfig struct {
	KeyPath string
}

func NewClientConfig(confPath string) (*ClientConfig, error) {
	// Инициализация конфига с значениями по умолчанию
	config := ClientConfig{}
	// Чтение файла
	data, err := os.ReadFile(confPath)
	if err != nil {
		slog.Warn("file not found", slog.String("error", err.Error()))
		return &config, nil
	}

	// Парсинг JSON
	if err := json.Unmarshal(data, &config); err != nil {
		slog.Warn("failed parsing", slog.String("error", err.Error()))
		return &config, nil
	}

	return &config, nil
}

func (c *ClientConfig) Save(confPath string) error {
	data, err := json.Marshal(&c)
	if err != nil {
		return err
	}
	return os.WriteFile(confPath, data, os.ModePerm)
}
