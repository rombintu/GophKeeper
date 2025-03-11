package keeper

import (
	"context"
	"log/slog"
	"net"
	"time"

	"github.com/rombintu/GophKeeper/internal"
	"github.com/rombintu/GophKeeper/internal/config"
	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
	"github.com/rombintu/GophKeeper/internal/storage"
	"github.com/rombintu/GophKeeper/lib/common"
	"github.com/rombintu/GophKeeper/lib/jwt"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

const (
	ServiceName = "KeeperService"
)

type KeeperService struct {
	kpb.UnimplementedKeeperServer
	store  storage.SecretManager
	config config.KeeperConfig
}

func NewKeeperService(store storage.SecretManager, cfg config.KeeperConfig) internal.Service {
	return &KeeperService{
		store:  store,
		config: cfg,
	}
}

func (s *KeeperService) HealthCheck(duration time.Duration) {
	ticker := time.NewTicker(s.config.HealthCheckDuration)
	defer ticker.Stop()

	for range ticker.C {

		// TODO: отправка в API статус сервиса
		slog.Debug("health check service", slog.String("service", ServiceName))
		ctx := context.Background()
		if err := s.store.Ping(ctx, true); err != nil {
			slog.Warn("ping failed", slog.String("error", err.Error()))
		}
	}
}

func (s *KeeperService) Start() error {
	listen, err := net.Listen(internal.TCP, s.config.KeeperServiceAddress)
	if err != nil {
		return err
	}
	// TODO: конфигурация и унификация для сервисов
	limiter := rate.NewLimiter(rate.Limit(10), 20)

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			common.RateLimitInterceptor(limiter),
			jwt.VerifyTokenInterceptor(s.config.Secret, []string{}),
		),
	}
	server := grpc.NewServer(opts...)
	kpb.RegisterKeeperServer(server, s)
	slog.Info("Service is starting",
		slog.String("service", ServiceName),
		slog.String("address", s.config.KeeperServiceAddress),
	)
	return server.Serve(listen)
}

func (s *KeeperService) Shutdown() error {
	return s.store.Close(context.Background())
}

func (s *KeeperService) Configure() error {
	return s.store.Open(context.Background())
}
