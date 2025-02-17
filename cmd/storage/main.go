package main

import (
	"log/slog"
	"os"

	"github.com/rombintu/GophKeeper/internal/config"
	"github.com/rombintu/GophKeeper/internal/storage"
	"github.com/rombintu/GophKeeper/lib/common"
	"github.com/rombintu/GophKeeper/lib/logger"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	baseConfig, err := config.NewConfig()
	if err != nil {
		slog.Error("Config load error", slog.String("error", err.Error()))
		os.Exit(0)
	}
	cfg := config.NewStorageConfig(baseConfig)
	cfg.Load()
	logger.InitLogger(cfg.Env)
	common.Version(buildVersion, buildDate, buildCommit, "storage")
	service := storage.NewStorageService(cfg)
	go service.HealthCheck(cfg.HealthCheckDuration)
	if err := service.Start(); err != nil {
		panic(err)
	}
}
