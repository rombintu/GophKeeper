package auth

import (
	"context"
	"log/slog"
	"net"
	"time"

	"github.com/rombintu/GophKeeper/internal"
	"github.com/rombintu/GophKeeper/internal/config"
	apb "github.com/rombintu/GophKeeper/internal/proto/auth"
	"github.com/rombintu/GophKeeper/internal/storage"
	"github.com/rombintu/GophKeeper/lib/common"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

const (
	ServiceName = "AuthService"
)

type AuthService struct {
	apb.UnimplementedAuthServer
	store  storage.UserManager
	config config.AuthConfig
}

func NewAuthService(store storage.UserManager, cfg config.AuthConfig) internal.Service {
	return &AuthService{
		store:  store,
		config: cfg,
	}
}

func (s *AuthService) HealthCheck(duration time.Duration) {
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

func (s *AuthService) Start() error {
	listen, err := net.Listen(internal.TCP, s.config.AuthServiceAddress)
	if err != nil {
		return err
	}
	// TODO: конфигурация и унификация для сервисов
	limiter := rate.NewLimiter(rate.Limit(10), 20)
	server := grpc.NewServer(grpc.UnaryInterceptor(
		common.RateLimitInterceptor(limiter),
	))
	apb.RegisterAuthServer(server, s)
	slog.Info("Service is starting",
		slog.String("service", ServiceName),
		slog.String("address", s.config.AuthServiceAddress),
	)
	return server.Serve(listen)
}

func (s *AuthService) Shutdown() error {
	return s.store.Close(context.Background())
}

func (s *AuthService) Configure() error {
	ctx := context.Background()
	if err := s.store.Open(ctx); err != nil {
		return err
	}
	if err := s.store.Configure(ctx); err != nil {
		return err
	}
	return nil
}
