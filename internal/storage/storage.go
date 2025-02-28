package storage

import (
	"context"
	"log/slog"
	"strings"

	apb "github.com/rombintu/GophKeeper/internal/proto/auth"
	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
	"github.com/rombintu/GophKeeper/internal/storage/drivers"
)

const (
	memDriver = "mem"
	pgxDriver = "postgres" // TODO
)

func parseDriver(driverPath string) (string, string) {
	data := strings.Split(driverPath, "://")
	if len(data) == 2 {
		return data[0], driverPath
	}
	slog.Warn("parse driver failed",
		slog.String("got", driverPath),
		slog.String("default", memDriver),
	)
	return memDriver, ""
}

type Driver interface {
	Open(ctx context.Context) error
	Close(ctx context.Context) error
	Ping(ctx context.Context, monitoring bool) error
	Configure(ctx context.Context) error
}

type UserManager interface {
	Driver
	UserGet(ctx context.Context, user *apb.User) (*apb.User, error)
	UserCreate(ctx context.Context, user *apb.User) error
}

type SecretManager interface {
	Driver
	SecretCreate(ctx context.Context, secret *kpb.Secret) error
	SecretList(ctx context.Context, userEmail string) ([]*kpb.Secret, error)
	SecretPurge(ctx context.Context, secret *kpb.Secret) error
}

func NewDriver(driverPath, serviceName string) Driver {
	driverName, driverURL := parseDriver(driverPath)
	switch driverName { // TODO
	case memDriver:
		return &drivers.MemoryDriver{}
	case pgxDriver:
		return drivers.NewPgxDriver(driverURL, serviceName)
	}
	return nil
}

func NewUserManager(driverPath, serviceName string) UserManager {
	return NewDriver(driverPath, serviceName).(UserManager)

}

func NewSecretManager(driverPath, serviceName string) SecretManager {
	return NewDriver(driverPath, serviceName).(SecretManager)

}
