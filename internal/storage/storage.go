package storage

import (
	"log/slog"
	"strings"

	"github.com/rombintu/GophKeeper/internal/proto"
	"github.com/rombintu/GophKeeper/internal/storage/drivers"
)

const (
	memDriver = "mem"
	pgxDriver = "pgx" // TODO
	// ServiceName = "StorageService"
)

func ParseDriver(driverPath string) (string, string) {
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
	Open() error
	Close() error
	Ping() error
}

type UserManager interface {
	Driver
	UserGet(user *proto.User) (*proto.User, error)
	UserCreate(user *proto.User) error
}

type SecretManager interface {
	Driver
	SecretGet(userID int64) (*proto.Secret, error)
	SecretCreate(secret *proto.Secret) error
	SecretList(userID int64) ([]*proto.Secret, error)
}

func NewUserManager(driverPath string) UserManager {
	driverName, _ := ParseDriver(driverPath)
	switch driverName { // TODO
	case memDriver:
		return &drivers.MemoryDriver{}
	}
	return nil
}
