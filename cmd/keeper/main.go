package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/rombintu/GophKeeper/internal/config"
	"github.com/rombintu/GophKeeper/internal/keeper"
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
	cfg := config.NewKeeperConfig(baseConfig)
	cfg.Load()
	logger.InitLogger(cfg.Env)
	common.Version(buildVersion, buildDate, buildCommit, keeper.ServiceName)

	store := storage.NewSecretManager(storage.DriverOpts{
		ServiceName: keeper.ServiceName,
		DriverPath:  cfg.DriverPath})
	service := keeper.NewKeeperService(store, cfg)

	if err := service.Configure(); err != nil {
		panic(err)
	}

	go service.HealthCheck(cfg.HealthCheckDuration)
	go func() {
		if err := service.Start(); err != nil {
			panic(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// Waiting for SIGINT (pkill -2) or SIGTERM
	<-stop

	if err := service.Shutdown(); err != nil {
		panic(err)
	}
	slog.Info("Service is shutdown", slog.String("service", "keeper"))
}
