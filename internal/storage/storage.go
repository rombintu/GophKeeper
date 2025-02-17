package storage

import (
	"log/slog"
	"net"
	"strings"
	"time"

	"github.com/rombintu/GophKeeper/internal"
	"github.com/rombintu/GophKeeper/internal/config"
	pb "github.com/rombintu/GophKeeper/internal/proto"
	"github.com/rombintu/GophKeeper/internal/storage/drivers"
	"google.golang.org/grpc"
)

const (
	memDriver   = "mem"
	pgxDriver   = "pgx" // TODO
	ServiceName = "StorageService"
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
	UserGet(user *pb.User) error
	UserCreate(user *pb.User) error

	SecretGet(userID int64) (*pb.Secret, error)
	SecretCreate(secret *pb.Secret) error
	SecretsGet(userID int64) ([]*pb.Secret, error)

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
			slog.Debug("health check service", slog.String("service", ServiceName))
		}
	}
}

func (s *StorageService) Start() error {
	listen, err := net.Listen(internal.TCP, s.config.StorageServiceAddress)
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	pb.RegisterStorageServer(server, s)
	slog.Info("Service is starting",
		slog.String("service", ServiceName),
		slog.String("address", s.config.StorageServiceAddress),
	)
	return server.Serve(listen)
}

func (s *StorageService) Shutdown() error {
	return nil
}
