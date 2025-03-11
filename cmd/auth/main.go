package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/rombintu/GophKeeper/internal/auth"
	"github.com/rombintu/GophKeeper/internal/config"
	"github.com/rombintu/GophKeeper/internal/storage"
	"github.com/rombintu/GophKeeper/lib/common"
	"github.com/rombintu/GophKeeper/lib/jwt"
	"github.com/rombintu/GophKeeper/lib/logger"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	genSecret := flag.Bool("secret", false, "Generate secret and exit")
	flag.Parse()

	if *genSecret {
		newSecret, err := jwt.GenerateHMACSecret(32)
		if err != nil {
			fmt.Println("error generate secret")
			os.Exit(1)
		}
		fmt.Printf("Generated: %s\n", newSecret)
		os.Exit(0)
	}

	baseConfig, err := config.NewConfig()
	if err != nil {
		slog.Error("Config load error", slog.String("error", err.Error()))
		os.Exit(0)
	}
	cfg := config.NewAuthConfig(baseConfig)
	cfg.Load()
	logger.InitLogger(cfg.Env)
	common.Version(buildVersion, buildDate, buildCommit, auth.ServiceName)

	store := storage.NewUserManager(storage.DriverOpts{
		ServiceName: auth.ServiceName,
		DriverPath:  cfg.DriverPath,
	})
	service := auth.NewAuthService(store, cfg)

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
	slog.Info("Service is shutdown", slog.String("service", "auth"))
}
