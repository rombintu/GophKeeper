package main

import (
	"log/slog"

	"github.com/rombintu/GophKeeper/internal/auth"
	"github.com/rombintu/GophKeeper/internal/config"
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
	}
	authConfig := config.NewAuthConfig(baseConfig)
	authConfig.Load()
	logger.InitLogger(authConfig.Env)
	slog.Info(
		"Init", slog.String("Binary", "auth"),
		slog.String("Build version", buildVersion),
		slog.String("Build date", buildDate),
		slog.String("Build commit", buildCommit),
	)

	service := auth.NewAuthService(authConfig)
	go service.HealthCheck(authConfig.HealthCheckDuration)
	if err := service.Start(); err != nil {
		panic(err)
	}
}
