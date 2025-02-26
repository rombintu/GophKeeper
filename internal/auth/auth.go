package auth

import (
	"log/slog"
	"net"
	"time"

	"github.com/rombintu/GophKeeper/internal"
	"github.com/rombintu/GophKeeper/internal/config"
	"github.com/rombintu/GophKeeper/internal/proto"
	"github.com/rombintu/GophKeeper/internal/storage"
	"github.com/rombintu/GophKeeper/lib/common"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

const (
	ServiceName = "AuthService"
)

type AuthService struct {
	proto.UnimplementedAuthServer
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

	for {
		select {
		case <-ticker.C:
			// TODO: отправка в API статус сервиса
			slog.Debug("health check service", slog.String("service", ServiceName))
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
	proto.RegisterAuthServer(server, s)
	slog.Info("Service is starting",
		slog.String("service", ServiceName),
		slog.String("address", s.config.AuthServiceAddress),
	)
	return server.Serve(listen)
}

func (s *AuthService) Shutdown() error {
	return s.store.Close()
}

func (s *AuthService) Configure() error {
	return s.store.Open()
}
