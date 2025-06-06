package storage

import (
	"context"
	"log/slog"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
	apb "github.com/rombintu/GophKeeper/internal/proto/auth"
	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
	"github.com/rombintu/GophKeeper/internal/storage/drivers"
)

const (
	memDriver = "mem"
	pgxDriver = "postgres" // TODO
	bltDriver = "bolt"
)

func parseDriver(driverPath string) (string, string, string) {
	data := strings.Split(driverPath, "://")
	if len(data) == 2 {
		return data[0], driverPath, data[1]
	}
	slog.Warn("parse driver failed",
		slog.String("got", driverPath),
		slog.String("default", memDriver),
	)
	return memDriver, "", ""
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
	SecretCreateBatch(ctx context.Context, secrets []*kpb.Secret) error
	SecretGetBatch(ctx context.Context) ([]*kpb.Secret, error)
	SecretList(ctx context.Context, userEmail string) ([]*kpb.Secret, error)
	SecretPurge(ctx context.Context, secret *kpb.Secret) error
}

type ClientManager interface {
	SecretManager
	Set(ctx context.Context, key []byte, value []byte) error
	Get(ctx context.Context, key []byte) ([]byte, error)
	GetMap(ctx context.Context) (map[string]string, error)
}

type DriverOpts struct {
	ServiceName string
	DriverPath  string
	CryptoKey   openpgp.EntityList
}

func NewDriver(opts DriverOpts) Driver {
	driverName, driverURL, driverPathFile := parseDriver(opts.DriverPath)
	switch driverName { // TODO
	case memDriver:
		return &drivers.MemoryDriver{}
	case pgxDriver:
		return drivers.NewPgxDriver(driverURL, opts.ServiceName)
	case bltDriver:
		return drivers.NewBoltDriver(driverPathFile, opts.CryptoKey)
	}
	return nil
}

func NewUserManager(opts DriverOpts) UserManager {
	return NewDriver(opts).(UserManager)

}

func NewSecretManager(opts DriverOpts) SecretManager {
	return NewDriver(opts).(SecretManager)
}

func NewClientManager(opts DriverOpts) ClientManager {
	return NewDriver(opts).(ClientManager)
}
