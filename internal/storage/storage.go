package storage

import (
	"log/slog"
	"net"
	"strings"
	"time"

	"github.com/rombintu/GophKeeper/internal"
	"github.com/rombintu/GophKeeper/internal/config"
	"github.com/rombintu/GophKeeper/internal/models/auth"
	models "github.com/rombintu/GophKeeper/internal/models/storage"
	"github.com/rombintu/GophKeeper/internal/storage/drivers"
	pb "github.com/rombintu/GophKeeper/internal/storage/proto"
	"google.golang.org/grpc"
)

const (
	memDriver   = "mem"
	pgxDriver   = "pgx" // TODO
	serviceName = "StorageService"
)

func parseDriver(driverPath string) (string, string) {
	data := strings.Split(driverPath, ":")
	if len(data) == 2 {
		return data[0], data[1]
	}
	slog.Warn("parse driver failed",
		slog.String("got", driverPath),
		slog.String("default", memDriver),
	)
	return memDriver, ""
}

type Driver interface {
	UserGet(user auth.User) (auth.User, error)
	UserCreate(user auth.User) error

	SecretGet(userID int64) (models.Secret, error)
	SecretCreate(secret models.Secret) error
	SecretsGet(userID int64) ([]models.Secret, error)

	Ping() error
}

func NewDriver(driver, path string) Driver {
	switch driver { // TODO
	case memDriver:
		return NewMemDriver()
	}
	return nil
}

func NewMemDriver() Driver {
	return &drivers.MemoryDriver{}
}

type StorageService struct {
	internal.Service
	pb.UnimplementedStorageServer
	driver Driver
	config config.StorageConfig
}

func NewStorageService(cfg config.StorageConfig) *StorageService {
	driver, path := parseDriver(cfg.DriverPath)
	return &StorageService{
		config: cfg,
		driver: NewDriver(driver, path),
	}
}

func (s *StorageService) HealthCheck(duration time.Duration) {
	ticker := time.NewTicker(s.config.HealthCheckDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// TODO: отправка в API статус сервиса
			if err := s.driver.Ping(); err != nil {
				slog.Error("ping database error", slog.String("error", err.Error()))
			}
			slog.Debug("health check service", slog.String("service", serviceName))
		}
	}
}

func (s *StorageService) Start() error {
	listen, err := net.Listen(internal.TCP, s.config.Address)
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	pb.RegisterStorageServer(server, s)
	slog.Info("Service is starting",
		slog.String("service", serviceName),
		slog.String("address", s.config.Address),
	)
	return server.Serve(listen)
}

func (s *StorageService) Shutdown() error {
	return nil
}
