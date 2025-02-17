package auth

import (
	"log/slog"
	"net"
	"time"

	"github.com/rombintu/GophKeeper/internal"
	"github.com/rombintu/GophKeeper/internal/config"
	pb "github.com/rombintu/GophKeeper/internal/proto"
	"google.golang.org/grpc"
)

const (
	ServiceName = "AuthService"
)

type AuthService struct {
	pb.UnimplementedAuthServer
	config config.AuthConfig
}

func NewAuthService(cfg config.AuthConfig) internal.Service {
	return &AuthService{
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
	server := grpc.NewServer()
	pb.RegisterAuthServer(server, s)
	slog.Info("Service is starting",
		slog.String("service", ServiceName),
		slog.String("address", s.config.AuthServiceAddress),
	)
	return server.Serve(listen)
}

func (s *AuthService) Shutdown() error {
	return nil
}
