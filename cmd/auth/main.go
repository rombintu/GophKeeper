package main

import (
	"fmt"
	"log/slog"

	"github.com/rombintu/GophKeeper/internal/config"
	"github.com/rombintu/GophKeeper/lib/common/logger"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	baseConfig, err := config.NewConfig()
	authConfig := config.NewAuthConfig(baseConfig)
	authConfig.Load()
	if err != nil {
		slog.Error("Config load error", slog.String("error", err.Error()))
	}
	logger.InitLogger(authConfig.Env)
	slog.Info(
		"Init", slog.String("Service", "auth"),
		slog.String("Build version", buildVersion),
		slog.String("Build date", buildDate),
		slog.String("Build commit", buildCommit),
	)

	slog.Debug(fmt.Sprintf("%+v", authConfig))
}
