package sync

import (
	"context"
	"log/slog"
	"net"
	"time"

	"github.com/rombintu/GophKeeper/internal"
	"github.com/rombintu/GophKeeper/internal/config"
	spb "github.com/rombintu/GophKeeper/internal/proto/sync"
	"github.com/rombintu/GophKeeper/internal/storage"
	"github.com/rombintu/GophKeeper/lib/common"
	"github.com/rombintu/GophKeeper/lib/jwt"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

const (
	ServiceName = "SyncService"
)

type SyncService struct {
	spb.UnimplementedSyncServer
	store  storage.ClientManager
	config config.SyncConfig
}

func NewSyncService(store storage.ClientManager, cfg config.SyncConfig) internal.Service {
	return &SyncService{
		store:  store,
		config: cfg,
	}
}

func (s *SyncService) HealthCheck(duration time.Duration) {
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

func (s *SyncService) Start() error {
	listen, err := net.Listen(internal.TCP, s.config.SyncServiceAddress)
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
	spb.RegisterSyncServer(server, s)
	slog.Info("Service is starting",
		slog.String("service", ServiceName),
		slog.String("address", s.config.SyncServiceAddress),
	)
	return server.Serve(listen)
}

func (s *SyncService) Shutdown() error {
	return s.store.Close(context.Background())
}

func (s *SyncService) Configure() error {
	return s.store.Open(context.Background())
}
