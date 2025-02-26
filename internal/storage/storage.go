package storage

import (
	"log/slog"
	"strings"

	apb "github.com/rombintu/GophKeeper/internal/proto/auth"
	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
	"github.com/rombintu/GophKeeper/internal/storage/drivers"
)

const (
	memDriver = "mem"
	pgxDriver = "pgx" // TODO
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
	UserGet(user *apb.User) (*apb.User, error)
	UserCreate(user *apb.User) error
}

type SecretManager interface {
	Driver
	SecretCreate(secret *kpb.Secret) error
	SecretList(userID int64, pattern string) ([]*kpb.Secret, error)
	SecretPurge(userID int64, secret *kpb.Secret) error
}

// TODO: Унификация. уменьшение кода
func NewUserManager(driverPath string) UserManager {
	driverName, _ := ParseDriver(driverPath)
	switch driverName { // TODO
	case memDriver:
		return &drivers.MemoryDriver{}
	}
	return nil
}

func NewSecretManager(driverPath string) SecretManager {
	driverName, _ := ParseDriver(driverPath)
	switch driverName { // TODO
	case memDriver:
		return &drivers.MemoryDriver{}
	}
	return nil
}
