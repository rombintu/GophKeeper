package sync

import (
	"log/slog"
	"net"
	"time"

	"github.com/rombintu/GophKeeper/internal"
	"github.com/rombintu/GophKeeper/internal/config"
	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
	spb "github.com/rombintu/GophKeeper/internal/proto/sync"
	"github.com/rombintu/GophKeeper/lib/connections"
	"github.com/rombintu/GophKeeper/lib/jwt"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

const (
	ServiceName = "SyncService"
)

type SyncService struct {
	spb.UnimplementedSyncServer
	config config.SyncConfig
	keeper kpb.KeeperClient
}

type SyncServiceOpts struct {
	KeeperClient kpb.KeeperClient
}

func NewSyncService(cfg config.SyncConfig, opts SyncServiceOpts) internal.Service {
	return &SyncService{
		config: cfg,
		keeper: opts.KeeperClient,
	}
}

func (s *SyncService) HealthCheck(duration time.Duration) {
	ticker := time.NewTicker(s.config.HealthCheckDuration)
	defer ticker.Stop()

	for range ticker.C {

		// TODO: отправка в API статус сервиса
		slog.Debug("health check service", slog.String("service", ServiceName))
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
			connections.RateLimitInterceptor(limiter),
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
	return nil
}

func (s *SyncService) Configure() error {
	return nil
}
